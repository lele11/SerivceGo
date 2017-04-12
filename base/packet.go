package base

import (
	"encoding/binary"
)

const (
	HEADERSIZE = 11
)

type Header struct {
	flag   uint8
	cmd    uint16
	length uint16
	refer  uint32
	kind   uint8 //接受者的类型
	sID    uint8
}

type Packet struct {
	head       Header
	body       []byte
	sess       *Session
	recvEntity IEntity
	sendEntity IEntity
}

func NewPacket() *Packet {
	return &Packet{
		head: Header{},
	}
}
func (packet *Packet) GetHead() []byte {
	h := make([]byte, HEADERSIZE)
	h[0] = packet.head.flag
	binary.LittleEndian.PutUint16(h[1:2], packet.head.cmd)
	binary.LittleEndian.PutUint16(h[3:4], packet.head.length)
	binary.LittleEndian.PutUint32(h[5:8], packet.head.refer)
	h[9] = packet.head.kind
	h[10] = packet.head.sID
	return h
}
func (packet *Packet) GetLength() uint16 {
	return packet.head.length
}
func (packet *Packet) GetRefer() uint32 {
	return packet.head.refer
}
func (packet *Packet) GetCmd() uint16 {
	return packet.head.cmd
}
func (packet *Packet) GetBody() []byte {
	return packet.body
}
func (packet *Packet) GetSession() *Session {
	return packet.sess
}

func (packet *Packet) SetSender(e IEntity) {
	packet.sendEntity = e
}
func (packet *Packet) GetSender() IEntity {
	return packet.sendEntity
}

func (packet *Packet) SetRecver(e IEntity) {
	packet.recvEntity = e
}
func (packet *Packet) GetRecver() IEntity {
	return packet.recvEntity
}
func (packet *Packet) GetSerivceInfo() uint32 {
	return uint32(packet.head.sID)
}
func (packet *Packet) SetHeader(h []byte) {
	packet.head.flag = uint8(h[0])
	packet.head.cmd = binary.LittleEndian.Uint16(h[1:3])
	packet.head.length = binary.LittleEndian.Uint16(h[3:5])
	packet.head.refer = binary.LittleEndian.Uint32(h[5:9])
	packet.head.kind = uint8(h[9])
	packet.head.sID = uint8(h[10])
}
