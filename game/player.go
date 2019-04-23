package game

import (
	"game/base"
	"game/config"
	"game/db"
	"game/protoMsg"
	"time"
)

type Player struct {
	uid          uint64
	gateId       uint64
	lastSave     time.Time
	loginChannel string
	pm           *PlayerManager
	modules      map[int]Module
	baseModule   *PlayerBaseModule
	resChange    *ResChange
	tmpPkRank    uint32
	tmpPKList    []uint32
	online       bool
}

func (p *Player) Init() {
	p.modules = map[int]Module{}
	p.NewPlayerBaseModule()
	p.newResChange()
	for _, m := range p.modules {
		m.Init()
	}
}
func (p *Player) Load() {
	for _, m := range p.modules {
		m.Load()
	}
	p.lastSave = time.Now().Add(10 * time.Second)
}

func (p *Player) loop() {

}

func (p *Player) Save(force bool) {
	//TODO 随机时间
	if !p.online {
		return
	}
	if time.Now().After(p.lastSave) || force {
		for _, m := range p.modules {
			m.Save()
		}
		p.lastSave = time.Now().Add(10 * time.Second)
	}
}
func (p *Player) RegModule(kind int, m Module) {
	p.modules[kind] = m
}
func (p *Player) SetGateId(i uint64) {
	p.gateId = i
}
func (p *Player) SendToClient(cmd protoMsg.S_CMD, msg base.IMsg) {
	p.pm.s.SendMsg(config.ServerKindGateway, p.gateId, p.uid, int32(cmd), msg)
}
func (p *Player) create() {
	if p.baseModule.create > 0 {
		return
	}

	p.baseModule.SetCreate(time.Now().Unix())
}

func (p *Player) login() {
	p.create()
	p.loginChannel = db.GetPlayerBaseInfo(p.uid, db.ChannelField)
	p.baseModule.SetLogin(time.Now().Unix())
}
func (p *Player) logout() {
	p.baseModule.SetLogout(time.Now().Unix())
}
