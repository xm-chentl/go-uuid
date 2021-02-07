package uuid

import "sync"

const DEFAULT = "default"

var (
	mt        sync.Mutex
	keyOfUUID = make(map[string]IUUID)
)

func SetDefault(inst IUUID) {
	mt.Lock()
	defer mt.Unlock()
	if inst != nil {
		deyOfUUID[DEFAULT] = inst
	}
}

func Default() IUUID {
	mt.Lock()
	defer mt.Unlock()
	
	if inst, ok := keyOfUUID[DEFAULT]; ok {
		return inst
	}
	panic("default instalce is not exists")
}

func Set(key string, inst IUUID) {
	mt.Lock()
	defer mt.Unlock()

	if key != "" && inst != nil {
		keyOfUUID[key] = inst
	}
}

func Get(key string) (IUUID, error) {
	mt.Lock()
	defer mt.Unlock()

	if inst, ok := keyOfUUID[key]; ok {
		return inst
	}
	panic(
		fmt.Sprintf("key (%s) instalce is not exists", key)
	)
}
