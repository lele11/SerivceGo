package discovery

import (
	"sync"
	"time"
)

const (
	GetServiceTypeMiniLoad = "MiniLoad"
	GetServiceTypeRandom   = "Random"
)

var poolMgr *ServicePoolMgr

func getServicePoolMgr() *ServicePoolMgr {
	if poolMgr == nil {
		poolMgr = &ServicePoolMgr{
			pool: map[string]*ServicePool{},
		}
	}
	return poolMgr
}

type ServicePoolMgr struct {
	pool map[string]*ServicePool
}

func (spm *ServicePoolMgr) NewServicePool(name string) {
	sp := &ServicePool{
		name: name,
	}
	spm.pool[sp.name] = sp
	sp.FreshData()
	go sp.Run()
}

func (spm *ServicePoolMgr) GetAllList(name string) []*ServiceDesc {
	s := spm.pool[name]
	if s == nil {
		return nil
	}
	return s.list
}

func (spm *ServicePoolMgr) GetService(name string, id string, flag string) *ServiceDesc {
	p, ok := spm.pool[name]
	if !ok {
		return nil
	}
	if id == "MiniLoad" {
		return p.GetMiniLoad(flag)
	}
	if id == "Random" {
		return p.GetOne(flag)
	}
	return p.GetById(id)
}

type ServicePool struct {
	name string //服务的名称
	list []*ServiceDesc
	lock sync.RWMutex
}

func (sp *ServicePool) Run() {
	timer := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-timer.C:
			sp.FreshData()
			//TODO 没有服务时 是否需要退出
		}
	}
}

func (sp *ServicePool) FreshData() {
	data := Default.Query(sp.name)
	sp.lock.Lock()
	sp.list = data
	sp.lock.Unlock()
}
func (sp *ServicePool) isEmpty() bool {
	return len(sp.list) == 0
}
func (sp *ServicePool) GetOne(flag string) *ServiceDesc {
	if sp.isEmpty() {
		return nil
	}
	sp.lock.RLock()
	defer sp.lock.RUnlock()
	for _, l := range sp.list {
		if l.State.Invalid {
			continue
		}
		if f, ok := l.Meta["flag"]; ok && f != flag {
			continue
		}
		return l
	}
	return nil
}
func (sp *ServicePool) GetMiniLoad(flag string) *ServiceDesc {
	if sp.isEmpty() {
		return nil
	}
	sp.lock.RLock()
	defer sp.lock.RUnlock()
	var mini *ServiceDesc
	for _, l := range sp.list {
		if f, ok := l.Meta["flag"]; ok && f != flag {
			continue
		}
		if mini == nil {
			mini = l
		}

		if !l.State.Invalid && mini.GetLoad() > l.GetLoad() {
			mini = l
		}
	}
	return mini
}
func (sp *ServicePool) GetById(id string) *ServiceDesc {
	if sp.isEmpty() {
		return nil
	}
	sp.lock.RLock()
	defer sp.lock.RUnlock()
	for _, l := range sp.list {
		if l.ID == id {
			return l
		}
	}
	return nil
}
