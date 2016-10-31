package gateway

import (
	"bemicro/proto"
	"bmw/lib"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestHttp(t *testing.T) {
	srvName := "greeter"
	etcdHost := "http://localhost:2379"
	cliFunc := map[string]interface{}{
		srvName: proto.NewGreeterClient,
	}

	c := NewConnContainer(cliFunc)
	c.InitConns(etcdHost, []string{srvName}, false)

	server := NewHTTPServer(c)

	w := &mockResponseWriter{}
	httpReq, _ := lib.NewHTTPReqeust("POST", "http://localhost:1234/grpc/greeter/SayHello", nil, nil, []byte(`{"name": "Tony", "times": 1}`))

	server.ServeHTTP(w, httpReq)

}

type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	fmt.Println(string(p))
	return len(p), nil
}

func (m *mockResponseWriter) WriteHeader(int) {}

// ========= test reflect =========

type Helper struct {
}

type Req struct {
	Name string `json:"name"`
}

func (h *Helper) Help(req *Req) {
	fmt.Println("helping", req.Name)
}

func TestReflect(t *testing.T) {
	h := &Helper{}
	hVal := reflect.ValueOf(h)
	fmt.Println(hVal.Kind(), hVal.NumMethod())
	mVal := hVal.MethodByName("Help")
	if !mVal.IsValid() {
		panic("mVal is invalid")
	}

	pType := mVal.Type().In(0)
	fmt.Println(pType.Kind())

	reader := bytes.NewBufferString(`{"name": "jason"}`)
	p := parseJSONData(reader, pType)

	pVal := reflect.ValueOf(p)
	mVal.Call([]reflect.Value{pVal})

}

func TestJson(t *testing.T) {
	v := map[int]string{}
	//v[1] = "test"
	//b, err := json.Marshal(v)

	b := []byte(`{"0": "value"}`)
	err := json.Unmarshal(b, &v)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(v)
}
