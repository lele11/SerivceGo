package base

import (
	"io"
	"net"
	"service/utils"
	"time"
)

const (
	SessionClose = iota
	SessionRecon
	SessionWorking
)

func newSession(c net.Conn) *Session {
	s := &Session{
		id:     utils.GenUUID(),
		conn:   c,
		signal: make(chan int, 1),
	}
	go s.run()
	return s
}

type Session struct {
	Entity
	id         uint64
	conn       net.Conn
	signal     chan int
	packetPool []*Packet
}

func (sess *Session) run() {
	ticket := time.NewTicker(time.Millisecond * 100)
	for {
		select {
		case s := <-sess.signal:
			sess.doSignal(s)
			return
		case <-ticket.C:
			sess.read()
		}
	}
}

func (sess *Session) doSignal(i int) {
	if i == SessionClose {
		sess.conn.Close()
	}
}

func (sess *Session) send(p *Packet) {
	buf := make([]byte, HEADERSIZE+p.GetLength())
	copy(buf[:HEADERSIZE], p.GetHead())
	copy(buf[HEADERSIZE:], p.GetBody())
	sess.conn.Write(buf)
}

func (sess *Session) read() {
	h := make([]byte, HEADERSIZE)
	if _, err := io.ReadFull(sess.conn, h); err != nil {
		return
	}
	packet := NewPacket()
	packet.SetHeader(h)
	if packet.head.length > 0 {
		packet.body = make([]byte, packet.head.length)
		if _, err := io.ReadFull(sess.conn, packet.body); err != nil {
			return
		}
	}
	//session 里填充session的信息
	packet.sess = sess
	sess.packetPool = append(sess.packetPool, packet)
}

func (sess *Session) close() {
	sess.signal <- SessionClose
}

func (sess *Session) GetPackets() (ss []*Packet) {

	l := len(sess.packetPool)
	ss = sess.packetPool[:l]
	sess.packetPool = sess.packetPool[l:]
	return
}
