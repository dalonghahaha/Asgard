package registry

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/dalonghahaha/avenger/components/logger"

	"Asgard/constants"
	"Asgard/runtimes"
)

var (
	leaderFlag bool
	client     *clientv3.Client
	r          *etcdRegistry
)

type etcdRegistry struct {
	ttl        int64
	schema     string
	exitSingel chan bool
	ticker     *time.Ticker
}

func IsLeader() bool {
	return leaderFlag
}

func RegisterRegistry(Endpoints []string) error {
	config := clientv3.Config{
		Endpoints: Endpoints,
	}
	var err error
	client, err = clientv3.New(config)
	if err != nil {
		return err
	}
	r = &etcdRegistry{
		ttl:    constants.MASTER_CLUSTER_TTL,
		schema: constants.MASTER_CLUSTER_SCHEMA,
		ticker: time.NewTicker(time.Second * time.Duration(constants.MASTER_CLUSTER_TTL)),
	}
	return nil
}

func Register(name, id, ip, port string) {
	go keepAlive(name, id, fmt.Sprintf("%s:%s", ip, port))
}

func UnRegister(name, id string) {
	r.ticker.Stop()
	_, err := client.Delete(context.Background(), "/"+r.schema+"/"+name+"/"+id)
	if err != nil {
		logger.Errorf("Etcd Delete failed:%+v", err)
	}
}

func keepAlive(name, id, addr string) {
	logger.Debugf("keepAlive ticker start!")
	runtimes.SubscribeExit(r.exitSingel)
	for {
		select {
		case <-r.exitSingel:
			logger.Debugf("keepAlive monitor ticker stop!")
			r.ticker.Stop()
			break
		case <-r.ticker.C:
			getResp, err := client.Get(context.Background(), "/"+r.schema+"/"+name+"/"+id)
			if err != nil {
				logger.Errorf("Etcd Get failed:%+v", err)
			} else if getResp.Count == 0 {
				err = withAlive(name, id, addr, r.ttl)
				if err != nil {
					logger.Errorf("withAlive failed:%+v", err)
				}
			}
		}
	}
}

func withAlive(name, id, addr string, ttl int64) error {
	leaseResp, err := client.Grant(context.Background(), ttl)
	if err != nil {
		return err
	}
	_, err = client.Put(context.Background(), "/"+r.schema+"/"+name+"/"+id, addr, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}
	_, err = client.KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		return err
	}
	return nil
}

func Campaign(pfx string, prop string) {
	for {
		session, err := concurrency.NewSession(client, concurrency.WithTTL(constants.MASTER_CLUSTER_CAMPAIGN_TTL))
		if err != nil {
			logger.Errorf("campaign new session failed:%+v", err)
			continue
		}
		election := concurrency.NewElection(session, pfx)
		if err = election.Campaign(context.Background(), prop); err != nil {
			logger.Errorf("campaign election failed:%+v", err)
			continue
		}
		logger.Info("campaign leader sccuess!")
		leaderFlag = true
		select {
		case <-session.Done():
			leaderFlag = false
		}
	}
}
