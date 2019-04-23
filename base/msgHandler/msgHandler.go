package msgHandler

import (
	"game/base/packet"
	"game/base/safelist"

	"github.com/cihub/seelog"
)

/*
	管理消息处理函数映射 cmd => func
*/
type handlerFunc func(packet packet.IPacket)

// MsgHandler 消息处理
type MsgHandler struct {
	handlers       map[uint16][]handlerFunc
	recvBuffer     *safelist.SafeList
	defaultHandler func(packet packet.IPacket)
}

// NewMsgHandler 创建
func NewMsgHandler() *MsgHandler {
	m := &MsgHandler{
		recvBuffer: safelist.NewSafeList(),
		handlers:   make(map[uint16][]handlerFunc),
	}
	return m
}

// SetDefaultHandler 设置默认消息处理函数
func (msgHandler *MsgHandler) SetDefaultHandler(f handlerFunc) {
	msgHandler.defaultHandler = f
}

// AttachHandler 添加消息映射函数
func (msgHandler *MsgHandler) AttachHandler(cmd uint16, f handlerFunc) {
	msgHandler.handlers[cmd] = append(msgHandler.handlers[cmd], f)
}

// DoConsumeMsg 处理消息逻辑
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

// Receive 接受消息
func (msgHandler *MsgHandler) Receive(packet packet.IPacket) {
	msgHandler.recvBuffer.Put(packet)
}
