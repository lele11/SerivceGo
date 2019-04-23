package game

import (
	"game/base/redis"
	"game/base/util"
)

func (p *Player) NewPlayerBaseModule() {
	m := &PlayerBaseModule{p: p}
	m.Init()
	p.baseModule = m
	p.RegModule(ModuleBaseInfo, m)
}

const (
	PlayerLogin  = "login"
	PlayerLogout = "logout"
	PlayerCreate = "create" //创角时间
)

type PlayerBaseModule struct {
	login  int64
	logout int64
	create int64
	update map[interface{}]interface{}
	p      *Player
}

func (pbm *PlayerBaseModule) Key() string {
	return "player:" + util.Uint64ToString(pbm.p.uid)
}
func (pbm *PlayerBaseModule) Init() {
	pbm.update = map[interface{}]interface{}{}
}
func (pbm *PlayerBaseModule) Load() {
	if redis.Exists(pbm.Key()) {
		info := redis.HGetAll(pbm.Key())
		pbm.SetLogin(util.StringToInt64(info[PlayerLogin]))
		pbm.SetLogout(util.StringToInt64(info[PlayerLogout]))
		pbm.create = util.StringToInt64(info[PlayerCreate])
	}
}
func (pbm *PlayerBaseModule) Save() {
	if len(pbm.update) > 0 {
		redis.HMSet(pbm.Key(), pbm.update)
		pbm.update = map[interface{}]interface{}{}
	}
}

func (pbm *PlayerBaseModule) SetLogin(login int64) {
	pbm.login = login
	pbm.update[PlayerLogin] = login
}
func (pbm *PlayerBaseModule) SetLogout(logout int64) {
	pbm.logout = logout
	pbm.update[PlayerLogout] = logout
}

func (pbm *PlayerBaseModule) SetCreate(create int64) {
	pbm.create = create
	pbm.update[PlayerCreate] = create
}

func (pm *PlayerManager) HandlerPlayerInfo(p *Player, data []byte) {

}
