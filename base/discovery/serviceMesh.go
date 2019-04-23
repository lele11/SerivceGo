package discovery

import (
	"sync"
	"time"
)

const (
	// GetServiceTypeMiniLoad 按照最低负载
	GetServiceTypeMiniLoad = "MiniLoad"
	// GetServiceTypeRandom 随机
	GetServiceTypeRandom = "Random"
)

var poolMgr *ServicePoolMgr

// getServicePoolMgr 获取某个服务池管理器
func getServicePoolMgr() *ServicePoolMgr {
	if poolMgr == nil {
		poolMgr = &ServicePoolMgr{
			pool: map[string]*ServicePool{},
		}
	}
	return poolMgr
}

// ServicePoolMgr 服务池管理器
type ServicePoolMgr struct {
	pool map[string]*ServicePool
}

// NewServicePool 创建某个服务的服务池
func (spm *ServicePoolMgr) NewServicePool(name string) {
	sp := &ServicePool{
		name: name,
	}
	spm.pool[sp.name] = sp
	sp.freshData()
	go sp.Run()
}

// GetAllList 获取某个服务的所有节点列表
func (spm *ServicePoolMgr) GetAllList(name string) []*ServiceDesc {
	s := spm.pool[name]
	if s == nil {
		return nil
	}
	return s.list
}

// GetService 获取某个服务的一个节点
func (spm *ServicePoolMgr) GetService(name string, id string, flag string) *ServiceDesc {
	p, ok := spm.pool[name]
	if !ok {
		return nil
	}
	if id == "MiniLoad" {
		return p.getMiniLoad(flag)
	}
	if id == "Random" {
		return p.getOne(flag)
	}
	return p.getById(id)
}

// ServicePool 服务池 缓存某个服务类型的所有节点信息
type ServicePool struct {
	name string //服务的名称
	list []*ServiceDesc
	lock sync.RWMutex
}

// Run  运行服务池逻辑
func (sp *ServicePool) Run() {
	timer := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-timer.C:
			sp.freshData()
			//TODO 没有服务时 是否需要退出
		}
	}
}

// freshData 刷新
func (sp *ServicePool) freshData() {
	data := Default.Query(sp.name)
	sp.lock.Lock()
	sp.list = data
	sp.lock.Unlock()
}

// isEmpty 是否为空
func (sp *ServicePool) isEmpty() bool {
	return len(sp.list) == 0
}

// getOne 获取一个节点
func (sp *ServicePool) getOne(flag string) *ServiceDesc {
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

// getMiniLoad 获取最小负载
func (sp *ServicePool) getMiniLoad(flag string) *ServiceDesc {
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

// getById 根据id 获取
func (sp *ServicePool) getById(id string) *ServiceDesc {
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
