package gateway

import (
	"errors"
	"game/base"
	"game/base/discovery"
	"game/base/packet"
	"game/base/proto"
	"game/base/serverMgr"
	"game/config"
	"game/db"
	"game/protoMsg"
	"time"

	"github.com/cihub/seelog"
)

// init 服务注册
func init() {
	serverMgr.Register(&GatewayServer{})
}

// GatewayServer 网关服务
type GatewayServer struct {
	*base.Service
	clientTransport map[uint64]uint64 //消息转发映射表 可以扩展不同层次的转发 如 按消息号区别目标服务类型
}

// GetKind 获取服务类型
func (gate *GatewayServer) GetKind() uint32 {
	return config.ServerKindGateway
}

// Init 初始化函数
func (gate *GatewayServer) Init(cfg *config.ServerConfig) {
	gate.clientTransport = map[uint64]uint64{}
	gate.Service = base.NewService(cfg.ID, cfg.Host, cfg.Port, cfg.Kind, cfg.Protocol, cfg)
	gate.SetSessionVerifyHandler(gate.HandlerSessionVerify)
	gate.SetDefaultHandler(gate.HandlerTransformPacket)
	gate.SetReporter(gate.nodeUpdate)
	gate.AttachHandler(uint16(innerMsg.InnerCmd_transport), gate.HandlerTransformPacketT)
	gate.AttachHandler(uint16(innerMsg.InnerCmd_clientConnect), gate.HandlerClientConnect)
	gate.AttachHandler(uint16(innerMsg.InnerCmd_closeSession), gate.HandlerCloseSession)
	gate.AttachHandler(uint16(innerMsg.InnerCmd_broadMsg), gate.HandlerBroadMsg)
}

// Run 运行函数
func (gate *GatewayServer) Run() {
	gate.Service.Run()
}

// HandlerTransformPacketT 收到的外发消息
func (gate *GatewayServer) HandlerTransformPacketT(packet packet.IPacket) {
	msg := &innerMsg.PacketTransport{}
	if e := msg.Unmarshal(packet.GetBody()); e != nil {
		return
	}
	gate.SendData(msg.GetTarget(), uint16(msg.GetCmd()), msg.GetData())
}

// HandlerTransformPacket 收到的外部消息
func (gate *GatewayServer) HandlerTransformPacket(packet packet.IPacket) {
	if packet.GetCmd() == uint16(protoMsg.C_CMD_C_GETSERVERTIME) {
		gate.SendMsg(packet.GetOrigin(), uint16(protoMsg.S_CMD_S_GETSERVERTIME), &protoMsg.S_GetServerTime{
			ServerTime: time.Now().Unix(),
		})
		return
	}
	if packet.GetCmd() == 0 {
		return
	}
	msg := &innerMsg.PacketTransport{}
	msg.Target = packet.GetOrigin()
	msg.Cmd = int32(packet.GetCmd())
	msg.Data = packet.GetBody()
	/*
		TODO 转发路由可以增加多种维度和条件
		1、通过消息号，获取目标类型映射表
		2、基于客户端id，获取节点id
	*/

	srvId := gate.clientTransport[packet.GetOrigin()]
	if srvId != 0 {
		gate.SendToService(config.ServerKindGame, srvId, uint16(innerMsg.InnerCmd_transport), msg)
	}
}

// 某个服务节点的连接断掉，表示远端服务异常
func (gate *GatewayServer) onServerClose(id uint64) {
	for _, p := range gate.clientTransport {
		if p == id {
			//服务器异常关闭
			gate.CloseSession(p, true)
		}
	}
}

// HandlerCloseSession 业务层直接关闭session ，来源于重复登录，强制下线等操作
func (gate *GatewayServer) HandlerCloseSession(packet packet.IPacket) {
	msg := &innerMsg.CloseSession{}
	if e := msg.Unmarshal(packet.GetBody()); e != nil {
		return
	}
	delete(gate.clientTransport, msg.Id)
	gate.CloseSession(msg.Id, false)
}

// HandlerClientConnect 底层连接状态变化通知，包括连接创建和销毁
func (gate *GatewayServer) HandlerClientConnect(packet packet.IPacket) {
	msg := &innerMsg.ClientConnect{}
	if msg.Unmarshal(packet.GetBody()) != nil {
		return
	}
	if msg.GetKind() == innerMsg.ConnectType_Server {
		gate.onServerClose(msg.GetId())
	} else {
		//告诉连接对象管理服务
		id := gate.clientTransport[msg.GetId()]
		gate.SendData(id, uint16(innerMsg.InnerCmd_clientConnect), packet.GetBody())
	}
}

var errSessionKey = errors.New("Session Key Error ")

// HandlerSessionVerify 远端接入验证函数 主要处理客户端接入
func (gate *GatewayServer) HandlerSessionVerify(pack packet.IPacket) (p packet.IPacket) {
	p = packet.GenPacket()
	p.SetCmd(6001)
	msg := &protoMsg.C_GameServer{}
	e := msg.Unmarshal(pack.GetBody())
	ret := &protoMsg.S_GameServer{}
	var sid uint64
	var sessionKey string
	var sessionExpire int64
	if e != nil || msg.Uid == 0 || msg.Sessionkey == "" {
		ret.Result = 2
		goto END
	}

	// 安全验证 TODO 考虑更安全 全面的验证方式
	sessionKey, sessionExpire, sid = db.GetSessionInfo(msg.Uid)
	if sessionExpire < time.Now().Unix() || sessionKey != msg.Sessionkey {
		ret.Result = 1
		p.SetError(errSessionKey)
		goto END
	}
	gate.verifyDone(msg.Uid, sid)

	seelog.Debug("Gate Session Verify ", msg.Uid)
	p.SetTarget(msg.Uid) //发送需要
END:
	d, _ := ret.Marshal()
	p.SetBody(d)
	return
}
func (gate *GatewayServer) verifyDone(uid uint64, sid uint64) {
	if sid == 0 {
		i := discovery.GetServiceMiniLoad(config.GetServiceName(config.ServerKindGame), gate.GetFlag())
		sid = i.GetServiceNodeID()
	}
	gate.clientTransport[uid] = sid
	gate.SendToService(config.ServerKindGame, sid, uint16(innerMsg.InnerCmd_clientConnect), &innerMsg.ClientConnect{
		Id:    uid,
		State: innerMsg.ConnectState_alive,
	})
	// 重复登录
	if d := db.GetGateWaySrv(uid); d != 0 {
		if d == gate.GetSID() {
			duplicate := packet.GenPacket()
			duplicate.SetCmd(10000)
			duplicate.SetTarget(uid)
			gate.SendDirect(duplicate)
			gate.CloseSession(uid, false)
		} else {
			//关闭其他网关的session
			gate.SendToService(config.ServerKindGateway, d, uint16(innerMsg.InnerCmd_closeSession), &innerMsg.CloseSession{Id: uid})
		}
	}
	db.SetGateWaySrv(uid, gate.GetSID())
}

// HandlerBroadMsg 广播消息处理 ，单服 全服
func (gate *GatewayServer) HandlerBroadMsg(packet packet.IPacket) {
	msg := &innerMsg.PacketTransport{}
	if e := msg.Unmarshal(packet.GetBody()); e != nil {
		return
	}
	// 广播消息 使用sessionMgr  区分类型广播
}

// nodeUpdate TODO 负责函数
func (gate *GatewayServer) nodeUpdate() (output *discovery.ServiceState, status string) {
	output, status = gate.ServiceReport()
	output.Load = uint32(len(gate.clientTransport))
	return
}
