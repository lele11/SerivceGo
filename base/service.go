package base

import (
	"errors"
	"game/base/discovery"
	"game/base/logger"
	"game/base/msgHandler"
	"game/base/network"
	"game/base/packet"
	"game/base/proto"
	"game/base/sessionmgr"
	"game/base/util"
	"game/config"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cihub/seelog"
)

// IMsg 可以被打包的消息结构 为了兼容proto
type IMsg interface {
	Unmarshal(data []byte) error
	Marshal() (dAtA []byte, err error)
	Size() (n int)
}

// IService 基础服务接口
type IService interface {
	Run()
	Init(serverConfig *config.ServerConfig)
	GetSID() uint64
	GetKind() uint32
}

// Service 基础服务结构
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
	sessionMgr        *sessionmgr.SessionMgr //连接管理器
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
	srv.sessionMgr = sessionmgr.NewSessionMgr()
	srv.MsgHandler = msgHandler.NewMsgHandler()
	srv.RemoteLog = logger.NewLogger()
	srv.Init()
	return srv
}

// Init 初始化函数
func (service *Service) Init() {
	service.SetReporter(service.ServiceReport)
	service.NetServer.Init(service.addr+":"+util.IntToString(service.port), "", "", false)
	service.NetServer.SetConnAcceptor(service.sessionMgr)
	service.NetClient.SetConnAcceptor(service.sessionMgr)
	service.sessionMgr.SetMsgReceiver(service.MsgHandler)
}

// GetKind 获取上层业务类型
func (service *Service) GetKind() uint32 {
	return service.kind
}

// GetSID 获取上层业务id
func (service *Service) GetSID() uint64 {
	return service.id
}

// GetFlag 获取上层业务标记
func (service *Service) GetFlag() string {
	return service.cfg.Meta["flag"]
}

// Close 关闭
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
	service.RemoteLog.Run()
	service.register()
	seelog.Infof("Service Start Node %d Addr %s port %d", service.id, service.addr, service.port)
	service.doLoop()
}

// doLoop 业务逻辑处理
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

// SetUpdate 设置上层业务内部逻辑处理函数
func (service *Service) SetUpdate(f func()) {
	service.loopFunc = f
}

// SetClose 设置上层业务关闭回调函数
func (service *Service) SetClose(f func()) {
	service.closeFunc = f
}

// SetSessionVerifyHandler 设置session自定义验证函数
func (service *Service) SetSessionVerifyHandler(f func(packet packet.IPacket) packet.IPacket) {
	service.sessionMgr.SetVerifyHandler(f)
}

// SendData 发送二进制数据
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

// SendDirect 直接发送至session
func (service *Service) SendDirect(p packet.IPacket) {
	s := service.sessionMgr.GetSessByID(p.GetTarget())
	if s != nil {
		s.SendPacket(p)
	}
}

// SendToService 发送到某个服务器 ，与sendMsg的区别
func (service *Service) SendToService(sType uint32, id uint64, cmd uint16, msg IMsg) {
	if id == 0 {
		return
	}
	if service.connectToServer(sType, id) != nil {
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

// connectToServer 连接其他节点
func (service *Service) connectToServer(sType uint32, serverID uint64) error {
	if service.GetSID() == serverID {
		return nil
	}
	if service.sessionMgr.GetSessByID(serverID) != nil {
		return nil
	}
	info := discovery.GetServiceByID(config.GetServiceName(sType), config.GetServiceID(sType, serverID))
	if info == nil {
		seelog.Error("Connect id Not Found ", serverID)
		return errors.New("error")
	}
	if e := service.dial("ws", info.Host, info.Port, serverID); e != nil {
		seelog.Errorf("Connect id %d Info %v Error %s", serverID, info, e)
		return e
	}
	m := &innerMsg.ClientConnect{}
	m.Id = service.id
	m.Kind = innerMsg.ConnectType_Server
	service.SendMsg(serverID, uint16(innerMsg.InnerCmd_clientConnect), m)
	return nil
}

// CloseSession 关闭session
func (service *Service) CloseSession(id uint64, normal bool) {
	service.sessionMgr.CloseSession(id, normal)
}

// dial 主动建立连接
func (service *Service) dial(protocol, addr string, port int, id uint64) error {
	addr = addr + ":" + util.IntToString(port)
	err := service.NetClient.Dial(protocol, addr, id)
	if err != nil {
		return err
	}
	return nil
}

// NewNetWork 创建网络服务
func (service *Service) NewNetWork(protocol, addr, certKey, certFile string, tls bool) network.NetServer {
	n := network.NewNetWork(protocol)
	n.Init(addr, certKey, certFile, tls)
	n.SetConnAcceptor(service.sessionMgr)
	return n
}

// SetReporter 设置状态上报函数
func (service *Service) SetReporter(f discovery.ReportFunc) {
	service.reportFunc = f
}

// register 注册服务
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

// ServiceReport 默认服务上报函数
func (service *Service) ServiceReport() (output *discovery.ServiceState, status string) {
	output = &discovery.ServiceState{}
	output.MemSys, output.GcPause = util.GetMemState()
	status = "passing"
	return
}
