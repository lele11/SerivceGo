package serverMgr

import "game/base"

var listByKind = map[uint32]base.IService{}

func Register(s base.IService) {
	listByKind[s.GetKind()] = s
}
func GetService(kind uint32) base.IService {
	return listByKind[kind]
}
