package pluginIO

import "sync"

type SafeMap struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		data: make(map[string]interface{}),
	}
}

// Set 写操作使用写锁
func (sm *SafeMap) Set(key string, value interface{}) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data[key] = value
}

// Get 读操作使用读锁
func (sm *SafeMap) Get(key string) (interface{}, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	value, exists := sm.data[key]
	return value, exists
}

// Delete 删除操作
func (sm *SafeMap) Delete(key string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.data, key)
}

// Range 遍历操作
func (sm *SafeMap) Range(f func(key string, value interface{}) bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	for k, v := range sm.data {
		if !f(k, v) {
			break
		}
	}
}
