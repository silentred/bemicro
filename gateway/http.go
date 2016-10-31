package gateway

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"golang.org/x/net/context"
)

// POST /grpc/{service_name}/{method_name}
// json.Decode(body) to corresponding request type

type HTTPServer struct {
	conns *ConnContainer
	types map[string]reflect.Type
	tLock *sync.RWMutex
}

func NewHTTPServer(conns *ConnContainer) *HTTPServer {
	return &HTTPServer{
		conns: conns,
		types: make(map[string]reflect.Type),
		tLock: &sync.RWMutex{},
	}
}

// TODO return ([]byte, error)
func (s HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Write([]byte(`{"error" : "method should be POST"}`))
		return
	}

	srv, method, grpcReq := s.parseHTTP(r)
	if grpcReq == nil {
		w.Write([]byte(fmt.Sprintf(`{"error" : "srv=%s method=%s req=%v"}`, srv, method, grpcReq)))
		return
	}

	tracePair := GetTraceIDPair()
	ctx := MergeStrings(context.Background(), tracePair)

	ret, err := s.conns.CallRPC(ctx, srv, method, grpcReq)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error" : %s}`, err)))
		return
	}

	b, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error" : %s}`, err)))
		return
	}

	w.Write(b)
}

// /grpc/v1/xx/{service_name}/{method_name}
// return *proto.Request interface
func (s *HTTPServer) parseHTTP(request *http.Request) (srv, method string, req interface{}) {
	ss := strings.Split(request.URL.Path, "/")

	if len(ss) < 4 {
		fmt.Printf("invalid uri %s \n", request.URL.Path)
		return
	}
	srv = ss[len(ss)-2]
	method = ss[len(ss)-1]

	cliVal, err := s.conns.getClient(srv)
	if err != nil {
		// log
		return
	}

	reqType := s.getRequestType(srv, method, cliVal)
	if reqType == nil {
		return
	}

	req = parseJSONData(request.Body, reqType)

	return
}

func (s *HTTPServer) getRequestType(srvName, method string, cliVal *reflect.Value) reflect.Type {
	// check cache
	key := getTypeKey(srvName, method)
	if t := s.getType(key); t != nil {
		return t
	}

	mVal := cliVal.MethodByName(method)
	if mVal.IsValid() {
		t := mVal.Type().In(1)
		s.setType(key, t)
		return t
	}

	return nil
}

func parseJSONData(data io.Reader, t reflect.Type) interface{} {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	p := reflect.New(t).Interface() // ptr interface
	json.NewDecoder(data).Decode(p)

	return p
}

func (s *HTTPServer) getType(k string) reflect.Type {
	s.tLock.RLock()
	defer s.tLock.RUnlock()

	if t, ok := s.types[k]; ok {
		return t
	}
	return nil
}

func (s *HTTPServer) setType(k string, t reflect.Type) {
	s.tLock.Lock()
	defer s.tLock.Unlock()

	if t, ok := s.types[k]; !ok {
		s.types[k] = t
	}
}

func getTypeKey(srv, method string) string {
	return fmt.Sprintf("%s.%s", srv, method)
}
