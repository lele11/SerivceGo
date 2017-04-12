package base

import (
	"os"
	"os/signal"
	"service/config"
	"service/utils"
	"time"
)

const (
	CentralType = 1
)

func CreateNode() *Node {
	n := &Node{
		IEndPoint:   newEndPoint(),
		services:    make(map[uint64]IService),
		closeSignal: make(chan os.Signal),
		entities:    make(map[uint32]IEntity),
	}
	n.log = utils.NewLogger("node_"+utils.ChangeUint32ToString(n.id), config.GetMustString("logPath"), 0)
	n.timer_1 = utils.NewTimer(1)
	return n
}

type Node struct {
	IEndPoint
	id          uint32
	services    map[uint64]IService
	timer_1     *utils.Timer
	log         *utils.Logger
	closeSignal chan os.Signal
	entities    map[uint32]IEntity //服务实例
	tmpId       uint32
}

func (n *Node) GetServiceByID(id uint64) IService {
	s := n.services[id]
	if s != nil {
		return s
	}
	return nil
}

func (n *Node) handle(p IPacket) {
	srv := n.GetServiceByID(p.GetRecver().GetID())
	if srv != nil {
		if p.GetCmd() == CMD_REGISTER {
			n.AddEntity(srv)
		}
		srv.Recv(p)
	} else {
		n.log.Info("Not Found service %d  Hanlde  %d from %d", p.GetRecver().GetID(), p.GetCmd(), p.GetSender().GetID())
	}
}

func (n *Node) getCentralEntity() IEntity {
	for _, e := range n.entities {
		if e.GetType() == CentralType {
			return e
		}
	}
	return nil
}

func (n *Node) Start() {
	//监听
	listenIp := config.GetMustString("addr")
	if listenIp == "" {
		n.log.Error("Not found listen ip in %s ", config.GetConfigFile())
		return
	}
	e := n.Listen(listenIp)
	if e != nil {
		n.log.Error("Listen Error %s", e.Error())
		return
	}
	n.log.Info("Liste Success %s", listenIp)
	//初始化服务
	n.tmpId++
	services := config.GetMustArray("service")
	for _, kind := range services {
		n.tmpId++
		s := newService(kind, n)
		go s.Run()
		n.services[s.GetID()] = s
	}
	n.doRun()
}

func (n *Node) doRun() {
	//启动服务
	signal.Notify(n.closeSignal)
	ticker := time.NewTicker(10 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			n.doFrame()
		case p := <-n.GetRChannel():
			n.handle(p)
		case <-n.closeSignal:
			n.shutDown()
			return
		}
	}
}

func (n *Node) shutDown() {
	n.log.Info("node closing...")
	//处理session
	n.Close()
	//处理service
	for _, s := range n.services {
		s.Close()
	}
	time.Sleep(1)
}

func (n *Node) doFrame() {
	now := uint64(time.Now().Unix())
	if n.timer_1.Update(now) {
		//TODO 底层的状态数据收集

		n.log.Debug("service %v", n.services)

	}
}

func (n *Node) AddEntity(s IService) {
	e := s.GenEntity()
	n.entities[uint32(e.GetID())] = e
}

func (n *Node) GetEntityByID(id uint32) IEntity {
	return n.entities[id]
}

func (n *Node) GetEntityByType(kind uint32) []IEntity {
	var d []IEntity
	for _, e := range n.entities {
		if e.GetType() == kind {
			d = append(d, e)
		}
	}
	return d
}

func (n *Node) SendMsg(sender IEntity, recver uint32, p IPacket) {
	if sender.GetID() == uint64(recver) {
		p.SetRecver(sender)
		n.handle(p)
		return
	}
	e := n.GetEntityByID(recver)
	if e == nil {
		n.log.Error("Not Found Recver %d", recver)
		return
	}
	p.SetSender(sender)
	p.SetRecver(e)
	n.GetWChannel() <- p

}
