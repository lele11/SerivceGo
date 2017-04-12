package base

// 服务 也作为一个实例
type IService interface {
	Start()
	Close()
	Send(uint32, uint16, []byte)
	Attach(uint16, func(*Packet))
	handler(IPacket)
	IsWorking() bool
	GetID() uint64
	Recv(IPacket)
	GenEntity() IEntity
}

// 全部实例 的基础
type IEntity interface {
	GetID() uint64
	GetType() uint32
	GetNode() uint32
}

type IEndPoint interface {
	Listen(string) error
	Run()
	Connect(string)
	Close()
	GetRChannel() chan IPacket
	GetWChannel() chan IPacket
}
type IPacket interface {
	GetCmd() uint16
	GetLength() uint16
	SetSender(IEntity)
	GetSender() IEntity
	SetRecver(IEntity)
	GetRecver() IEntity
}
