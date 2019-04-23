package network

import "game/base/network/netConn"

// ConnRunner 连接服务保持对象
type ConnRunner interface {
	Start()
}

// ConnAcceptor 接入连接处理器
type ConnAcceptor interface {
	Accept(conn netconn.Conn, defaultID uint64) ConnRunner
}

// NetServer 网络服务对象
type NetServer interface {
	// 设置处理接入连接的处理对象
	SetConnAcceptor(acceptor ConnAcceptor)
	// 开始服务
	Run()
	// 关闭服务
	Close()
	//初始化服务器
	Init(addr, certKey, certFile string, tls bool)
}

// NewNetWork 创建网络服务
func NewNetWork(protocol string) NetServer {
	switch protocol {
	case "tcp":
		return &TCPServer{}
	case "ws", "wss":
		return &WssServer{}
	default:
		return nil
	}
}
