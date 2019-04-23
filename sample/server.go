package sample

import (
	"server/base"
	"server/base/packet"
	"server/base/proto"
	"server/base/serverMgr"
	"server/config"

	"github.com/cihub/seelog"
)

func init() {
	serverMgr.Register(&SampleServer{})
}

// 基础服务结构
type SampleServer struct {
	*base.Service
}

func (game *SampleServer) GetKind() uint32 {
	return config.ServerKindGame
}
func (game *SampleServer) Init(cfg *config.ServerConfig) {
	game.Service = base.NewService(cfg.ID, cfg.Host, cfg.Port, cfg.Kind, cfg.Protocol)
	game.SetReporter(game.reporter)
	game.SetUpdate(game.Loop)

}

func (game *SampleServer) Run() {
	game.AttachHandler(uint16(innerMsg.InnerCmd_transport), game.HandlerPacketTransport)
	game.AttachHandler(uint16(innerMsg.InnerCmd_clientConnect), game.HandlerClientConnect)
	game.Service.Run()
}

// 服务器主循环函数
func (game *SampleServer) Loop() {

}

// 客户端业务消息处理
func (game *SampleServer) HandlerPacketTransport(packet packet.IPacket) {
	msg := &innerMsg.PacketTransport{}
	if e := msg.Unmarshal(packet.GetBody()); e != nil {
		return
	}
	seelog.Errorf("recv transport  ", msg)
}

// 客户端连接消息通知 客户端接入，客户端断开
func (game *SampleServer) HandlerClientConnect(packet packet.IPacket) {
	msg := &innerMsg.ClientConnect{}
	if e := msg.Unmarshal(packet.GetBody()); e != nil {
		return
	}
	seelog.Errorf("recv client connect ", msg)
}
func (game *SampleServer) nodeUpdate() *innerMsg.ServerUpdate {
	return &innerMsg.ServerUpdate{
		SID: game.GetSID(),
	}
}
func (game *SampleServer) reporter() (output string, status string) {
	output = "Game Ok"
	status = "passing"
	return
}
