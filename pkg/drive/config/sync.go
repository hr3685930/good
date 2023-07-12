package config

import (
	"good/pkg/drive"
	"good/pkg/drive/cache"
	"good/pkg/drive/cache/sync"
	"reflect"
)

// SyncDrive SyncDrive
type SyncDrive struct {
	App App
}

// Connect Connect
func (m SyncDrive) Connect(key string, options interface{}, app interface{}) error {
	var appTypeInfo = reflect.TypeOf(app)
	var appValInfo = reflect.ValueOf(app)
	for i := 0; i < appTypeInfo.NumField(); i++ {
		switch appTypeInfo.Field(i).Name {
		case "Name":
			m.App.Name = appValInfo.Field(i).String()
			break
		case "Env":
			m.App.Env = appValInfo.Field(i).String()
			break
		case "Debug":
			m.App.Debug = appValInfo.Field(i).Bool()
			break
		}

	}
	c := sync.New()
	cache.CacheMap.Store(key, c)
	return nil
}

// Default Default
func (SyncDrive) Default(key string) {
	drive.Cache = cache.GetCache(key)
}