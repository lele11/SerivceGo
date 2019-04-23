package discovery

import (
	"fmt"
	"game/base/util"
	"sort"
	"strings"
)

type ServiceState struct {
	Load    uint32 `json:"load"`
	Invalid bool   `json:"invalid"`
	MemSys  uint64 `json:"mem_sys"`
	GcPause uint64 `json:"gc_pause"`
}

// 注册到服务发现的服务描述
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

func (desc *ServiceDesc) GetServiceNodeID() uint64 {
	return desc.GetMetaAsUint64("id")
}
func (desc *ServiceDesc) GetLoad() uint32 {
	return desc.State.Load
}
func (desc *ServiceDesc) SetMeta(key, value string) {
	if desc.Meta == nil {
		desc.Meta = make(map[string]string)
	}

	desc.Meta[key] = value
}

func (desc *ServiceDesc) GetMeta(name string) string {
	if desc.Meta == nil {
		return ""
	}

	return desc.Meta[name]
}

func (desc *ServiceDesc) GetMetaAsInt(name string) int {
	return util.StringToInt(desc.GetMeta(name))
}
func (desc *ServiceDesc) GetMetaAsUint64(name string) uint64 {
	return util.StringToUint64(desc.GetMeta(name))
}
func (desc *ServiceDesc) GetMetaAsUint32(name string) uint32 {
	return util.StringToUint32(desc.GetMeta(name))
}
func (desc *ServiceDesc) Address() string {
	return fmt.Sprintf("%s:%d", desc.Host, desc.Port)
}

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
