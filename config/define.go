package config

import (
	"game/base/discovery"
	"game/base/util"

	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	ID        uint64 `yaml:"id"`
	Kind      uint32 `yaml:"kind"`
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Protocol  string `yaml:"protocol"`
	Name      string `yaml:"name"`
	OpenStamp string `yaml:"openStamp"`
	State     uint32 `yaml:"state"`
	Domain    string `yaml:"domain"`
	Meta      map[string]string
	Tag       []string
}

type CommonConfig struct {
	Redis     *RedisConfig `yaml:"redis"`
	LogLevel  string       `yaml:"logLevel"`
	LogPath   string       `yaml:"logPath"`
	RemoteLog string       `yaml:"remoteLog"`
}
type RedisConfig struct {
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
}

func GetServiceID(kind uint32, sid uint64) string {
	return GameName + "_" + ServerName[kind] + "_" + util.Uint64ToString(sid)
}
func GetServiceName(kind uint32) string {
	return GameName + "_" + ServerName[kind]
}
func LoadCommonConfig() (*CommonConfig, error) {
	common := ""
	if e := discovery.Default.GetValue(commonKey, &common); e != nil {
		return nil, e
	}
	c := &CommonConfig{}
	if e := yaml.Unmarshal([]byte(common), &c); e != nil {
		return nil, e
	}
	return c, nil
}
