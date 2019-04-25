package sessionmgr

import (
	"game/base/network/netConn"
	"game/base/proto"
	"game/base/safelist"
	"sync"

	"game/base/packet"
	"game/base/uuid"

	"game/base/network"

	log "github.com/cihub/seelog"
)

/*
	一个独立的会话管理器
	1. 接收链接，创建会话
	2. 会话生命周期管理
	3. 会话验证
	4. 收发消息
*/

// IMsgReceiver 消息接收接口，收到的消息将由对应的结构处理
type IMsgReceiver interface {
	Receive(packet packet.IPacket)
}

// HeartBeatCmd 心跳协议号
const HeartBeatCmd = 5999

// NewSessionMgr 创建session管理器
func NewSessionMgr() *SessionMgr {
	return &SessionMgr{
		sessionByID:    &sync.Map{},
		pendingSession: &sync.Map{},
		sendBuffer:     safelist.NewSafeList(),
	}
}

// SessionMgr 管理 Sess
type SessionMgr struct {
	IMsgReceiver
	sessionByID    *sync.Map
	pendingSession *sync.Map
	verifyHandler  func(packet.IPacket) packet.IPacket
	heartBeat      bool
	sendBuffer     *safelist.SafeList
}

// SetVerifyHandler 设置session的验证函数
func (sessMgr *SessionMgr) SetVerifyHandler(f func(packet.IPacket) packet.IPacket) {
	sessMgr.verifyHandler = f
}

// SetMsgReceiver 设置消息接收结构
func (sessMgr *SessionMgr) SetMsgReceiver(r IMsgReceiver) {
	sessMgr.IMsgReceiver = r
}

// CloseSession 主动关闭一个session，normal 控制是否正常关闭，正常关闭会通知业务层
func (sessMgr *SessionMgr) CloseSession(id uint64, normal bool) {
	session := sessMgr.GetSessByID(id)
	if session != nil {
		session.Close(normal)
	}
}

// GetSessByID 获取一个session
func (sessMgr *SessionMgr) GetSessByID(id uint64) *Session {
	if i, ok := sessMgr.sessionByID.Load(id); ok {
		return i.(*Session)
	}
	return nil
}

// IsSessionExist 检查session是否存在
func (sessMgr *SessionMgr) IsSessionExist(id uint64) bool {
	_, ok := sessMgr.sessionByID.Load(id)
	return ok
}

// Run 开启执行任务，消息发送业务
func (sessMgr *SessionMgr) Run() {
	go sessMgr.doSend()
}

// Close 关闭，主要是关闭session
func (sessMgr *SessionMgr) Close() {
	sessMgr.sessionByID.Range(func(key, value interface{}) bool {
		sess := value.(*Session)
		sess.Close(true)
		return true
	})
	sessMgr.pendingSession.Range(func(key, value interface{}) bool {
		sess := value.(*Session)
		sess.Close(false)
		return true
	})
}

// getTempID session临时id 使用uuid
func (sessMgr *SessionMgr) getTempID() uint64 {
	id, _ := uuid.NextID()
	return uint64(id)
}

// verifySession 验证session
func (sessMgr *SessionMgr) verifySession(sess *Session, p packet.IPacket) bool {
	var ret bool
	if sessMgr.verifyHandler == nil || p.GetCmd() == uint16(innerMsg.InnerCmd_clientConnect) {
		//没有设置验证函数 或者 来自内部的连接，使用自定义验证消息
		m := &innerMsg.ClientConnect{}
		if m.Unmarshal(p.GetBody()) != nil {
			return false
		}
		sess.SetID(m.GetId())
		sess.kind = m.GetKind()
		// 内部连接 取消心跳
		sess.heartBeat = false
		ret = true
	} else {
		// 正常的session验证
		p = sessMgr.verifyHandler(p)
		if p.GetError() != nil {
			log.Error("Verify Session Error ", p.GetError())
			ret = false
		} else {
			sess.SetID(p.GetTarget())
			ret = true
			sess.SendPacket(p)
		}
	}
	// 验证过后，需要清理挂起
	sessMgr.pendingSession.Delete(sess.GetID())
	if ret {
		// 通过的session 进入session保存
		sessMgr.sessionByID.Store(sess.GetID(), sess)
	}
	return ret
}

// removeSession 移除session
func (sessMgr *SessionMgr) removeSession(id uint64) {
	sessMgr.sessionByID.Delete(id)
}

// Accept 实现 network 中的 connectAcceptor 接口，基于连接创建session
func (sessMgr *SessionMgr) Accept(conn netconn.Conn, id uint64) network.ConnRunner {
	sess := newSession(conn, sessMgr)
	if id == 0 { //被动会话
		sess.SetID(sessMgr.getTempID())
		sessMgr.pendingSession.Store(sess.GetID(), sess)
	} else { //主动发起
		sess.SetID(id)
		sess.heartBeat = false
		sess.isVerify = true
		sess.clearHeartBeat()
		sess.kind = innerMsg.ConnectType_Server
		sessMgr.sessionByID.Store(sess.GetID(), sess)
	}

	return sess
}

// Send 接收来自外部的消息
func (sessMgr *SessionMgr) Send(p packet.IPacket) {
	sessMgr.sendBuffer.Put(p)
}
func (sessMgr *SessionMgr) doSend() {
	for {
		<-sessMgr.sendBuffer.C
		for {
			info, err := sessMgr.sendBuffer.Pop()
			if err != nil {
				break
			}
			p := info.(packet.IPacket)
			if sess := sessMgr.GetSessByID(p.GetTarget()); sess != nil {
				sess.SendPacket(p)
			} else {
				log.Errorf("Not Found Target %d cmd %d ", p.GetTarget(), p.GetCmd())
			}
		}
	}
}
