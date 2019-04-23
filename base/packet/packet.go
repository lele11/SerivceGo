package packet

import (
	"encoding/binary"
	"errors"
)

var packetError = errors.New("Packet Length Error ")

type IPacket interface {
	PackData() []byte
	GetCmd() uint16
	SetCmd(uint16)
	GetBody() []byte
	SetBody([]byte)
	GetTarget() uint64
	SetTarget(uint64)
	GetType() uint8
	SetType(uint8)
	GetError() error
	SetError(error)
	GetOrigin() uint64
	SetOrigin(uint64)
}

func GenPacket() IPacket {
	return &Packet{}
}

type Packet struct {
	cmd    uint16 // 消息号
	body   []byte //包体数据
	target uint64 //目的id
	mType  uint8  // 消息类型 0 来自客户端 1 来自服务器
	err    error
	origin uint64 //数据来源的sessionID
}

func (packet *Packet) SetOrigin(id uint64) {
	packet.origin = id
}
func (packet *Packet) GetOrigin() uint64 {
	return packet.origin
}
func (packet *Packet) GetError() error {
	return packet.err
}

func (packet *Packet) SetError(e error) {
	packet.err = e
}

func (packet *Packet) SetTarget(id uint64) {
	packet.target = id
}
func (packet *Packet) GetTarget() uint64 {
	return packet.target
}

func (packet *Packet) SetType(mType uint8) {
	packet.mType = mType
}
func (packet *Packet) GetType() uint8 {
	return packet.mType
}
func (packet *Packet) GetCmd() uint16 {
	return packet.cmd
}
func (packet *Packet) SetCmd(cmd uint16) {
	packet.cmd = cmd
}
func (packet *Packet) GetBody() []byte {
	return packet.body
}

func (packet *Packet) SetBody(body []byte) {
	packet.body = body
}

func (packet *Packet) UnpackData(data []byte) {
	if len(data) < 2 {
		packet.err = packetError
		return
	}
	packet.cmd = binary.LittleEndian.Uint16(data[0:2])
	packet.body = data[2:]
	return
}

//TODO  在发送时候直接写入 不需要中间变量
func (packet *Packet) PackData() []byte {
	buf := make([]byte, len(packet.body)+2)

	binary.LittleEndian.PutUint16(buf[0:2], packet.cmd)
	copy(buf[2:], packet.body)
	return buf
}
