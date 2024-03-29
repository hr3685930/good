package sync

import (
	"good/pkg/drive/cache"
)

// Drive Drive
type Drive struct {}

// Connect Connect
func (m Drive) Connect(key string) error {
	c := New()
	cache.CacheMap.Store(key, c)
	return nil
}

// Register Register
func (Drive) Register(key string) {
	cache.Cached = cache.GetCache(key)
}