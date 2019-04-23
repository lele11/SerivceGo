package network

import (
	"game/base/network/netConn"
	"net"
	"net/http"
	"time"

	"github.com/cihub/seelog"
	"github.com/gorilla/websocket"
)

type WssServer struct {
	addr      string
	listener  net.Listener
	acceptor  ConnAcceptor
	maxConns  int
	maxMsgLen uint32
	upgrader  websocket.Upgrader
	certFile  string
	certKey   string
	cType     uint8 //
	tls       bool
}

// Init 初始化
func (wssServer *WssServer) Init(addr, certkey, certfile string, tls bool) {
	wssServer.addr = addr
	wssServer.certKey = certkey
	wssServer.certFile = certfile
	wssServer.tls = tls
}

// SetConnAcceptor 设置接受
func (wssServer *WssServer) SetConnAcceptor(ca ConnAcceptor) {
	wssServer.acceptor = ca
}

// ServeHTTP
func (wssServer *WssServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	conn, err := wssServer.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	//conn.SetReadLimit(int64(wssServer.maxMsgLen))
	c := netconn.NewWSConn(conn)
	//i := r.Header.Get("X-Forwarded-For") 多层代理 记录所有的代理ip
	i := r.Header.Get("X-real-ip")
	c.SetRealIP(i)
	runner := wssServer.acceptor.Accept(c, 0)
	runner.Start()
}

// Run 运行
func (wssServer *WssServer) Run() {
	var err error
	wssServer.listener, err = net.Listen("tcp", wssServer.addr)
	if err != nil {
		seelog.Error("Listen error ", wssServer.addr, err)
		return
	}
	wssServer.upgrader = websocket.Upgrader{
		HandshakeTimeout: 10 * time.Second,
		CheckOrigin:      func(_ *http.Request) bool { return true },
	}
	httpServer := &http.Server{
		Addr:           wssServer.addr,
		Handler:        wssServer,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1024,
	}
	if wssServer.tls {
		go httpServer.ServeTLS(wssServer.listener, wssServer.certFile, wssServer.certKey)
	} else {
		go httpServer.Serve(wssServer.listener)
	}
	seelog.Info("Start WssServer Ok ,", wssServer.addr)
	return
}

// Close 关闭
func (wssServer *WssServer) Close() {
	if wssServer.listener == nil {
		return
	}
	wssServer.listener.Close()
}
