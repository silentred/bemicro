package main

import (
	"bemicro/discovery"
	"bemicro/gateway"
	"bemicro/proto"
	"log"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

const srvName = "greeter"
const etcdHost = "http://localhost:2379"

func main() {
	byhand()
	//byConns()
}

func byConns() {
	// new connContainer
	cliFunc := map[string]interface{}{
		srvName: proto.NewGreeterClient,
	}
	c := gateway.NewConnContainer(cliFunc)

	c.InitConns(etcdHost, []string{srvName}, false)

	req := proto.HelloReq{Name: "jason"}

	for range time.Tick(time.Second / 2) {

		tracePair := gateway.GetTraceIDPair()
		// authPair := gateway.GetAuthInfoPair()
		ctx := gateway.MergeStrings(context.Background(), tracePair)

		resp, err := c.CallRPC(ctx, srvName, "SayHello", &req)
		if err != nil {
			log.Println(err)
		}

		if resp, ok := resp.(*proto.HelloResp); ok {
			log.Printf("Greeting: %s", resp.Message)
		} else {
			log.Printf("not valid resp %#v", resp)
		}
	}
}

func byhand() {
	// new greeter with balancer
	resolver := discovery.NewResolver(srvName, discovery.DefaultPrefix)
	b := grpc.RoundRobin(resolver)

	opt := grpc.WithBalancer(b)

	conn, err := grpc.Dial(etcdHost, grpc.WithInsecure(), opt)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := proto.NewGreeterClient(conn)

	// Contact the server and print out its response.
	traceIDPair := gateway.GetTraceIDPair()
	authPair := gateway.GetAuthInfoPair()

	userToken, err := gateway.GenerateUserToken(1, "jason", "aaa@163.com")
	if err != nil {
		panic(err)
	}
	userPair := gateway.GetUserClaimPair(userToken)

	ctx := gateway.MergeStrings(context.Background(), traceIDPair, authPair, userPair)

	r, err := c.SayHello(ctx, &proto.HelloReq{Name: "jason"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)

	// for i := 0; i < 20; i++ {
	// 	go func(i int) {
	// 		// Contact the server and print out its response.
	// 		ctx := gateway.PlantTraceID(context.Background())
	// 		//ctx = gateway.PlantAuthInfo(ctx)

	// 		r, err := c.SayHello(ctx, &proto.HelloReq{Name: "jason", Times: int32(i)})
	// 		if err != nil {
	// 			log.Fatalf("could not greet: %v", err)
	// 		}
	// 		log.Printf("No. %d, Greeting: %s", i, r.Message)
	// 	}(i)
	// }

	// select {}
}
