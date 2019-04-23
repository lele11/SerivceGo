package config

import (
	"game/base/discovery"
	"game/base/util"

	"gopkg.in/yaml.v2"
)

// ServerConfig 服务配置
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

// CommonConfig 通用配置
type CommonConfig struct {
	Redis     *RedisConfig `yaml:"redis"`
	LogLevel  string       `yaml:"logLevel"`
	LogPath   string       `yaml:"logPath"`
	RemoteLog string       `yaml:"remoteLog"`
}

// RedisConfig redis配置
type RedisConfig struct {
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
}

// GetServiceID 获取服务的ID
func GetServiceID(kind uint32, sid uint64) string {
	return GameName + "_" + ServerName[kind] + "_" + util.Uint64ToString(sid)
}

// GetServiceName 获取服务的名称
func GetServiceName(kind uint32) string {
	return GameName + "_" + ServerName[kind]
}

// LoadCommonConfig 加载服务的通用配置
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
