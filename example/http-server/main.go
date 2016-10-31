package main

import (
	"bemicro/gateway"
	"bemicro/proto"
	"net/http"
)

func main() {
	srvName := "greeter"
	etcdHost := "http://localhost:2379"
	cliFunc := map[string]interface{}{
		srvName: proto.NewGreeterClient,
	}

	c := gateway.NewConnContainer(cliFunc)
	c.InitConns(etcdHost, []string{srvName}, false)

	server := gateway.NewHTTPServer(c)

	http.ListenAndServe(":8088", server)
}
