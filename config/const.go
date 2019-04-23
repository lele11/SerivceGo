package config

const (
	ServerKindGateway = 2
	ServerKindGame    = 3
	ServerKindLogin   = 4
)

var ServerName = map[uint32]string{
	ServerKindGateway: "Gateway",
	ServerKindGame:    "Game",
	ServerKindLogin:   "Login",
}

const (
	CertKey  = ""
	CertFile = ""
)
const (
	GameName  = "dinosaur"
	commonKey = "/" + GameName + "/common"
)
