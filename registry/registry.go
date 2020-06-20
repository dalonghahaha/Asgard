package registry

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/dalonghahaha/avenger/components/logger"

	"Asgard/constants"
	"Asgard/runtimes"
)

var r *etcdRegistry

type etcdRegistry struct {
	ttl        int64
	schema     string
	exitSingel chan bool
	ticker     *time.Ticker
	client     *clientv3.Client
}

func RegisterRegistry(Endpoints []string) error {
	config := clientv3.Config{
		Endpoints: Endpoints,
	}
	client, err := clientv3.New(config)
	if err != nil {
		return err
	}
	r = &etcdRegistry{
		ttl:    constants.MASTER_TTL,
		schema: constants.MASTER_SCHEMA,
		client: client,
		ticker: time.NewTicker(time.Second * time.Duration(constants.MASTER_TTL)),
	}
	return nil
}

func Register(name, id, ip, port string) {
	go r.keepAlive(name, id, fmt.Sprintf("%s:%s", ip, port))
}

func UnRegister(name, id string) {
	r.ticker.Stop()
	_, err := r.client.Delete(context.Background(), "/"+r.schema+"/"+name+"/"+id)
	if err != nil {
		logger.Errorf("Etcd Delete failed:%+v", err)
	}
}

func (r *etcdRegistry) keepAlive(name, id, addr string) {
	logger.Debugf("keepAlive ticker start!")
	runtimes.SubscribeExit(r.exitSingel)
	for {
		select {
		case <-r.exitSingel:
			logger.Debugf("keepAlive monitor ticker stop!")
			r.ticker.Stop()
			break
		case <-r.ticker.C:
			getResp, err := r.client.Get(context.Background(), "/"+r.schema+"/"+name+"/"+id)
			if err != nil {
				logger.Errorf("Etcd Get failed:%+v", err)
			} else if getResp.Count == 0 {
				err = r.withAlive(name, id, addr, r.ttl)
				if err != nil {
					logger.Errorf("withAlive failed:%+v", err)
				}
			}
		}
	}
}

func (r *etcdRegistry) withAlive(name, id, addr string, ttl int64) error {
	leaseResp, err := r.client.Grant(context.Background(), ttl)
	if err != nil {
		return err
	}
	_, err = r.client.Put(context.Background(), "/"+r.schema+"/"+name+"/"+id, addr, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}
	_, err = r.client.KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		return err
	}
	return nil
}
