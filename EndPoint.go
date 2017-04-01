package nServer

import (
	"bytes"
	"common/config"
	"fmt"
	"io"
	"net"
)

type EndPoint struct {
	sesses     map[uint32]*Session
	refer      uint32
	acceptChan chan net.Conn
	recv       chan *Packet
	out        chan *Packet
	listener   net.Listener
	buffer     *bytes.Buffer
}

func NewEndPoint() *EndPoint {
	return &EndPoint{
		sesses:     make(map[uint32]*Session),
		recv:       make(chan *Packet, 1024),
		out:        make(chan *Packet, 1024),
		acceptChan: make(chan net.Conn, 1024),
		buffer:     new(bytes.Buffer),
	}
}
func (ep *EndPoint) addSession(c net.Conn) {
	ep.refer++
	s := Session{
		conn:  c,
		refer: ep.refer,
	}
	ep.sesses[s.refer] = &s
	fmt.Println("set new session")
}
func (ep *EndPoint) removeSession(id uint32) {
	delete(ep.sesses, id)
}

func (ep *EndPoint) RecvPacket() {
	defer fmt.Println("recv loop return")
	for {
		for _, v := range ep.sesses {
			if v.close {
				continue
			}
			p := v.read()
			if p != nil {
				ep.recv <- p
				fmt.Println("recv packet")
			}
		}
	}
}

func (ep *EndPoint) SendPacket() {
	defer fmt.Println("send loop return")
	for {
		select {
		case p := <-ep.out:
			sess := ep.sesses[p.GetRefer()]
			if sess == nil {
				return
			}
			if sess.close {
				continue
			}
			sess.send(p)
		}
	}
}

func (ep *EndPoint) AcceptConn() {
	defer fmt.Println("accept loop return")
	if ep.listener == nil {
		fmt.Println("not listener")
		return
	}
	fmt.Println("begin wait conn")
	for {
		c, e := ep.listener.Accept()
		if e == nil {
			ep.addSession(c)
		}
	}
}
func (ep *EndPoint) Connect(addr string) {
	c, e := net.Dial("tcp", addr)
	if e != nil {
		fmt.Println("connect err ", e)
		return
	}
	ep.addSession(c)
}
func (ep *EndPoint) Init() {
	var err error
	ep.listener, err = net.Listen("tcp", config.GetMustString("addr"))
	if err != nil {
		fmt.Println("listen err ", err)
		return
	}
	fmt.Println("listen ok", ep.listener.Addr().String())
}

func (ep *EndPoint) Run() {
	ep.Init()
	go ep.AcceptConn()
	go ep.RecvPacket()
	go ep.SendPacket()
}

type Session struct {
	refer uint32
	conn  net.Conn
	close bool
}

func (sess *Session) send(p *Packet) {
	buf := make([]byte, HEADERSIZE+p.GetLength())
	copy(buf[:HEADERSIZE], p.GetHead())
	copy(buf[HEADERSIZE:], p.GetBody())
	sess.conn.Write(buf)
}

func (sess *Session) read() *Packet {
	h := make([]byte, HEADERSIZE)
	if _, err := io.ReadFull(sess.conn, h); err != nil {
		return nil
	}
	packet := NewPacket()
	packet.SetHeader(h)
	fmt.Println("read packet head", packet.head)
	if packet.head.length > 0 {
		packet.body = make([]byte, packet.head.length)
		if _, err := io.ReadFull(sess.conn, packet.body); err != nil {
			fmt.Println("read body err", err)
			return nil
		}
	}
	fmt.Println("read packet ok")
	return packet
}
