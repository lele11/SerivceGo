package discovery

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/cihub/seelog"
)

// ReportFunc consul健康监测上报方法，上报中携带必要的负载信息
type ReportFunc func() (state *ServiceState, status string)

// localServices 本地服务管理
var localServices sync.Map

// LocalService 本地上报服务
type LocalService struct {
	name string
	ttl  time.Duration
	f    ReportFunc
}

// update 服务上报进程
func (ls *LocalService) update() {
	if ls.f == nil {
		return
	}
	for {
		output, status := ls.f()
		o, _ := json.Marshal(output)
		e := Default.UpdateService(ls.name, status, string(o))
		if e != nil {
			seelog.Error(e)
		}
		time.Sleep(ls.ttl)
	}
}
