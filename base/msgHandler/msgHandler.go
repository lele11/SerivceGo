package msgHandler

import (
	"game/base/packet"
	"game/base/safelist"

	"github.com/cihub/seelog"
)

/*
	管理消息处理函数映射
*/
type MsgHandler struct {
	handlers       map[uint16][]func(packet packet.IPacket)
	recvBuffer     *safelist.SafeList
	defaultHandler func(packet packet.IPacket)
}

func NewMsgHandler() *MsgHandler {
	m := &MsgHandler{
		recvBuffer: safelist.NewSafeList(),
		handlers:   make(map[uint16][]func(packet packet.IPacket)),
	}
	return m
}
func (msgHandler *MsgHandler) SetDefaultHandler(f func(packet packet.IPacket)) {
	msgHandler.defaultHandler = f
}
func (msgHandler *MsgHandler) AttachHandler(cmd uint16, f func(packet packet.IPacket)) {
	msgHandler.handlers[cmd] = append(msgHandler.handlers[cmd], f)
}

func (msgHandler *MsgHandler) DoConsumeMsg() {
	for {
		info, err := msgHandler.recvBuffer.Pop()
		if err != nil {
			break
		}
		p := info.(packet.IPacket)
		if f, ok := msgHandler.handlers[p.GetCmd()]; ok {
			for _, d := range f {
				d(p)
			}
		} else {
			if msgHandler.defaultHandler == nil {
				seelog.Errorf("Not Found handler for cmd %d", p.GetCmd())
			} else {
				msgHandler.defaultHandler(p)
			}
		}
	}
}

func (msgHandler *MsgHandler) Receive(packet packet.IPacket) {
	msgHandler.recvBuffer.Put(packet)
}
