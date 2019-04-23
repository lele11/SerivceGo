package game

import (
	"game/base"
	"game/base/discovery"
	"game/base/packet"
	"game/base/proto"
	"game/base/serverMgr"
	"game/config"
	"time"

	"github.com/cihub/seelog"
)

func init() {
	servermgr.Register(&GameServer{})
}

type GameServer struct {
	*base.Service              // 基础service
	secondTicker  *time.Ticker //时间控制器
	playerMgr     *PlayerManager
}

func (game *GameServer) GetKind() uint32 {
	return config.ServerKindGame
}

func (game *GameServer) Init(cfg *config.ServerConfig) {
	game.Service = base.NewService(cfg.ID, cfg.Host, cfg.Port, config.ServerKindGame, cfg.Protocol, cfg)
	game.SetClose(game.close)   //设置关闭函数
	game.SetUpdate(game.Update) //设置本服务的循环函数
	game.SetReporter(game.nodeUpdate)
	game.InitHandler()
	game.secondTicker = time.NewTicker(1 * time.Second)

	// 业务代码
	game.playerMgr = NewPlayerManager(game)
}

func (game *GameServer) Run() {
	game.Service.Run()
}

func (game *GameServer) InitHandler() {
	game.AttachHandler(uint16(innerMsg.InnerCmd_transport), game.HandlerUserPacket)
	game.AttachHandler(uint16(innerMsg.InnerCmd_clientConnect), game.HandlerPlayerConnect)
}

func (game *GameServer) HandlerPlayerConnect(packet packet.IPacket) {
	msg := &innerMsg.ClientConnect{}
	if e := msg.Unmarshal(packet.GetBody()); e != nil {
		seelog.Error("Player Connect Error ", e)
	}
	if msg.GetState() == innerMsg.ConnectState_alive {
		game.playerMgr.playerLogin(msg.Id, packet.GetOrigin())
	}
	if msg.GetState() == innerMsg.ConnectState_dead {
		game.playerMgr.playerLogout(msg.Id)
	}

}
func (game *GameServer) HandlerUserPacket(packet packet.IPacket) {
	msg := &innerMsg.PacketTransport{}
	if e := msg.Unmarshal(packet.GetBody()); e != nil {
		seelog.Error("UnMarshal Error ", e)
		return
	}
	game.playerMgr.dispatcher(msg.GetTarget(), msg.GetCmd(), msg.GetData())
}

func (game *GameServer) Update() {
	select {
	case <-game.secondTicker.C:
		game.playerMgr.update()
	default:

	}
}
func (game *GameServer) SendMsg(sType uint32, serverId uint64, uid uint64, cmd int32, msg base.IMsg) {
	userPacket := &innerMsg.PacketTransport{}
	userPacket.Cmd = cmd
	userPacket.Target = uid
	if msg != nil {
		var data []byte
		var e error
		data, e = msg.Marshal()
		if e != nil {
			seelog.Errorf("Send cmd %d Msg %v Error %s", cmd, msg, e)
			return
		}
		userPacket.Data = data
	}
	seelog.Debugf(" SendMsg to %d cmd %d Msg %v", uid, cmd, msg)
	game.SendToService(sType, serverId, uint16(innerMsg.InnerCmd_transport), userPacket)
}

// nodeUpdate TODO 负责函数
func (game *GameServer) nodeUpdate() (output *discovery.ServiceState, status string) {
	output, status = game.ServiceReport()
	return
}

func (game *GameServer) close() {
	//TODO 清理玩家数据
}
