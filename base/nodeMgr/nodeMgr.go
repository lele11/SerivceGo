package nodeMgr

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"game/base/proto"
	"game/config"
	"sync"
	"time"

	"github.com/cihub/seelog"
)

/*
 管理系统中所有的节点信息
*/
func NewNodeMgr() *NodeMgr {
	return &NodeMgr{
		nodeInfo: map[uint64]*innerMsg.ServerInfoBroad_ServerInfo{},
		closed:   make(chan bool, 1),
		upTimer:  2,
	}
}

type NodeMgr struct {
	lock     sync.RWMutex
	nodeInfo map[uint64]*innerMsg.ServerInfoBroad_ServerInfo
	ticker   *time.Ticker
	closed   chan bool
	callback func() *innerMsg.ServerUpdate
	upTimer  uint32
	c        bool
}

func (nm *NodeMgr) Init() {
	nm.getAllNode()
}

func (nm *NodeMgr) GetMinLoadNodeId(kind uint32) uint64 {
	var min uint32
	var info *innerMsg.ServerInfoBroad_ServerInfo
	for _, node := range nm.nodeInfo {
		if node.Id / 1000 != uint64(kind) {
			continue
		}
		if min == 0 || min > node.GetLoad() {
			min = node.GetLoad()
			info = node
		}
	}
	if info == nil {
		return 0
	}
	return info.GetId()
}

func (nm *NodeMgr) getAllNode() {
	v := url.Values{}
	r, e := http.PostForm("http://"+config.CenterHttp+"/getAllNode", v)
	if e != nil {
		seelog.Error("getAllNode Error ", e)
		return
	}

	data, e := ioutil.ReadAll(r.Body)
	if e != nil {
		seelog.Error("getAllNode Error ", e)
		return
	}
	r.Body.Close()
	if len(data) == 0 {
		return
	}

	msg := &innerMsg.ServerInfoBroad{}
	if e := msg.Unmarshal(data); e != nil {
		seelog.Error("HandlerServerInfoSync packet error ", e)
		return
	}
	nm.lock.Lock()
	for _, node := range msg.List {
		nm.nodeInfo[node.GetId()] = node
	}
	nm.lock.Unlock()
}
func (nm *NodeMgr) SetUpdateTime(t uint32) {
	nm.upTimer = t
}
func (nm *NodeMgr) NodeUpdateFunc(f func() *innerMsg.ServerUpdate) {
	nm.callback = f
}
func (nm *NodeMgr) Run() {
	if nm.c {
		return
	}
	nm.ticker = time.NewTicker(time.Duration(nm.upTimer) * time.Second)
	go func() {
		for {
			nm.Update()
			select {
			case <-nm.ticker.C:
				nm.Update()
				nm.getAllNode()
			case <-nm.closed:
				return
			}
		}
	}()
}
func (nm *NodeMgr) close() {
	if nm.c {
		return
	}
	if nm.ticker != nil {
		nm.ticker.Stop()
	}
	nm.closed <- true
	close(nm.closed)
	nm.c = true
}

func (nm *NodeMgr) Close() {
	nm.close()
}
func (nm *NodeMgr) Update() {
	if nm.callback == nil {
		return
	}
	info := nm.callback()
	s, _ := info.Marshal()
	v := url.Values{}
	v.Add("info", string(s))
	_, e := http.PostForm("http://"+config.CenterHttp+"/upNodeState", v)
	if e != nil {
		seelog.Error("GetNodeConfig Error ", e)
		return
	}
}

func (nm *NodeMgr) GetNodeById(id uint64) uint64 {
	nm.lock.RLock()
	i := nm.nodeInfo[id]
	nm.lock.RUnlock()
	if i == nil {
		return 0
	}
	return i.Node
}
func (nm *NodeMgr) GetNodeInfoById(id uint64) *innerMsg.ServerInfoBroad_ServerInfo {
	nm.lock.RLock()
	i := nm.nodeInfo[id]
	nm.lock.RUnlock()
	if i == nil {
		return nil
	}
	return i
}
