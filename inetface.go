package nServer

type IService interface {
	Run()
	Stop()
	Send()
	Connect()
	Listen()
	Accept()
	Attach()
	Handle()
}
