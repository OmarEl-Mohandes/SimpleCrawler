package DataStructures

import "sync"


type ConcurrentCache struct {
	sync.RWMutex
	Cache map[string]bool
}

func NewConcurrentCache() *ConcurrentCache {
	return &ConcurrentCache{
		Cache: make(map[string]bool),
	}
}

func (mp *ConcurrentCache) Store(key string) {
	mp.Lock()
	defer mp.Unlock()
	mp.Cache[key] = true
}

func (mp *ConcurrentCache) Load(key string)(bool, bool){
	mp.RLock()
	defer mp.RUnlock()
	res , ok := mp.Cache[key]
	return res, ok
}
