package dict

import (
	"sync"
)

type SyncDict struct {
	dicMap sync.Map
}

func MakeSyncDict() *SyncDict {
	return &SyncDict{}
}

func (dict *SyncDict) Get(key string) (val interface{}, exists bool) {
	val, exists = dict.dicMap.Load(key)
	return
}

func (dict *SyncDict) Len() int {
	length := 0
	dict.dicMap.Range(func(key, value any) bool {
		length++
		return true
	})
	return length
}

func (dict *SyncDict) Put(key string, val interface{}) (result int) {
	_, exists := dict.dicMap.Load(key)
	dict.dicMap.Store(key, val)
	if exists {
		return 1
	}
	return 0
}

func (dict *SyncDict) PutIfAbsent(key string, val interface{}) (result int) {
	_, exists := dict.dicMap.Load(key)
	if exists {
		return 0
	}
	dict.dicMap.Store(key, val)
	return 1
}

func (dict *SyncDict) PutIfExists(key string, val interface{}) (result int) {
	_, exists := dict.dicMap.Load(key)
	if !exists {
		return 0
	}
	dict.dicMap.Store(key, val)
	return 1
}

func (dict *SyncDict) Remove(key string) (result int) {
	_, exists := dict.dicMap.Load(key)
	if exists {
		dict.dicMap.Delete(key)
		return 1
	}
	return 0
}

func (dict *SyncDict) ForEach(consumer Consumer) {
	dict.dicMap.Range(func(key, value any) bool {
		keyString, ok := key.(string)
		if ok {
			return consumer(keyString, value)
		}
		return false
	})
}

func (dict *SyncDict) Keys() []string {
	keys := make([]string, 0, dict.Len())
	dict.dicMap.Range(func(key any, _ any) bool {
		keyString, ok := key.(string)
		if ok {
			keys = append(keys, keyString)
		}
		return true
	})
	return keys
}

func (dict *SyncDict) RandomKeys(limit int) []string {
	result := make([]string, limit)
	for i := 0; i < limit; i++ {
		dict.dicMap.Range(func(key, value any) bool {
			result[i], _ = key.(string)
			return false
		})
	}
	return result
}

func (dict *SyncDict) RandomDistinctKeys(limit int) []string {
	result := make([]string, limit)
	i := 0
	dict.dicMap.Range(func(key, value any) bool {
		result[i], _ = key.(string)
		i++
		if i != limit {
			return true
		}
		return false
	})
	return result
}

func (dict *SyncDict) Clear() {
	*dict = *MakeSyncDict()
}
