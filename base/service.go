package base

import (
	"fmt"
	"service/config"
	"service/utils"
	"time"
)

func newService(kind uint32, n *Node) *Service {
	s := &Service{
		sType:    kind,
		node:     n,
		state:    ServiceStateInit,
		handlers: make(map[uint16]func(*Packet)),
		timer_1:  utils.NewTimer(1),
	}
	s.sId = uint64(n.tmpId)
	if kind == 1 {
		s.sId = 1
	}
	return s
}

type Service struct {
	IEntity
	node     *Node
	log      *utils.Logger
	timer_1  *utils.Timer
	shutdown chan int
	handlers map[uint16]func(*Packet)
	sType    uint32
	sId      uint64
	recv     chan *Packet
	state    int
}

func (service *Service) GetID() uint64 {
	return service.sId
}
func (service *Service) Start() {
	service.Attach(CMD_REGISTER, service.init)
}

func (service *Service) GenEntity() IEntity {
	e := &Entity{
		id:   service.sId,
		kind: service.sType,
	}
	return e

}

func (service *Service) Close() {

}

func (service *Service) Recv(p IPacket) {
	service.recv <- p.(*Packet)
}

func (service *Service) IsWorking() bool {
	return false
}

func (service *Service) Run() {
	service.Start()
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			service.do()
		case p := <-service.recv:
			service.handler(p)
		}
	}
}

func (service *Service) init(p *Packet) {
	service.log = utils.NewLogger("service", config.GetMustString("logPath"), 0)
}

func (service *Service) register() {
	service.Send(CentralID, CMD_REGISTER, []byte{})
}

func (service *Service) do() {
	now := uint64(time.Now().Unix())
	if service.timer_1.Update(now) {
		if service.state == ServiceStateInit {
			service.register()
		}
	}
}

func (service *Service) handler(packet IPacket) {
	fmt.Println("hanle packet cmd", packet.GetCmd())

	h := service.handlers[packet.GetCmd()]
	if h != nil {
		h(packet.(*Packet))
	}
}

func (service *Service) Attach(cmd uint16, func_ func(*Packet)) {
	service.handlers[cmd] = func_
}

//send 有多种  发给某一个端 发给某一类端 发给一类的某个端
func (service *Service) Send(id uint32, cmd uint16, data []byte) {
	packet := &Packet{
		head: Header{
			cmd: cmd,
		},
		body: data,
	}
	service.node.SendMsg(service, id, packet)
}
