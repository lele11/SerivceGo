package discovery

import (
	"time"
)

// Default 默认对象
var Default Discovery

// Discovery 对上层业务暴露的接口 ，提供查找服务，注册服务，监听服务 ，获取配置数据，修改配置
type Discovery interface {
	QueryServices() []string
	// 注册服务
	Register(*ServiceDesc) error
	// 解注册服务
	Deregister(svcid string) error
	// 查询
	Query(name string) []*ServiceDesc
	// 服务状态上报
	UpdateService(ID, status, info string) error
	// 设置值
	SetValue(key string, value interface{}) error
	// 取值，并赋值到变量
	GetValue(key string, valuePtr interface{}) error
	// 删除值
	DeleteValue(key string) error
}

// WatchServiceAll 监控缓存所有服务信息 TODO 过滤
func WatchServiceAll() {
	services := Default.QueryServices()
	for _, name := range services {
		WatchService(name)
	}
}

// WatchService 监控某个服务类型
func WatchService(serviceName string) {
	getServicePoolMgr().NewServicePool(serviceName)
}

// GetServiceList 获取某个服务类型的服务列表
func GetServiceList(serviceName string) []*ServiceDesc {
	return getServicePoolMgr().GetAllList(serviceName)
}

// GetServiceByID 获取某个服务信息
func GetServiceByID(serviceName string, serviceID string) *ServiceDesc {
	return getServicePoolMgr().GetService(serviceName, serviceID, "")
}

// GetServiceMiniLoad 根据标记获取最小负载的服务
func GetServiceMiniLoad(serviceName string, flag string) *ServiceDesc {
	return getServicePoolMgr().GetService(serviceName, GetServiceTypeMiniLoad, flag)
}

// GetServiceOne 根据标记随机获取一个服务
func GetServiceOne(serviceName string, flag string) *ServiceDesc {
	return getServicePoolMgr().GetService(serviceName, GetServiceTypeRandom, flag)
}

// DeregisterService 注销服务
func DeregisterService(serviceID string) error {
	return Default.Deregister(serviceID)
}

// RegisterService 注册服务
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
