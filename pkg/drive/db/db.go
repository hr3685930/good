package db

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
	"time"
)

var (
	// Orm Orm
	Orm *gorm.DB
	// ConnStore ConnStore
	ConnStore sync.Map
	// DefaultLogLevel DefaultLogLevel
	DefaultLogLevel = logger.Error
)

// DB DB
type DB interface {
	Connect() (*gorm.DB, error)
}


// GetConnect GetConnect
func GetConnect(con string) *gorm.DB {
	v, ok := ConnStore.Load(con)
	if ok {
		return v.(*gorm.DB)
	}
	return nil
}

// ListenDriveConnectFail ListenDriveConnectFail
func ListenDriveConnectFail(fn func()) {
	ConnStore.Range(func(key, value interface{}) bool {
		k := key.(string)
		d, _ := GetConnect(k).DB()
		go func() {
			for {
				err := d.PingContext(context.Background())
				if err != nil {
					fmt.Println(k + " connect error")
					fn()
				}
				time.Sleep(time.Second * 5)
			}
		}()
		return true
	})
}
