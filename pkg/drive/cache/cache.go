package cache

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var (
	// Cached Cached
	Cached Cache
	// CacheMap CacheMap
	CacheMap sync.Map
)

// Cache Cache
type Cache interface {
	// Contains check if a cached key exists
	Contains(ctx context.Context, key string) bool
	// Delete remove the cached key
	Delete(ctx context.Context, key string) error
	// Fetch retrieve the cached key value
	Fetch(ctx context.Context, key string) (string, error)
	// FetchMulti retrieve multiple cached keys value
	FetchMulti(ctx context.Context, keys []string) map[string]string
	// Flush remove all cached keys
	Flush(ctx context.Context) error
	// Save cache a value by key
	Save(ctx context.Context, key string, value string, lifeTime time.Duration) error
	// AddTracingHook Hook
	AddTracingHook()
	// AddMetricHook Hook
	AddMetricHook()
	// Ping ping
	Ping() error
}

// GetCache GetCache
func GetCache(c string) Cache {
	v, ok := CacheMap.Load(c)
	if ok {
		return v.(Cache)
	}
	return nil
}

// ListenDriveConnectFail ListenDriveConnectFail
func ListenDriveConnectFail(fn func()) {
	CacheMap.Range(func(key, value interface{}) bool {
		k := key.(string)
		d := GetCache(k)
		go func() {
			for {
				if d.Ping() != nil {
					fmt.Println(k + " connect error")
					fn()
				}
				time.Sleep(time.Second * 5)
			}
		}()
		return true
	})
}