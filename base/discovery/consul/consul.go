package consulsd

import (
	"game/base/discovery"
	"sync"
	"time"

	"github.com/cihub/seelog"
	"github.com/hashicorp/consul/api"
)

type consulDiscovery struct {
	consulAddr string
	client     *api.Client
	kvCache    sync.Map // map[string]*KVPair
}

//  WaitReady 等待consul服务启动成功
func (self *consulDiscovery) WaitReady() {
	for {
		_, _, err := self.client.Health().Service("consul", "", false, nil)
		if err == nil {
			break
		}
		seelog.Error(err)
		time.Sleep(time.Second * 2)
	}
}

func NewDiscovery(address string, kvPrefix string) discovery.Discovery {
	if address == "" {
		address = "127.0.0.1:8500"
	}
	self := &consulDiscovery{
		consulAddr: address,
	}
	c := api.DefaultConfig()
	c.Address = address
	var err error
	self.client, err = api.NewClient(c)

	if err != nil {
		panic(err)
	}
	self.WaitReady()
	self.startWatchKV(kvPrefix)

	return self
}
