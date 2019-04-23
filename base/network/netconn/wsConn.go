package netconn

import (
	"bytes"
	"io"
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WSConn websocket 对象
type WSConn struct {
	io.Reader //Read(p []byte) (n int, err error)
	io.Writer //Write(p []byte) (n int, err error)
	sync.Mutex
	bufLocks  chan error //当有写入一次数据设置一次
	buffer    bytes.Buffer
	conn      *websocket.Conn
	readFirst bool
	closeFlag bool
	cType     uint8
	realIP    string
}

// NewWSConn 创建ws链接结构
func NewWSConn(conn *websocket.Conn) *WSConn {
	wsConn := new(WSConn)
	wsConn.conn = conn
	wsConn.bufLocks = make(chan error)
	wsConn.readFirst = false

	return wsConn
}

// Close 关闭
func (wsConn *WSConn) Close() error {
	if wsConn.closeFlag {
		return nil
	}
	return wsConn.conn.Close()
}

// SetType 设置链接类型
func (wsConn *WSConn) SetType(t uint8) {
	wsConn.cType = t
}

// GetType 获取链接类型
func (wsConn *WSConn) GetType() uint8 {
	return wsConn.cType
}

// Destroy 销毁函数
func (wsConn *WSConn) Destroy() {
	wsConn.Lock()
	defer wsConn.Unlock()
	if wsConn.closeFlag {
		return
	}
	wsConn.closeFlag = true
	wsConn.conn.Close()
}

// Write  写入
func (wsConn *WSConn) Write(p []byte) (int, error) {
	err := wsConn.conn.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

// ReadMessage 读取消息
func (wsConn *WSConn) ReadMessage() (d []byte, err error) {
	_, d, err = wsConn.conn.ReadMessage()
	if err != nil {
		wsConn.Destroy()
		return
	}
	return

}

// Read 读二进制
func (wsConn *WSConn) Read(p []byte) (n int, err error) {
	err = <-wsConn.bufLocks //等待写入数据
	if err != nil {
		//读取数据出现异常了
		return
	}
	if wsConn.buffer.Len() == 0 {
		//再等一次
		err = <-wsConn.bufLocks //等待写入数据
		if err != nil {
			//读取数据出现异常了
			return
		}
	}
	return wsConn.buffer.Read(p)
}

// LocalAddr 本地地址
func (wsConn *WSConn) LocalAddr() net.Addr {
	return wsConn.conn.LocalAddr()
}

// RemoteAddr 远端地址
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

// RealIP 真实ip 在经过网络转发后，远端地址不是客户端的真实IP
func (wsConn *WSConn) RealIP() string {
	return wsConn.realIP
}

// SetRealIP 设置真实IP 通过代理参数设置
func (wsConn *WSConn) SetRealIP(ip string) {
	wsConn.realIP = ip
}
