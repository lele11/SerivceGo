package base

import (
	"errors"
	"game/base/discovery"
	"game/base/logger"
	"game/base/msgHandler"
	"game/base/network"
	"game/base/packet"
	"game/base/proto"
	"game/base/sessionMgr"
	"game/base/util"
	"game/config"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cihub/seelog"
)

//可以被打包的消息结构 为了兼容proto
type IMsg interface {
	Unmarshal(data []byte) error
	Marshal() (dAtA []byte, err error)
	Size() (n int)
}

type IService interface {
	Run()
	Init(serverConfig *config.ServerConfig)
	GetSID() uint64
	GetKind() uint32
}

type Service struct {
	cfg                    *config.ServerConfig
	id                     uint64
	addr                   string
	protocol               string
	port                   int
	kind                   uint32
	*msgHandler.MsgHandler //消息处理器
	network.NetServer      //网络服务器Server
	network.NetClient
	*logger.RemoteLog                        //日志记录
	sessionMgr        *sessionMgr.SessionMgr //连接管理器
	loopFunc          func()                 //上层业务的Loop方法
	closeFunc         func()                 //上层业务的关闭方法
	reportFunc        discovery.ReportFunc
}

// NewService 创建一个新的服务器
func NewService(id uint64, addr string, port int, kind uint32, protocol string, cfg *config.ServerConfig) *Service {
	srv := &Service{
		id:       id,
		addr:     addr,
		protocol: protocol,
		port:     port,
		kind:     kind,
		cfg:      cfg,
	}
	srv.NetServer = network.NewNetWork(srv.protocol)
	if srv.NetServer == nil {
		return nil
	}
	srv.NetClient = network.NetNetClient()
	srv.sessionMgr = sessionMgr.NewSessionMgr()
	srv.MsgHandler = msgHandler.NewMsgHandler()
	srv.RemoteLog = logger.NewLogger()
	srv.Init()
	return srv
}

func (service *Service) Init() {
	service.SetReporter(service.ServiceReport)
	service.NetServer.Init(service.addr+":"+util.IntToString(service.port), config.CertKey, config.CertFile, false)
	service.NetServer.SetConnAcceptor(service.sessionMgr)
	service.NetClient.SetConnAcceptor(service.sessionMgr)
	service.sessionMgr.SetMsgReceiver(service.MsgHandler)
}
func (service *Service) GetKind() uint32 {
	return service.kind
}
func (service *Service) GetSID() uint64 {
	return service.id
}
func (service *Service) GetFlag() string {
	return service.cfg.Meta["flag"]
}
func (service *Service) Close() {
	if service.closeFunc != nil {
		service.closeFunc()
	}
	service.NetServer.Close()
	service.sessionMgr.Close()
	seelog.Infof("Service %d Close ", service.id)
	seelog.Flush()
}

// Run 逻辑入口
func (service *Service) Run() {
	service.NetServer.Run()
	service.sessionMgr.Run()
	service.RemoteLog.Run()
	service.register()
	seelog.Infof("Service Start Node %d Addr %s port %d", service.id, service.addr, service.port)
	service.doLoop()
}

func (service *Service) doLoop() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, os.Interrupt)
	ticker := time.NewTicker(time.Millisecond * time.Duration(1000/30))
	for {
		select {
		case sig := <-c:
			seelog.Info(sig.String())
			service.Close()
			return
		case <-ticker.C:
			service.DoConsumeMsg()
			if service.loopFunc != nil {
				service.loopFunc()
			}
		}
	}
}

func (service *Service) SetUpdate(f func()) {
	service.loopFunc = f
}

func (service *Service) SetClose(f func()) {
	service.closeFunc = f
}

func (service *Service) SetSessionVerifyHandler(f func(packet packet.IPacket) packet.IPacket) {
	service.sessionMgr.SetVerifyHandler(f)
}

func (service *Service) SendData(target uint64, cmd uint16, data []byte) {
	p := &packet.Packet{}
	p.SetCmd(cmd)
	p.SetTarget(target)
	if data != nil {
		p.SetBody(data)
	}
	if target == service.GetSID() {
		service.Receive(p)
	} else {
		service.sessionMgr.Send(p)
	}
}

//SendMsg 发送消息
func (service *Service) SendMsg(id uint64, cmd uint16, msg IMsg) {
	var data []byte
	var e error
	if msg != nil {
		data, e = msg.Marshal()
		if e != nil {
			seelog.Error("Marshal Msg Error ", e)
			return
		}
	}

	service.SendData(id, cmd, data)
}
func (service *Service) SendDirect(p packet.IPacket) {
	s := service.sessionMgr.GetSessByID(p.GetTarget())
	if s != nil {
		s.SendPacket(p.PackData())
	}
}

//SendToService 发送到某个服务器 ，与sendMsg的区别
func (service *Service) SendToService(sType uint32, id uint64, cmd uint16, msg IMsg) {
	if id == 0 {
		return
	}
	if service.ConnectToServer(sType, id) != nil {
		return
	}
	var data []byte
	var e error
	if msg != nil {
		data, e = msg.Marshal()
		if e != nil {
			seelog.Error("Marshal Msg Error ", e)
			return
		}
	}
	service.SendData(id, cmd, data)
}
func (service *Service) ConnectToServer(sType uint32, serverId uint64) error {
	if service.GetSID() == serverId {
		return nil
	}
	if service.sessionMgr.GetSessByID(serverId) != nil {
		return nil
	}
	info := discovery.GetServiceById(config.GetServiceName(sType), config.GetServiceID(sType, serverId))
	if info == nil {
		seelog.Error("Connect id Not Found ", serverId)
		return errors.New("error")
	}
	if e := service.Dial("ws", info.Host, info.Port, serverId); e != nil {
		seelog.Errorf("Connect id %d Info %v Error %s", serverId, info, e)
		return e
	}
	m := &innerMsg.ClientConnect{}
	m.Id = service.id
	m.Kind = innerMsg.ConnectType_Server
	service.SendMsg(serverId, uint16(innerMsg.InnerCmd_clientConnect), m)
	return nil
}

func (service *Service) CloseSession(id uint64, normal bool) {
	service.sessionMgr.CloseSession(id, normal)
}

func (service *Service) Dial(protocol, addr string, port int, id uint64) error {
	addr = addr + ":" + util.IntToString(port)
	err := service.NetClient.Dial(protocol, addr, id)
	if err != nil {
		return err
	}
	return nil
}

func (service *Service) NewNetWork(protocol, addr, certKey, certFile string, tls bool) network.NetServer {
	n := network.NewNetWork(protocol)
	n.Init(addr, certKey, certFile, tls)
	n.SetConnAcceptor(service.sessionMgr)
	return n
}

func (service *Service) SetReporter(f discovery.ReportFunc) {
	service.reportFunc = f
}
func (service *Service) register() {
	desc := &discovery.ServiceDesc{
		ID:       config.GetServiceID(service.GetKind(), service.GetSID()),
		Name:     config.GetServiceName(service.GetKind()),
		Host:     service.addr,
		Port:     service.port,
		Meta:     service.cfg.Meta,
		Tag:      service.cfg.Tag,
		Reporter: service.reportFunc,
	}

	e := discovery.RegisterService(desc)
	if e != nil {
		seelog.Error("regester Service Error ", e)
	}
}

func (service *Service) ServiceReport() (output *discovery.ServiceState, status string) {
	output = &discovery.ServiceState{}
	output.MemSys, output.GcPause = util.GetMemState()
	status = "passing"
	return
}
