package discovery

import (
	"fmt"
	"game/base/util"
	"sort"
	"strings"
)

// ServiceState 服务状态
type ServiceState struct {
	Load    uint32 `json:"load"`
	Invalid bool   `json:"invalid"`
	MemSys  uint64 `json:"mem_sys"`
	GcPause uint64 `json:"gc_pause"`
}

// ServiceDesc 注册到服务发现的服务描述
type ServiceDesc struct {
	Name  string
	ID    string // 所有service中唯一的id GameName_ServiceKind_ID
	Host  string
	Port  int
	Meta  map[string]string // 细节配置
	State *ServiceState     // 状态信息
	Tag   []string

	Reporter ReportFunc
}

// GetServiceNodeID 获取服务的节点id
func (desc *ServiceDesc) GetServiceNodeID() uint64 {
	return desc.GetMetaAsUint64("id")
}

// GetLoad 获取服务的负载
func (desc *ServiceDesc) GetLoad() uint32 {
	return desc.State.Load
}

// SetMeta 设置meta数据
func (desc *ServiceDesc) SetMeta(key, value string) {
	if desc.Meta == nil {
		desc.Meta = make(map[string]string)
	}

	desc.Meta[key] = value
}

// GetMeta 获取meta
func (desc *ServiceDesc) GetMeta(name string) string {
	if desc.Meta == nil {
		return ""
	}

	return desc.Meta[name]
}

// GetMetaAsInt 获取整数
func (desc *ServiceDesc) GetMetaAsInt(name string) int {
	return util.StringToInt(desc.GetMeta(name))
}

// GetMetaAsUint64 获取uint64
func (desc *ServiceDesc) GetMetaAsUint64(name string) uint64 {
	return util.StringToUint64(desc.GetMeta(name))
}

// GetMetaAsUint32 获取uint32
func (desc *ServiceDesc) GetMetaAsUint32(name string) uint32 {
	return util.StringToUint32(desc.GetMeta(name))
}

// Address 获取服务地址
func (desc *ServiceDesc) Address() string {
	return fmt.Sprintf("%s:%d", desc.Host, desc.Port)
}

// String 格式化输出
func (desc *ServiceDesc) String() string {
	var sb strings.Builder
	if len(desc.Meta) > 0 {
		sb.WriteString("meta: [ ")
		for key, value := range desc.Meta {
			sb.WriteString(key)
			sb.WriteString("=")
			sb.WriteString(value)
			sb.WriteString(" ")
		}
		sb.WriteString("]")
	}

	return fmt.Sprintf("%s host: %s port: %d %s", desc.ID, desc.Host, desc.Port, sb.String())
}

// FormatString 标准化输出
func (desc *ServiceDesc) FormatString() string {
	var sb strings.Builder
	if len(desc.Meta) > 0 {
		type pair struct {
			key   string
			value string
		}

		var pairs []pair
		for key, value := range desc.Meta {
			pairs = append(pairs, pair{key, value})
		}

		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i].key < pairs[j].key
		})

		sb.WriteString("meta: [ ")
		for _, kv := range pairs {
			sb.WriteString(kv.key)
			sb.WriteString("=")
			sb.WriteString(kv.value)
			sb.WriteString(" ")
		}
		sb.WriteString("]")
	}

	return fmt.Sprintf("%25s host: %15s port: %5d %s", desc.ID, desc.Host, desc.Port, sb.String())
}
