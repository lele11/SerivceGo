package logger

import (
	"bytes"
	"encoding/json"
	"net/http"
	"game/base/safelist"

	"crypto/tls"
	"net"

	"github.com/cihub/seelog"
	"golang.org/x/net/http2"
)

// TODO 是否可以做成单例 不需要业务层创建
// TODO 是否兼容文本日志

//NewLogger 创建logger对象
func NewLogger() *RemoteLog {
	return &RemoteLog{
		list: safelist.NewSafeList(),
	}
}

//Logger logger结构，http发送日志数据到日志服务
type RemoteLog struct {
	remoteURL string             //日志服务地址
	list      *safelist.SafeList //日志数据缓存队列 TODO 性能需要检验
	r         bool               //是否启动
	hc        *http.Client       //缓存http2 client
}

//Push 对外接口，发送日志
func (l *RemoteLog) PushRemoteLog(data interface{}) {
	if !l.r {
		return
	}
	s, e := json.Marshal(data)
	if e != nil {
		seelog.Error("Send Remote Logger Error ", e, data)
		return
	}
	l.list.Put(s)
}

//Init 初始化日志对象
func (l *RemoteLog) SetRemoteLogUrl(url string) *RemoteLog {
	if url == "" {
		return l
	}
	//TODO 检测地址是否有效
	l.remoteURL = url
	tr := &http2.Transport{
		DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) {
			return net.Dial(netw, addr)
		},
	}
	l.hc = &http.Client{Transport: tr}
	l.r = true
	return l
}

//Run 启动
func (l *RemoteLog) Run() {
	if !l.r {
		return
	}
	go func() {
		for {
			<-l.list.C
			for {
				d, e := l.list.Pop()
				if e != nil {
					break
				}
				rr := bytes.NewReader(d.([]byte))
				_, e = l.hc.Post(l.remoteURL, "application/json", rr)
				if e != nil {
					seelog.Error("send log Error ", e)
				}
			}
		}
	}()
}
