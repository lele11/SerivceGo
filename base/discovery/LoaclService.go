package discovery

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/cihub/seelog"
)

type ReportFunc func() (state *ServiceState, status string)

var localServices sync.Map

type LocalService struct {
	name string
	ttl  time.Duration
	f    ReportFunc
}

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
