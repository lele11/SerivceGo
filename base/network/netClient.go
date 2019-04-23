package network

import (
	"game/base/network/netConn"
	"net"

	"github.com/gorilla/websocket"
)

// NetClient 客户端接口
type NetClient interface {
	// 设置处理接入连接的处理对象
	SetConnAcceptor(acceptor ConnAcceptor)
	Dial(protocol string, addr string, id uint64) error
}

// NetNetClient 创建
func NetNetClient() NetClient {
	return &NClient{}
}

// NClient 客户端
type NClient struct {
	acceptor ConnAcceptor
}

// SetConnAcceptor 设置链接接受对象
func (netClient *NClient) SetConnAcceptor(ac ConnAcceptor) {
	netClient.acceptor = ac
}

// Dial 发起链接
func (netClient *NClient) Dial(protocol, addr string, id uint64) error {
	var runner ConnRunner
	switch protocol {
	case "ws", "wss":
		c, _, err := websocket.DefaultDialer.Dial(protocol+"://"+addr, nil)
		if err != nil {
			return err
		}
		runner = netClient.acceptor.Accept(netconn.NewWSConn(c), id)
	case "tcp":
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return err
		}
		runner = netClient.acceptor.Accept(netconn.NewTCPConn(conn), id)
	}
	runner.Start()

	return nil
}
