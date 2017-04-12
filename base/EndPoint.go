package base

import (
	"errors"
	"fmt"
	"net"
)

//连接管理 消息收集 消息发送
type EndPoint struct {
	sesses   map[uint64]*Session
	recv     chan *Packet
	out      chan *Packet
	listener net.Listener
}

func newEndPoint() IEndPoint {
	return &EndPoint{
		sesses: make(map[uint64]*Session),
		recv:   make(chan *Packet, 1024),
		out:    make(chan *Packet, 1024),
	}
}

func (ep *EndPoint) Listen(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.New(err.Error())
	}
	ep.listener = l
	return nil
}

func (ep *EndPoint) Run() {
	go ep.acceptConn()
	go ep.collPacket()
	go ep.sendPacket()
}

func (ep *EndPoint) Close() {
	for _, s := range ep.sesses {
		s.close()
	}
}

func (ep *EndPoint) GetRChannel() chan IPacket {
	//	return ep.recv
	return nil
}
func (ep *EndPoint) GetWChannel() chan IPacket {
	//	return ep.recv
	return nil
}

//处理连接的有效性
func (ep *EndPoint) handleConnect(c net.Conn) {
	se := newSession(c)
	ep.sesses[se.id] = se
}
func (ep *EndPoint) removeSession(id uint64) {
	delete(ep.sesses, id)
}

func (ep *EndPoint) collPacket() {
	for {
		for _, v := range ep.sesses {
			for _, p := range v.GetPackets() {
				ep.recv <- p
			}
		}
	}
}

func (ep *EndPoint) sendPacket() {
	for {
		select {
		case p := <-ep.out:
			sess := ep.getSessionByRefer(p.GetRefer())
			if sess != nil {
				sess.send(p)
			}
		}
	}
}

func (ep *EndPoint) acceptConn() {
	if ep.listener == nil {
		fmt.Println("not listener")
		return
	}
	for {
		c, e := ep.listener.Accept()
		if e == nil {
			ep.handleConnect(c)
		}
	}
}
func (ep *EndPoint) Connect(addr string) {
	c, e := net.Dial("tcp", addr)
	if e != nil {
		fmt.Println("connect err ", e)
		return
	}
	ep.handleConnect(c)
}

func (ep *EndPoint) getSessionByRefer(refer uint32) *Session {
	//	return ep.sessByRefer[refer]
	return nil
}
