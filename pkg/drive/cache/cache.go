package cache

import (
    "context"
    "sync"
    "time"
)

// CacheMap CacheMap
var CacheMap sync.Map

// Cached Cached
var Cached Cache

// GetCache GetCache
func GetCache(c string) Cache {
    v, ok := CacheMap.Load(c);
    if ok {
        return v.(Cache)
    }
    return nil
}

// AddTracingHook AddTracingHook
func AddTracingHook() {
    Cached.AddTracingHook()
}

// AddMetricHook AddMetricHook
func AddMetricHook() {
    Cached.AddMetricHook()
}

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

