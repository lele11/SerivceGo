package network

import (
	"crypto/tls"
	"net"
	"game/base/network/netConn"
	"sync"
	"time"

	"github.com/cihub/seelog"
)

type TCPServer struct {
	addr       string
	tls        bool //是否支持tls
	certFile   string
	certKey    string
	MaxConnNum int
	acceptor   ConnAcceptor
	ln         net.Listener
	mutexConns sync.Mutex
	wgLn       sync.WaitGroup
	cType      uint8
}

func (server *TCPServer) Init(addr, certKey, certFile string, tls bool) {
	server.addr = addr
	server.certKey = certKey
	server.certFile = certFile
	server.tls = tls
}

func (server *TCPServer) SetConnAcceptor(ca ConnAcceptor) {
	server.acceptor = ca
}
func (server *TCPServer) Run() {
	ln, err := net.Listen("tcp", server.addr)
	if err != nil {
		return
	}

	if server.tls {
		tlsConf := new(tls.Config)
		tlsConf.Certificates = make([]tls.Certificate, 1)
		tlsConf.Certificates[0], err = tls.LoadX509KeyPair(server.certFile, server.certKey)
		if err == nil {
			ln = tls.NewListener(ln, tlsConf)
		}
	}

	server.ln = ln
	seelog.Info("Start Listen Tcp ", server.addr, server.tls)
	go func() {
		for {
			var tempDelay time.Duration
			conn, err := server.ln.Accept()
			if err != nil {
				if ne, ok := err.(net.Error); ok && ne.Temporary() {
					if tempDelay == 0 {
						tempDelay = 5 * time.Millisecond
					} else {
						tempDelay *= 2
					}
					if max := 1 * time.Second; tempDelay > max {
						tempDelay = max
					}
					time.Sleep(tempDelay)
					continue
				}
				return
			}
			tempDelay = 0

			tcpConn := netConn.NewTCPConn(conn)
			runner := server.acceptor.Accept(tcpConn, 0)
			go runner.Start()
		}
	}()
	return
}

func (server *TCPServer) Close() {
	if server.ln == nil {
		return
	}
	server.ln.Close()
}
