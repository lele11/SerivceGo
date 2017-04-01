package nServer

import (
	"encoding/binary"
)

const (
	HEADERSIZE = 9
)

type Header struct {
	flag   uint8
	cmd    uint16
	length uint16
	refer  uint32
}

type Packet struct {
	head Header
	body []byte
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

func (packet *Packet) SetHeader(h []byte) {
	packet.head.flag = uint8(h[0])
	packet.head.cmd = binary.LittleEndian.Uint16(h[1:3])
	packet.head.length = binary.LittleEndian.Uint16(h[3:5])
	packet.head.refer = binary.LittleEndian.Uint32(h[5:9])
}
