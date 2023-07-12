package db

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
)

// ConnStore ConnStore
var ConnStore sync.Map

// DefaultLogLevel DefaultLogLevel
var DefaultLogLevel = logger.Error

// GetConnect GetConnect
func GetConnect(con string) *gorm.DB {
	v, ok := ConnStore.Load(con)
	if ok {
		return v.(*gorm.DB)
	}
	return nil
}

// DB DB
type DB interface {
	Connect() (*gorm.DB, error)
}
