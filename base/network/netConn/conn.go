package netConn

import (
	"net"
)

//TODO 为了兼容wss的不定长包 应该使用Read方法 按照tcpConn的读取方式 需要定义包结构
//TODO 目前在tcpConn中定义 包结构读取 ，以后项目应该通用包结构

// Conn 链接结构
type Conn interface {
	//连接实例
	net.Conn
	//读取消息
	ReadMessage() (d []byte, err error)
	//关闭连接
	Destroy()
	//获取连接类型
	GetType() uint8
	//设置连接类型
	SetType(uint8)
	// 获取真实ip
	RealIP() string
	// 设置真实 ip
	SetRealIP(string)
}
