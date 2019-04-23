package consulsd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/consul/api"
)

// check包括：node的check和service check
func isServiceHealth(entry *api.ServiceEntry) bool {
	for _, check := range entry.Checks {
		if check.Status != api.HealthPassing {
			return false
		}
	}

	return true
}

// BytesToAny 二进制流转换函数
func BytesToAny(data []byte, dataPtr interface{}) error {
	switch ret := dataPtr.(type) {
	case *int:
		v, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return err
		}
		*ret = int(v)
		return nil
	case *float32:
		v, err := strconv.ParseFloat(string(data), 32)
		if err != nil {
			return err
		}
		*ret = float32(v)
		return nil
	case *float64:
		v, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			return err
		}
		*ret = float64(v)
		return nil
	case *bool:
		v, err := strconv.ParseBool(string(data))
		if err != nil {
			return err
		}
		*ret = v
		return nil
	case *string:
		*ret = string(data)
		return nil
	default:
		return json.Unmarshal(data, dataPtr)
	}
}

// AnyToBytes 二进制流转换函数
func AnyToBytes(data interface{}) ([]byte, error) {
	switch v := data.(type) {
	case int, int32, int64, uint32, uint64, float32, float64, bool:
		return []byte(fmt.Sprint(data)), nil
	case string:
		return []byte(v), nil
	default:
		raw, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		return raw, nil
	}
}
