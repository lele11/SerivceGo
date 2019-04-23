package netConn

import (
	"bytes"
	"io"
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WSConn struct {
	io.Reader //Read(p []byte) (n int, err error)
	io.Writer //Write(p []byte) (n int, err error)
	sync.Mutex
	buf_lock  chan error //当有写入一次数据设置一次
	buffer    bytes.Buffer
	conn      *websocket.Conn
	readfirst bool
	closeFlag bool
	cType     uint8
	realIp    string
}

func NewWSConn(conn *websocket.Conn) *WSConn {
	wsConn := new(WSConn)
	wsConn.conn = conn
	wsConn.buf_lock = make(chan error)
	wsConn.readfirst = false

	return wsConn
}
func (wsConn *WSConn) Close() error {
	if wsConn.closeFlag {
		return nil
	}
	return wsConn.conn.Close()
}

func (wsConn *WSConn) SetType(t uint8) {
	wsConn.cType = t
}
func (wsConn *WSConn) GetType() uint8 {
	return wsConn.cType
}

func (wsConn *WSConn) Destroy() {
	wsConn.Lock()
	defer wsConn.Unlock()
	if wsConn.closeFlag {
		return
	}
	wsConn.closeFlag = true
	wsConn.conn.Close()
}

func (wsConn *WSConn) Write(p []byte) (int, error) {
	err := wsConn.conn.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (wsConn *WSConn) ReadMessage() (d []byte, err error) {
	_, d, err = wsConn.conn.ReadMessage()
	if err != nil {
		wsConn.Destroy()
		return
	}
	return

}

// goroutine not safe
func (wsConn *WSConn) Read(p []byte) (n int, err error) {
	err = <-wsConn.buf_lock //等待写入数据
	if err != nil {
		//读取数据出现异常了
		return
	}
	if wsConn.buffer.Len() == 0 {
		//再等一次
		err = <-wsConn.buf_lock //等待写入数据
		if err != nil {
			//读取数据出现异常了
			return
		}
	}
	return wsConn.buffer.Read(p)
}

func (wsConn *WSConn) LocalAddr() net.Addr {
	return wsConn.conn.LocalAddr()
}

func (wsConn *WSConn) RemoteAddr() net.Addr {
	return wsConn.conn.RemoteAddr()
}

// A zero value for t means I/O operations will not time out.
func (wsConn *WSConn) SetDeadline(t time.Time) error {
	err := wsConn.conn.SetWriteDeadline(t)
	if err != nil {
		return err
	}
	return wsConn.conn.SetWriteDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls.
// A zero value for t means Read will not time out.
func (wsConn *WSConn) SetReadDeadline(t time.Time) error {
	return wsConn.conn.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (wsConn *WSConn) SetWriteDeadline(t time.Time) error {
	return wsConn.conn.SetWriteDeadline(t)
}

func (wsConn *WSConn) RealIP() string {
	return wsConn.realIp
}
func (wsConn *WSConn) SetRealIP(ip string) {
	wsConn.realIp = ip
}
