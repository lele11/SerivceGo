package consulsd

import (
	"errors"
	"sync"

	"github.com/cihub/seelog"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/watch"
)

func (self *consulDiscovery) SetValue(key string, dataPtr interface{}) error {
	raw, err := AnyToBytes(dataPtr)
	if err != nil {
		return err
	}
	_, err = self.client.KV().Put(&api.KVPair{
		Key:   key,
		Value: raw,
	}, nil)

	return err
}

func (self *consulDiscovery) GetValue(key string, valuePtr interface{}) error {
	data, err := self.GetRawValue(key)
	if err != nil {
		return err
	}

	return BytesToAny(data, valuePtr)
}

func (self *consulDiscovery) GetRawValue(key string) ([]byte, error) {
	if raw, ok := self.kvCache.Load(key); ok {
		meta := raw.(*KVMeta)
		return meta.Value(), nil
	}
	// cache中没找到直接获取
	kvpair, _, err := self.client.KV().Get(key, nil)
	if err != nil {
		return nil, err
	}
	if kvpair == nil {
		return nil, ErrValueNotExists
	}
	return kvpair.Value, nil
}

func (self *consulDiscovery) DeleteValue(key string) error {
	_, err := self.client.KV().Delete(key, nil)
	return err
}

var (
	ErrValueNotExists = errors.New("value not exists")
)

func (self *consulDiscovery) startWatchKV(prefix string) {
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

	plan.Handler = self.onKVListChanged
	go plan.Run(self.consulAddr)
}

type KVMeta struct {
	value      []byte
	valueGuard sync.RWMutex
	Plan       *watch.Plan
}

func (self *KVMeta) SetValue(v []byte) {
	self.valueGuard.Lock()
	self.value = v
	self.valueGuard.Unlock()
}

func (self *KVMeta) Value() []byte {
	self.valueGuard.RLock()
	defer self.valueGuard.RUnlock()
	return self.value
}

func (self *consulDiscovery) onKVListChanged(u uint64, data interface{}) {
	kvNames, ok := data.(api.KVPairs)
	if !ok {
		return
	}

	for _, kv := range kvNames {
		// 已经在cache里的,肯定添加过watch了
		if _, ok := self.kvCache.Load(kv.Key); ok {
			continue
		}
		plan, err := watch.Parse(map[string]interface{}{
			"type": "key",
			"key":  kv.Key,
		})

		if err == nil {
			plan.Handler = self.onKVChanged
			go plan.Run(self.consulAddr)
			self.kvCache.Store(kv.Key, &KVMeta{
				value: kv.Value,
				Plan:  plan,
			})
		}
	}

	var foundKey []string
	self.kvCache.Range(func(key, value interface{}) bool {
		kvKey := key.(string)
		if !existsInPairs(kvNames, kvKey) {
			meta := value.(*KVMeta)
			meta.Plan.Stop()
			foundKey = append(foundKey, kvKey)
		}
		return true
	})

	for _, k := range foundKey {
		self.kvCache.Delete(k)
	}

}

func existsInPairs(kvp api.KVPairs, key string) bool {
	for _, kv := range kvp {
		if kv.Key == key {
			return true
		}
	}
	return false
}

func (self *consulDiscovery) onKVChanged(u uint64, data interface{}) {
	kv, ok := data.(*api.KVPair)
	if !ok {
		return
	}

	if raw, ok := self.kvCache.Load(kv.Key); ok {
		raw.(*KVMeta).SetValue(kv.Value)
	}
}
