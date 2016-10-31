package gateway

import (
	"bemicro/discovery"
	"fmt"
	"log"
	"reflect"
	"sync"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

// ConnContainer is a thread safe conn container. Key is service name, value is grpc.ClientConn
type ConnContainer struct {
	// srvName => connection
	conns map[string]*grpc.ClientConn
	lock  *sync.RWMutex

	// srvName.method => proto.XXXClient
	clients map[string]*reflect.Value
	cLock   *sync.RWMutex

	// srvName => newClientFunc
	cliFunc map[string]interface{}
}

// NewConnContainer return a new thread safe ConnContainer
func NewConnContainer(cliFuncMap map[string]interface{}) *ConnContainer {
	return &ConnContainer{
		lock:    new(sync.RWMutex),
		conns:   make(map[string]*grpc.ClientConn),
		cLock:   new(sync.RWMutex),
		clients: make(map[string]*reflect.Value),
		cliFunc: cliFuncMap,
	}
}

func (m *ConnContainer) getClient(k string) (*reflect.Value, error) {
	m.cLock.RLock()

	if val, ok := m.clients[k]; ok {
		m.cLock.RUnlock()
		return val, nil
	}
	m.cLock.RUnlock()

	conn := m.Get(k)
	if conn == nil {
		return nil, fmt.Errorf("service conn '%s' not found", k)
	}
	newCliFunc := m.getCliFunc(k)
	// get NewServiceClient's reflect.Value
	vClient := reflect.ValueOf(newCliFunc)
	var vParameter []reflect.Value
	vParameter = append(vParameter, reflect.ValueOf(conn))

	// c[0] is serviceServer reflect.Value
	c := vClient.Call(vParameter)

	m.setClient(k, &c[0])

	return &c[0], nil

}

func (m *ConnContainer) setClient(k string, v *reflect.Value) {
	m.cLock.Lock()
	if _, ok := m.clients[k]; !ok {
		m.clients[k] = v
	}
	m.cLock.Unlock()
}

func (m *ConnContainer) getCliFunc(k string) interface{} {
	m.cLock.RLock()
	defer m.cLock.RUnlock()
	if val, ok := m.cliFunc[k]; ok {
		return val
	}

	return nil
}

// Get from maps return the k's value
func (m *ConnContainer) Get(k string) *grpc.ClientConn {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if val, ok := m.conns[k]; ok {
		return val
	}
	return nil
}

// Set Maps the given key and value. Returns false if the key is already in the map and changes nothing.
func (m *ConnContainer) Set(k string, v *grpc.ClientConn) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	if val, ok := m.conns[k]; !ok {
		m.conns[k] = v
	} else if val != v {
		m.conns[k] = v
	} else {
		return false
	}
	return true
}

// SetService sets grpc.ClientConn to container by serivce name
func (m *ConnContainer) SetService(etcdAddrs, srvName string) error {
	r := discovery.NewResolver(srvName, discovery.DefaultPrefix)
	b := grpc.RoundRobin(r)

	conn, err := grpc.Dial(etcdAddrs, grpc.WithInsecure(), grpc.WithBalancer(b))
	if err != nil {
		log.Printf(`connect to '%s' service failed: %v`, srvName, err)
		return err
	}

	if m.Set(srvName, conn) {
		return nil
	}

	return fmt.Errorf("service %s exists, Set failed", srvName)
}

// Exists returns true if k is exist in the map.
func (m *ConnContainer) Exists(k string) bool {
	m.lock.RLock()
	_, ok := m.conns[k]
	m.lock.RUnlock()

	return ok
}

// Delete removes key
func (m *ConnContainer) Delete(k string) {
	m.lock.Lock()
	if conn, ok := m.conns[k]; ok {
		conn.Close()
		delete(m.conns, k)
	}

	m.lock.Unlock()
}

// All returns full map of connections
func (m *ConnContainer) All() map[string]*grpc.ClientConn {
	return m.conns
}

// InitConns initialize service connections
func (m *ConnContainer) InitConns(etcdAddrs string, serviceList []string, force bool) {
	wg := sync.WaitGroup{}
	wg.Add(len(serviceList))

	for _, serviceName := range serviceList {
		go func(srvName string) {
			defer wg.Done()

			if !force && m.Exists(srvName) {
				return
			}

			err := m.SetService(etcdAddrs, srvName)
			if err != nil {
				log.Println(err)
			}
		}(serviceName)
	}

	wg.Wait()
}

// CloseAll connections
func (m *ConnContainer) CloseAll() {
	for key := range m.All() {
		m.Delete(key)
	}
}

// CallRPC is helper func that make life easier
// client: grpc client Constructor function
func (m *ConnContainer) CallRPC(ctx context.Context, serviceName string, method string, req interface{}) (ret interface{}, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("call RPC '%s' error: %v", method, x)
		}
	}()

	cli, err := m.getClient(serviceName)
	if err != nil {
		fmt.Println(err)
	}

	// rpc param
	v := make([]reflect.Value, 2)
	v[0] = reflect.ValueOf(ctx)
	v[1] = reflect.ValueOf(req)
	// rpc method call
	f := cli.MethodByName(method)
	resp := f.Call(v)
	if !resp[1].IsNil() {
		return nil, resp[1].Interface().(error)
	}
	return resp[0].Interface(), nil
}
