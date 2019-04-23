package netconn

import (
	"encoding/binary"
	"net"
	"sync"
	"time"
)

// TCPConn tcp 连接
type TCPConn struct {
	//io.Reader //Read(p []byte) (n int, err error)
	//io.Writer //Write(p []byte) (n int, err error)
	sync.Mutex
	//bufLocks  chan bool //当有写入一次数据设置一次
	//buffer    bytes.Buffer
	conn      net.Conn
	closeFlag bool
	cType     uint8 //连接的类型
}

// NewTCPConn 创建tcp链接结构
func NewTCPConn(conn net.Conn) *TCPConn {
	tcpConn := new(TCPConn)
	tcpConn.conn = conn

	return tcpConn
}

// doDestroy 销毁
func (tcpConn *TCPConn) doDestroy() {
	tcpConn.conn.(*net.TCPConn).SetLinger(0)
	tcpConn.conn.Close()

	if !tcpConn.closeFlag {
		tcpConn.closeFlag = true
	}
}

// Destroy 销毁函数
func (tcpConn *TCPConn) Destroy() {
	tcpConn.Lock()
	defer tcpConn.Unlock()

	tcpConn.doDestroy()
}

// SetType 设置链接类型
func (tcpConn *TCPConn) SetType(t uint8) {
	tcpConn.cType = t
}

// GetType 获取链接类型
func (tcpConn *TCPConn) GetType() uint8 {
	return tcpConn.cType
}

// Close 关闭
func (tcpConn *TCPConn) Close() error {
	tcpConn.Lock()
	defer tcpConn.Unlock()
	if tcpConn.closeFlag {
		return nil
	}

	tcpConn.closeFlag = true
	return tcpConn.conn.Close()
}

// Write  写入
func (tcpConn *TCPConn) Write(b []byte) (n int, err error) {
	tcpConn.Lock()
	defer tcpConn.Unlock()
	if tcpConn.closeFlag || b == nil {
		return
	}
	l := uint32(len(b))
	d := make([]byte, l+4)
	binary.LittleEndian.PutUint32(d[:4], l)
	copy(d[4:], b)
	return tcpConn.conn.Write(d)
}

// ReadMessage 读取消息
func (tcpConn *TCPConn) ReadMessage() (d []byte, err error) {
	//TODO 读取数据包 为了兼容wss  效率 消息号格式
	l := make([]byte, 4)
	tcpConn.conn.Read(l)
	d = make([]byte, binary.LittleEndian.Uint32(l))
	_, err = tcpConn.conn.Read(d)
	if err != nil {
		return
	}
	return
}

// Read 读二进制
func (tcpConn *TCPConn) Read(b []byte) (int, error) {
	return tcpConn.conn.Read(b)
}

// LocalAddr 本地地址
func (tcpConn *TCPConn) LocalAddr() net.Addr {
	return tcpConn.conn.LocalAddr()
}

// RemoteAddr 远端地址
func (tcpConn *TCPConn) RemoteAddr() net.Addr {
	return tcpConn.conn.RemoteAddr()
}

// SetDeadline A zero value for t means I/O operations will not time out.
func (tcpConn *TCPConn) SetDeadline(t time.Time) error {
	return tcpConn.conn.SetDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls.
// A zero value for t means Read will not time out.
func (tcpConn *TCPConn) SetReadDeadline(t time.Time) error {
	return tcpConn.conn.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (tcpConn *TCPConn) SetWriteDeadline(t time.Time) error {
	return tcpConn.conn.SetWriteDeadline(t)
}

// RealIP 真实ip 在经过网络转发后，远端地址不是客户端的真实IP
func (tcpConn *TCPConn) RealIP() string {
	return ""
}

// SetRealIP 设置真实IP 通过代理参数设置
func (tcpConn *TCPConn) SetRealIP(ip string) {

}
