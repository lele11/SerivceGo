package network

import (
	"net"
	"game/base/network/netConn"

	"github.com/gorilla/websocket"
)

type NetClient interface {
	// 设置处理接入连接的处理对象
	SetConnAcceptor(acceptor ConnAcceptor)
	Dial(protocol string, addr string, id uint64) error
}

func NetNetClient() NetClient {
	return &NClient{}
}

type NClient struct {
	acceptor ConnAcceptor
}

func (netClient *NClient) SetConnAcceptor(ac ConnAcceptor) {
	netClient.acceptor = ac
}
func (netClient *NClient) Dial(protocol, addr string, id uint64) error {
	var runner ConnRunner
	switch protocol {
	case "ws", "wss":
		c, _, err := websocket.DefaultDialer.Dial(protocol+"://"+addr, nil)
		if err != nil {
			return err
		}
		runner = netClient.acceptor.Accept(netConn.NewWSConn(c), id)
	case "tcp":
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return err
		}
		runner = netClient.acceptor.Accept(netConn.NewTCPConn(conn), id)
	}
	runner.Start()

	return nil
}
