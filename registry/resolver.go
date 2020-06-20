package registry

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/dalonghahaha/avenger/components/logger"
	"google.golang.org/grpc/resolver"

	"Asgard/constants"
)

type etcdResolver struct {
	lock     sync.Mutex
	addrMap  map[string]string
	addrList []resolver.Address
	client   *clientv3.Client
	conn     resolver.ClientConn
}

func NewResolver(Endpoints []string) (resolver.Builder, error) {
	config := clientv3.Config{
		Endpoints:   Endpoints,
		DialTimeout: 15 * time.Second,
	}
	client, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}
	return &etcdResolver{
		addrMap:  map[string]string{},
		addrList: []resolver.Address{},
		client:   client,
	}, nil
}

func (r *etcdResolver) Build(target resolver.Target, conn resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r.conn = conn
	keyPrefix := "/" + target.Scheme + "/" + target.Endpoint + "/"
	result, err := r.client.Get(context.Background(), keyPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	} else {
		addrList := []resolver.Address{}
		for i := range result.Kvs {
			key := string(result.Kvs[i].Key)
			value := string(result.Kvs[i].Value)
			r.addrMap[key] = value
			addrList = append(addrList, resolver.Address{Addr: value})
		}
		r.addrList = addrList
	}
	if len(r.addrMap) == 0 {
		return nil, fmt.Errorf("address list empty")
	}
	logger.Debugf("resolver addrMap:%+v", r.addrMap)
	r.conn.UpdateState(resolver.State{Addresses: r.addrList})
	go r.watch(keyPrefix)
	return r, nil
}

func (r *etcdResolver) Scheme() string {
	return constants.MASTER_SCHEMA
}

func (r *etcdResolver) ResolveNow(rn resolver.ResolveNowOptions) {

}

func (r *etcdResolver) Close() {

}

func (r *etcdResolver) watch(keyPrefix string) {
	rch := r.client.Watch(context.Background(), keyPrefix, clientv3.WithPrefix())
	for n := range rch {
		for _, event := range n.Events {
			key := string(event.Kv.Key)
			value := string(event.Kv.Value)
			switch {
			case event.IsCreate():
				_, ok := r.addrMap[key]
				if !ok {
					logger.Debugf("resolver add address:%s", value)
					r.lock.Lock()
					r.addrMap[key] = value
					r.addrList = append(r.addrList, resolver.Address{Addr: value})
					r.lock.Unlock()
					r.conn.UpdateState(resolver.State{Addresses: r.addrList})
				}
			case event.Type == clientv3.EventTypeDelete:
				_, ok := r.addrMap[key]
				if ok {
					logger.Debugf("resolver remove address:%s", r.addrMap[key])
					r.lock.Lock()
					delete(r.addrMap, key)
					addrList := []resolver.Address{}
					for _, value := range r.addrMap {
						addrList = append(addrList, resolver.Address{Addr: value})
					}
					r.addrList = addrList
					r.lock.Unlock()
					r.conn.UpdateState(resolver.State{Addresses: r.addrList})
				}
			}
		}
	}
}
