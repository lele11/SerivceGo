package servermgr

import "game/base"

var listByKind = map[uint32]base.IService{}

// Register 服务注册
func Register(s base.IService) {
	listByKind[s.GetKind()] = s
}

// GetService 获取服务实例
func GetService(kind uint32) base.IService {
	return listByKind[kind]
}
