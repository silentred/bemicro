package main

import (
	"bmw/lib"
	"fmt"
)

func main() {

	config := lib.NewReqeustConfig(nil, nil, 5, []byte(`{"name": "Tony", "times": 1}`), nil)
	ret, err := lib.HTTPPost("http://localhost:8088/grpc/greeter/SayHello", config)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(ret))
}
