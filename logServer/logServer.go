package main

import (
	"encoding/json"
	"game/base/logger"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/cihub/seelog"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	e := logger.NewRemoteServer("")
	if e != nil {
		seelog.Error(e)
	}
	logger.RemoteServerWorkerNum(10)
	logger.RemoteServerRun()
	s := &LogServer{}
	s.startHTTP("127.0.0.1:1000")
}

type LogServer struct {
}

// close 停止服务器
func (l *LogServer) close() {
	logger.RemoteServerClose()
}

func (l *LogServer) startHTTP(addr string) {
	server := &http.Server{
		Addr:         addr,
		Handler:      l,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  1 * time.Minute,
	}

	listener, e := net.Listen("tcp", addr)
	if e != nil {
		seelog.Error("Start Listen Error ", e)
		return
	}
	go server.Serve(listener)
}

func (l *LogServer) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	content, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		seelog.Error("Log Request Error ", err)
		return
	}
	if len(content) == 0 {
		seelog.Error("Log Request Error ", err)
		return
	}
	param := &logger.LogData{}
	err = json.Unmarshal(content, param)
	if err != nil {
		seelog.Error("Log Request Error ", err)
		return
	}
	logger.RemoteLogAddLog(param)
}
