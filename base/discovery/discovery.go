package discovery

import (
	"time"
)

var Default Discovery

// 对上层业务暴露的接口 ，提供查找服务，注册服务，监听服务 ，获取配置数据，修改配置
type Discovery interface {
	QueryServices() []string
	// 注册服务
	Register(*ServiceDesc) error
	// 解注册服务
	Deregister(svcid string) error
	Query(name string) []*ServiceDesc
	UpdateService(ID, status, info string) error
	// 设置值
	SetValue(key string, value interface{}) error
	// 取值，并赋值到变量
	GetValue(key string, valuePtr interface{}) error
	// 删除值
	DeleteValue(key string) error
}

func WatchServiceAll() {
	services := Default.QueryServices()
	for _, name := range services {
		WatchService(name)
	}
}

func WatchService(serviceName string) {
	getServicePoolMgr().NewServicePool(serviceName)
}

func GetServiceList(serviceName string) []*ServiceDesc {
	return getServicePoolMgr().GetAllList(serviceName)
}

func GetServiceById(serviceName string, serviceID string) *ServiceDesc {
	return getServicePoolMgr().GetService(serviceName, serviceID, "")
}
func GetServiceMiniLoad(serviceName string, flag string) *ServiceDesc {
	return getServicePoolMgr().GetService(serviceName, GetServiceTypeMiniLoad, flag)
}
func GetServiceOne(serviceName string, flag string) *ServiceDesc {
	return getServicePoolMgr().GetService(serviceName, GetServiceTypeRandom, flag)
}
func DeregisterService(serviceID string) error {
	return Default.Deregister(serviceID)
}
func RegisterService(desc *ServiceDesc) error {
	err := Default.Register(desc)
	if err != nil {
		return err
	}
	ls := &LocalService{
		name: desc.ID,
		ttl:  5 * time.Second,
		f:    desc.Reporter,
	}
	localServices.Store(ls.name, ls)
	go ls.update()
	return nil
}
