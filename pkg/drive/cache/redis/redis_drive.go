package redis

import (
	"good/pkg/drive"
	"good/pkg/drive/cache"
)

// Drive Drive
type Drive struct {
	Dsn string `default:"" mapstructure:"dsn"`
}

// Connect Connect
func (m Drive) Connect(key string) error {
	if drive.IgnoreErr {
		return nil
	}
	c, err := New(m.Dsn)
	if err != nil {
		return err
	}
	cache.CacheMap.Store(key, c)
	return nil
}

// Register Register
func (Drive) Register(key string) {
	cache.Cached = cache.GetCache(key)
}
