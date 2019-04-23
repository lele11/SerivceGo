package consulsd

import (
	"errors"
	"sync"

	"github.com/cihub/seelog"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/watch"
)

// SetValue 设置键值数据
func (cd *consulDiscovery) SetValue(key string, dataPtr interface{}) error {
	raw, err := AnyToBytes(dataPtr)
	if err != nil {
		return err
	}
	_, err = cd.client.KV().Put(&api.KVPair{
		Key:   key,
		Value: raw,
	}, nil)

	return err
}

// GetValue 获取KV数据，根据参数做转换
func (cd *consulDiscovery) GetValue(key string, valuePtr interface{}) error {
	data, err := cd.GetRawValue(key)
	if err != nil {
		return err
	}

	return BytesToAny(data, valuePtr)
}

// GetRawValue 获取二进制数据
func (cd *consulDiscovery) GetRawValue(key string) ([]byte, error) {
	if raw, ok := cd.kvCache.Load(key); ok {
		meta := raw.(*KVMeta)
		return meta.Value(), nil
	}
	// cache中没找到直接获取
	kvpair, _, err := cd.client.KV().Get(key, nil)
	if err != nil {
		return nil, err
	}
	if kvpair == nil {
		return nil, ErrValueNotExists
	}
	return kvpair.Value, nil
}

// DeleteValue 删除数据
func (cd *consulDiscovery) DeleteValue(key string) error {
	_, err := cd.client.KV().Delete(key, nil)
	return err
}

var (
	// ErrValueNotExists 不存在
	ErrValueNotExists = errors.New("value not exists")
)

// startWatchKV 开始监控
func (cd *consulDiscovery) startWatchKV(prefix string) {
	if prefix == "" {
		prefix = "/"
	}
	plan, err := watch.Parse(map[string]interface{}{
		"type":   "keyprefix",
		"prefix": prefix,
	})
	if err != nil {
		seelog.Error("startWatchKV:", err)
		return
	}

	plan.Handler = cd.onKVListChanged
	go plan.Run(cd.consulAddr)
}

// KVMeta KV数据结构
type KVMeta struct {
	value      []byte
	valueGuard sync.RWMutex
	Plan       *watch.Plan
}

// SetValue 设置值
func (km *KVMeta) SetValue(v []byte) {
	km.valueGuard.Lock()
	km.value = v
	km.valueGuard.Unlock()
}

// Value 获取值
func (km *KVMeta) Value() []byte {
	km.valueGuard.RLock()
	defer km.valueGuard.RUnlock()
	return km.value
}

// onKVListChanged 变化回调
func (cd *consulDiscovery) onKVListChanged(u uint64, data interface{}) {
	kvNames, ok := data.(api.KVPairs)
	if !ok {
		return
	}

	for _, kv := range kvNames {
		// 已经在cache里的,肯定添加过watch了
		if _, ok := cd.kvCache.Load(kv.Key); ok {
			continue
		}
		plan, err := watch.Parse(map[string]interface{}{
			"type": "key",
			"key":  kv.Key,
		})

		if err == nil {
			plan.Handler = cd.onKVChanged
			go plan.Run(cd.consulAddr)
			cd.kvCache.Store(kv.Key, &KVMeta{
				value: kv.Value,
				Plan:  plan,
			})
		}
	}

	var foundKey []string
	cd.kvCache.Range(func(key, value interface{}) bool {
		kvKey := key.(string)
		if !existsInPairs(kvNames, kvKey) {
			meta := value.(*KVMeta)
			meta.Plan.Stop()
			foundKey = append(foundKey, kvKey)
		}
		return true
	})

	for _, k := range foundKey {
		cd.kvCache.Delete(k)
	}

}

// existsInPairs
func existsInPairs(kvp api.KVPairs, key string) bool {
	for _, kv := range kvp {
		if kv.Key == key {
			return true
		}
	}
	return false
}

// onKVChanged
func (cd *consulDiscovery) onKVChanged(u uint64, data interface{}) {
	kv, ok := data.(*api.KVPair)
	if !ok {
		return
	}

	if raw, ok := cd.kvCache.Load(kv.Key); ok {
		raw.(*KVMeta).SetValue(kv.Value)
	}
}
