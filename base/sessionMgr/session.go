package sessionMgr

import (
	"game/base/network/netConn"
	"game/base/proto"
	"game/base/safelist"
	"time"

	"game/base/packet"

	"github.com/cihub/seelog"
)

func newSession(conn netConn.Conn, mgr *SessionMgr) *Session {
	s := &Session{
		conn:       conn,
		sendBuffer: safelist.NewSafeList(),
		mgr:        mgr,
		heartBeat:  true,
	}
	s.closeChan = make(chan bool, 1)
	s.conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	return s
}

// Session 代表一个网络连接
type Session struct {
	id         uint64
	conn       netConn.Conn
	closeChan  chan bool
	isVerify   bool
	isClosed   bool
	sendBuffer *safelist.SafeList
	mgr        *SessionMgr
	heartBeat  bool
	kind       innerMsg.ConnectType
}

// Start
func (session *Session) Start() {
	go session.sendLoop()
	go session.recvLoop()
}
func (session *Session) SendPacket(data []byte) {
	if session.isClosed {
		return
	}
	session.sendBuffer.Put(data)
}
func (session *Session) clearHeartBeat() {
	session.conn.SetReadDeadline(time.Time{})
}

// Touch 记录心跳状态
func (session *Session) Touch() {
	if session.heartBeat {
		session.conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	}
}

func (session *Session) recvLoop() {
	for {
		if session.isClosed {
			return
		}
		data, err := session.conn.ReadMessage()
		if err != nil {
			session.Close(true)
			return
		}

		pack := &packet.Packet{}
		pack.UnpackData(data)
		if e := pack.GetError(); e != nil {
			session.Close(true)
			return
		}
		pack.SetOrigin(session.GetID())
		if session.isVerify {
			if pack.GetCmd() == HeartBeatCmd {
				session.sendBuffer.Put(pack.PackData())
			} else {
				session.mgr.Receive(pack)
			}
		} else {
			if session.mgr.verifySession(session, pack) {
				session.clearHeartBeat()
				session.isVerify = true
				seelog.Info("Session Verify result ", session.GetID(), session.isVerify, session.RemoteAddr())
			} else {
				session.Close(false)
				return
			}
		}
		session.Touch()
	}
}

func (session *Session) sendLoop() {
	for {
		select {
		case <-session.closeChan:
			return
		case <-session.sendBuffer.C:
			for {
				data, err := session.sendBuffer.Pop()
				if err != nil {
					break
				}
				_, e := session.conn.Write(data.([]byte))
				if e != nil {
					seelog.Error("Receive Error ", e)
					session.Close(true)
					return
				}
			}
		}
	}
}
func (session *Session) onClose() {
	p := &packet.Packet{}
	p.SetOrigin(session.GetID())
	p.SetCmd(uint16(innerMsg.InnerCmd_clientConnect))
	m := &innerMsg.ClientConnect{
		Id:    session.GetID(),
		State: innerMsg.ConnectState_dead,
		Kind:  session.kind,
	}
	d, _ := m.Marshal()
	p.SetBody(d)
	session.mgr.Receive(p)
}

// Close 关闭
func (session *Session) Close(normal bool) {
	if session.isClosed {
		return
	}
	session.isClosed = true
	if normal {
		session.onClose()
	}
	session.mgr.removeSession(session.GetID())
	go func() {
		//等待消息发送结束 关闭
		times := 0
		for {
			times++
			if session.sendBuffer.IsEmpty() && times > 100 {
				session.closeChan <- true
				close(session.closeChan)
				session.conn.Destroy()
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		seelog.Infof("Session Close %d %s %t", session.GetID(), session.RemoteAddr(), normal)
	}()
}

// SetID 设置ID
func (session *Session) SetID(id uint64) {
	session.id = id
}

// GetID 获取ID
func (session *Session) GetID() uint64 {
	return session.id
}

func (session *Session) RemoteAddr() string {
	return session.conn.RemoteAddr().String()
}
func (session *Session) SetVerify(v bool) {
	session.isVerify = v
}