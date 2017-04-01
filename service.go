package nServer

import (
	"fmt"
	"time"
)

func NewService() *Service {
	s := &Service{
		EndPoint: NewEndPoint(),
		handlers: make(map[uint16]func([]byte)),
	}
	s.EndPoint.Run()
	return s
}

type Service struct {
	*EndPoint
	handlers map[uint16]func([]byte)
}

func (service *Service) Run() {
	ticker := time.NewTicker(10 * time.Millisecond)
	for {
		select {
		case packet := <-service.recv:
			service.Handle(packet)
		case <-ticker.C:
			continue
		}
	}
}

func (service *Service) Attach(cmd uint16, func_ func([]byte)) {
	service.handlers[cmd] = func_
}

func (service *Service) Handle(packet *Packet) {
	fmt.Println("hanle packet cmd", packet.GetCmd())
	h := service.handlers[packet.GetCmd()]
	if h != nil {
		h(packet.GetBody())
	}
}

func (service *Service) Send(refer uint32, cmd uint16, data []byte) {
	packet := &Packet{}
	packet.body = data
	packet.head = Header{
		cmd:   cmd,
		refer: refer,
	}
	service.out <- packet
}
