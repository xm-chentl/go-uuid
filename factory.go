package uuid

import (
	"fmt"
	"sync"
)

// DEFAULT 默认实例名
const DEFAULT = "default"

var (
	mt        sync.Mutex
	keyOfUUID = make(map[string]IUUID)
)

// SetDefault 设置默认实例
func SetDefault(inst IUUID) {
	mt.Lock()
	defer mt.Unlock()
	if inst != nil {
		keyOfUUID[DEFAULT] = inst
	}
}

// Default 默认实例
func Default() IUUID {
	mt.Lock()
	defer mt.Unlock()

	if inst, ok := keyOfUUID[DEFAULT]; ok {
		return inst
	}
	panic("default instalce is not exists")
}

// Set 设置实例
func Set(key string, inst IUUID) {
	mt.Lock()
	defer mt.Unlock()

	if key != "" && inst != nil {
		keyOfUUID[key] = inst
	}
}

// Get 获取实例
func Get(key string) (IUUID, error) {
	mt.Lock()
	defer mt.Unlock()

	if inst, ok := keyOfUUID[key]; ok {
		return inst, nil
	}
	panic(
		fmt.Sprintf("key (%s) instalce is not exists", key),
	)
}
