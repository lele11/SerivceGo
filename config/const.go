package config

const (
	// ServerKindGateway 网关服务
	ServerKindGateway = 2
	// ServerKindGame 游戏服务
	ServerKindGame = 3
	// ServerKindLogin 登录服务
	ServerKindLogin = 4
)

// ServerName 服务类型的名称
var ServerName = map[uint32]string{
	ServerKindGateway: "Gateway",
	ServerKindGame:    "Game",
	ServerKindLogin:   "Login",
}

const (
	// GameName 游戏名称
	GameName = "dinosaur"
	// commonKey 游戏通用配置的索引
	commonKey = "/" + GameName + "/common"
)
