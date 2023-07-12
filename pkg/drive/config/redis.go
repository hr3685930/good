package config

import (
	"good/pkg/drive"
	"good/pkg/drive/cache"
	"good/pkg/drive/cache/redis"
	"reflect"
)

// RedisDrive RedisDrive
type RedisDrive struct {
	Dsn string
	App App
}

// Connect Connect
func (m RedisDrive) Connect(key string, options interface{}, app interface{}) error {
	var typeInfo = reflect.TypeOf(options)
	var valInfo = reflect.ValueOf(options)
	for i := 0; i < typeInfo.NumField(); i++ {
		switch typeInfo.Field(i).Name {
		case "Dsn":
			m.Dsn = valInfo.Field(i).String()
			break
		}
	}

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
	if IgnoreErr {
		return nil
	}
	c, err := redis.New(m.Dsn)
	if err != nil {
		return err
	}
	cache.CacheMap.Store(key, c)
	return nil
}

// Default Default
func (RedisDrive) Default(key string) {
	drive.Cache = cache.GetCache(key)
}
