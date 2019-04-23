package main

import (
	"flag"
	"game/base/discovery"
	"game/base/discovery/consul"
	"game/base/logger"
	"game/base/redis"
	"game/base/serverMgr"
	"game/base/util"
	"game/config"
	_ "game/gateway"
	_ "game/login"
	"net/http"
	_ "net/http/pprof"

	"github.com/cihub/seelog"
)

func main() {
	uu := ""
	for k, v := range config.ServerName {
		uu += "		" + util.Uint32ToString(k) + "		" + v + "\n\r"
	}

	id := flag.String("ID", "", "启动的服务id ")
	kind := flag.Uint("kind", 0, "启动的服务类型：\n\r "+uu)
	consulAddr := flag.String("consul", "", "Consul的地址")
	flag.Parse()
	go func() {
		addr := "localhost:" + *id
		e := http.ListenAndServe(addr, nil)
		if e != nil {
			seelog.Error("Start PProf Error ", e)
		}
	}()
	//*kind = 3
	//*id = "3001"
	//*consulAddr = "127.0.0.1:8500"
	discovery.Default = consulsd.NewDiscovery(*consulAddr, "/"+config.GameName+"/")
	discovery.WatchServiceAll()
	//载入通用配置
	common, e := config.LoadCommonConfig()
	if e != nil {
		panic("Load Common Config Error , Not Found Info In Consul K/V")
		return
	}
	//日志初始化
	dd := logger.Init(common.LogPath+config.ServerName[uint32(*kind)], common.LogLevel)
	if dd != nil {
		defer dd.Flush()
	}
	//redis初始化
	redis.SetConfig(common.Redis.Host, common.Redis.Password)
	srv := serverMgr.GetService(uint32(*kind))
	if srv == nil {
		seelog.Errorf("Not Found Service Kind %d, you Need Use %s  ", *kind, uu)
		return
	}

	var cfg *config.ServerConfig
	info := discovery.GetServiceById(config.GetServiceName(uint32(*kind)), config.GetServiceID(uint32(*kind), util.StringToUint64(*id)))
	if info == nil {
		seelog.Error("LoadServer Config error ", config.GetServiceID(uint32(*kind), util.StringToUint64(*id)))
		return
	} else {
		seelog.Info(info.String())
		cfg = &config.ServerConfig{
			ID:        info.GetMetaAsUint64("id"),
			Kind:      uint32(*kind),
			Host:      info.Host,
			Port:      info.Port,
			Protocol:  "ws",
			Meta:      info.Meta,
			Tag:       info.Tag,
			Name:      info.GetMeta("name"),
			OpenStamp: info.GetMeta("open"),
			State:     info.GetMetaAsUint32("state"),
			Domain:    info.GetMeta("domain"),
		}
		if info.Port == 0 {
			cfg.Port, _ = util.GetFreePort()
		}
	}
	srv.Init(cfg)
	srv.Run()
}
