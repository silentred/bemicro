package discovery

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/coreos/etcd/client"
)

var (
	DefaultPrefix string = "being/services"
)

// EtcdPublisher publish sevice info to etcd
type EtcdPublisher struct {
	// should be in format as being/workers/{ServiceName}
	Prefix string
	TTL    time.Duration
	Client *client.Client
	Kapi   client.KeysAPI
}

// NewEtcdPublisher returns the publisher which refresh every ttl seconds
func NewEtcdPublisher(hosts []string, ttl int) *EtcdPublisher {
	return NewEtcdPublisherWithPrefix(hosts, DefaultPrefix, ttl)
}

func NewEtcdPublisherWithPrefix(hosts []string, prefix string, ttl int) *EtcdPublisher {
	cfg := client.Config{
		Endpoints:               hosts,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	cli, err := client.New(cfg)
	if err != nil {
		panic(err)
	}

	kapi := client.NewKeysAPI(cli)

	if len(prefix) == 0 {
		prefix = DefaultPrefix
	}

	return &EtcdPublisher{
		Prefix: prefix,
		TTL:    time.Duration(ttl) * time.Second,
		Client: &cli,
		Kapi:   kapi,
	}
}

// Register stores the info of service at registry, and keep it, refresh it
func (ep *EtcdPublisher) Register(service *Service) error {
	if service.ID == 0 {
		var prevVal, value string
		var currentID int
		for !ep.saveIDIndex(service.Name, prevVal, value) {
			if id, err := ep.lookupID(service.Name); err == nil {
				prevVal = strconv.Itoa(id)
				currentID = id + 1
				value = strconv.Itoa(currentID)
			} else {
				currentID = 1
				value = "1"
			}
		}
		service.ID = uint64(currentID)
	}

	path := ep.getFullPath(service)

	opt := &client.SetOptions{TTL: ep.TTL}
	if service.lastIndex > 0 {
		opt.PrevIndex = service.lastIndex
	}

	resp, err := ep.Kapi.Set(context.Background(), path, service.String(), opt)
	if err != nil {
		log.Println(err)
		return err
	}
	service.SetIndex(resp.Index)

	return nil
}

// Unregister removes the Publisher.FullKey at registry
func (ep *EtcdPublisher) Unregister(service *Service) error {
	path := ep.getFullPath(service)

	opt := &client.DeleteOptions{PrevIndex: service.GetIndex()}
	resp, err := ep.Kapi.Delete(context.Background(), path, opt)
	if err != nil {
		log.Println(err)
		return err
	}

	service.SetIndex(resp.Index)
	service.Stop()

	return nil
}

// Heartbeat blocks and refresh TTL every {ttl} seconds until the service is Unregistered
func (ep *EtcdPublisher) Heartbeat(service *Service) {
	ticker := time.NewTicker(ep.TTL / 2)
	path := ep.getFullPath(service)
	opt := &client.SetOptions{
		Refresh:   true,
		PrevExist: client.PrevExist,
		TTL:       ep.TTL,
	}

	for range ticker.C {
		select {
		case <-service.quit:
			return
		default:
			resp, err := ep.Kapi.Set(context.Background(), path, "", opt)
			if err != nil {
				log.Println(err)
			} else {
				if resp.Index > 0 {
					service.SetIndex(resp.Index)
				}
			}
		}
	}
}

func (ep *EtcdPublisher) getFullPath(service *Service) string {
	return fmt.Sprintf("%s/%s/%d", ep.Prefix, service.Name, service.ID)
}

func (ep *EtcdPublisher) getIDKey(srvName string) string {
	return fmt.Sprintf("%s/%s/ID", ep.Prefix, srvName)
}

func getDir(prefix, srvName string) string {
	return fmt.Sprintf("%s/%s", prefix, srvName)
}

func (ep *EtcdPublisher) lookupID(srvName string) (id int, err error) {
	if len(srvName) == 0 {
		panic("srvName cannot be empty")
	}

	var resp *client.Response
	opt := &client.GetOptions{}
	resp, err = ep.Kapi.Get(context.Background(), ep.getIDKey(srvName), opt)
	if err != nil {
		log.Println(err.Error())
		// if e, ok := err.(client.Error); ok && e.Code == client.ErrorCodeKeyNotFound {
		// 	log.Println("key not found")
		// }
		return
	}

	id, err = strconv.Atoi(resp.Node.Value)
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (ep *EtcdPublisher) saveIDIndex(srvName, prevValue, value string) bool {
	if len(value) == 0 {
		return false
	}

	var opt *client.SetOptions

	if len(prevValue) > 0 {
		opt = &client.SetOptions{PrevValue: prevValue}
	} else {
		opt = &client.SetOptions{PrevExist: client.PrevNoExist}
	}

	_, err := ep.Kapi.Set(context.Background(), ep.getIDKey(srvName), value, opt)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
