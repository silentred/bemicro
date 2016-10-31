package main

import (
	"bemicro/discovery"
	"bemicro/middleware"
	"bemicro/proto"
	"flag"
	"fmt"
	"log"
	"net"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

const srvName = "greeter"
const etcdHost = "http://localhost:2379"

var port string
var id uint64

func init() {
	flag.StringVar(&port, "p", ":1234", "listen port")
	flag.Uint64Var(&id, "i", 0, "service ID")
}

func main() {
	flag.Parse()

	// register service
	service := discovery.NewService(id, srvName, fmt.Sprintf("localhost%s", port))
	publisher := discovery.NewEtcdPublisher([]string{etcdHost}, 10)
	publisher.Register(service)
	go publisher.Heartbeat(service)

	// use interceptors...
	chain := middleware.UnaryInterceptorChain(middleware.Recovery, middleware.Logging,
		middleware.Auth, middleware.UserVerify)
	opt := grpc.UnaryInterceptor(chain)

	// start grpc service
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(opt)
	proto.RegisterGreeterServer(s, &server{})
	s.Serve(lis)
}

type server struct{}

// Hello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *proto.HelloReq) (*proto.HelloResp, error) {
	// log ctx values)
	return &proto.HelloResp{Message: "Hello " + in.Name}, nil
}
