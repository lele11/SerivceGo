package login

import (
	"encoding/json"
	"game/base/discovery"
	"game/base/network"
	"game/base/network/netConn"
	"game/base/packet"
	"game/base/serverMgr"
	"game/base/util"
	"game/config"
	"game/db"
	"game/protoMsg"
	"math/rand"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/cihub/seelog"
)

func init() {
	serverMgr.Register(&LoginServer{})
}

// LoginServer 登录服务
type LoginServer struct {
	out        network.NetServer
	randSource *rand.Rand
	cfg        *config.ServerConfig
}

// Init 初始化
func (login *LoginServer) Init(cfg *config.ServerConfig) {
	login.out = network.NewNetWork(cfg.Protocol)
	login.out.SetConnAcceptor(login)
	login.randSource = rand.New(rand.NewSource(time.Now().UnixNano()))
	addr := cfg.Host + ":" + util.IntToString(cfg.Port)
	login.out.Init(addr, "", "", false)
	login.cfg = cfg
}

// GetSID 获取服务id
func (login *LoginServer) GetSID() uint64 {
	return login.cfg.ID
}

// GetKind 获取服务类型
func (login *LoginServer) GetKind() uint32 {
	return config.ServerKindLogin
}

// Run 主逻辑
func (login *LoginServer) Run() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	e := discovery.RegisterService(&discovery.ServiceDesc{
		Name:     config.GetServiceName(login.GetKind()),
		ID:       config.GetServiceID(login.GetKind(), login.GetSID()),
		Host:     login.cfg.Host,
		Port:     login.cfg.Port,
		Meta:     login.cfg.Meta,
		Tag:      login.cfg.Tag,
		Reporter: login.nodeUpdate,
	})
	if e != nil {
		seelog.Error("regester Service Error ", e)
	}
	login.out.Run()
	<-c
}

// Accept 链接接收函数
func (login *LoginServer) Accept(conn netConn.Conn, defaultId uint64) network.ConnRunner {
	for {
		data, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			return login
		}

		p := &packet.Packet{}
		p.UnpackData(data)
		if p.GetError() != nil {
			seelog.Error("Recv Cmd ", p.GetError())
			return login
		}
		switch p.GetCmd() {
		case 6000:
			ret := login.login(p.GetBody())
			if conn != nil {
				conn.Write(ret.PackData())
				conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			}
		case 5100:
			ret := login.GetServerList(p.GetBody())
			if conn != nil {
				conn.Write(ret.PackData())
				conn.SetReadDeadline(time.Now().Add(2 * time.Second))
			}
		default:
			seelog.Error("unknown cmd ", p.GetCmd())
		}
	}
}

// Start 接口要求
func (login *LoginServer) Start() {

}

// login 登录函数
func (login *LoginServer) login(param []byte) (r packet.IPacket) {
	r = &packet.Packet{}
	r.SetCmd(6000)
	var e error
	var newPeople bool
	msg := &protoMsg.C_LoginServer{}
	e = msg.Unmarshal(param)
	if e != nil {
		seelog.Errorf("UnMarshal Msg Error ", e)
		return
	}

	var account string

	retMsg := &protoMsg.S_LoginServer{}
	// 验证区服是否可以提供服务
	if msg.ServerID != 0 {
		info := discovery.GetServiceById(config.ServerName[config.ServerKindGame], config.GetServiceID(config.ServerKindGame, msg.ServerID))
		if info == nil {
			retMsg.Result = 4
			goto END
		}
	}

	//TODO 验证渠道 非法渠道需要处理
	if msg.LoginChannel == "" {
		msg.LoginChannel = "1"
		//retMsg.Result = 5
		//goto END
	}
	if msg.GetAccount() == "" {
		retMsg.Result = 1
		goto END
	}

	account = msg.GetAccount()
	retMsg.Uid, newPeople = db.GetAccountUID(msg.LoginChannel, account, msg.ServerID)
	retMsg.Sessionkey = util.RandomString(login.randSource)
	seelog.Infof("Login success name:%s,openId:%s,channel:%s,uid:%d,sessKey[%s]", msg.UserInfo.NickName, account, "", retMsg.Uid, retMsg.Sessionkey)
	retMsg.IsNew = 0
	if newPeople {
		retMsg.IsNew = 1
	}
	db.StoreSessionInfo(retMsg.Uid, retMsg.Sessionkey, msg.ServerID)
	db.StorePlayerInfo(retMsg.Uid, msg.GetUserInfo())
	if msg.ServerID > 0 {
		db.AddAccountRecent(msg.LoginChannel, msg.Account, msg.ServerID)
	}
	if login.cfg.Domain == "" {
		retMsg.Result = 3
		goto END
	}

	retMsg.Address = login.cfg.Domain
	if !strings.HasPrefix(retMsg.Address, "ws") {
		retMsg.Address = "wss://" + retMsg.Address
	}
	retMsg.Account = ""
	retMsg.Result = 0
END:
	b, e := retMsg.Marshal()
	if e != nil {
		seelog.Errorf("Login Error ", e)
	}
	r.SetBody(b)
	return
}

// nodeUpdate TODO 负责函数
func (login *LoginServer) nodeUpdate() (output *discovery.ServiceState, status string) {
	output = &discovery.ServiceState{}
	output.MemSys, output.GcPause = util.GetMemState()
	status = "passing"
	return
}

// ServerInfo 区服信息
type ServerInfo struct {
	Id      uint64 `json:"id"`
	Default int    `json:"default"`
	Group   string `json:"group"`
	Name    string `json:"name"`
	IsNew   int    `json:"is_new"`
	Login   string `json:"login"`
}

// ServerListRet 区服列表
type ServerListRet struct {
	ServerList []*ServerInfo `json:"server_list"`
	Recent     []uint64      `json:"recent"`
}

// ServerListReq 区服列表请求
type ServerListReq struct {
	Account string `json:"account"`
	Channel string `json:"channel"`
}

// GetServerList 获取区服列表
func (login *LoginServer) GetServerList(param []byte) (r *packet.Packet) {
	r = &packet.Packet{}
	r.SetCmd(5100)
	req := &ServerListReq{}
	if e := json.Unmarshal(param, req); e != nil {
		seelog.Error(e)
		return
	}

	list := discovery.GetServiceList(config.ServerName[config.ServerKindGame])
	loginInfo := discovery.GetServiceOne(config.ServerName[config.ServerKindLogin], login.cfg.Meta["flag"])
	ret := &ServerListRet{}
	for _, l := range list {
		info := &ServerInfo{}
		info.Id = l.GetMetaAsUint64("id")
		info.Name = l.GetMeta("name")
		info.Default = 0
		info.Group = l.GetMeta("group")
		info.Login = l.GetMeta("login")
		if info.Login == "" {
			info.Login = loginInfo.GetMeta("domain")
		}
		ret.ServerList = append(ret.ServerList, info)
	}
	ret.Recent = db.GetAccountRecent(req.Channel, req.Account)
	d, e := json.Marshal(ret)
	if e != nil {
		seelog.Errorf("GetServerList ", e)
	}
	r.SetBody(d)
	return
}
