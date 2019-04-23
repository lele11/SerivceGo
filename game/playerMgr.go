package game

import (
	"game/protoMsg"

	"github.com/cihub/seelog"
)

func NewPlayerManager(s *GameServer) *PlayerManager {
	p := &PlayerManager{
		s:       s,
		handler: map[int32]func(*Player, []byte){},
		player:  map[uint64]*Player{},
	}
	p.AttachPlayerHandler()
	return p
}

type PlayerManager struct {
	s       *GameServer
	handler map[int32]func(*Player, []byte)
	player  map[uint64]*Player
}

func (pm *PlayerManager) playerLogin(uid uint64, gateID uint64) {
	p := pm.getPlayer(uid)
	if p == nil {
		p = &Player{
			uid: uid,
			pm:  pm,
		}
		p.Init()
	}
	p.Load()
	p.SetGateId(gateID)
	p.login()
	p.online = true
	pm.player[p.uid] = p
	seelog.Infof("Player %d Login From %d ", uid, gateID)
}
func (pm *PlayerManager) playerLogout(uid uint64) {
	p := pm.getPlayer(uid)
	if p == nil {
		return
	}
	p.logout()
	p.Save(true)
	p.online = false
	delete(pm.player, uid)
	seelog.Infof("Player %d Logout ", uid)
}
func (pm *PlayerManager) getPlayer(uid uint64) *Player {
	return pm.player[uid]
}
func (pm *PlayerManager) dispatcher(uid uint64, cmd int32, data []byte) {
	player := pm.getPlayer(uid)
	if player == nil {
		seelog.Error("Not Found Player ", uid)
		return
	}
	seelog.Debug("Recv Msg ", cmd)
	if f, ok := pm.handler[cmd]; ok {
		f(player, data)
		player.ResChange()
	} else {
		seelog.Errorf("Not Found Handler UserPacket uid %d  cmd %d", uid, cmd)
	}
}
func (pm *PlayerManager) AttachHandler(cmd protoMsg.C_CMD, f func(*Player, []byte)) {
	pm.handler[int32(cmd)] = f
}
func (pm *PlayerManager) update() {
	for _, p := range pm.player {
		p.Save(false)
		p.loop()
	}
}
