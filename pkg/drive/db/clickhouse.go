package db

import (
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

// ClickHouse ClickHouse
type ClickHouse struct {
	dsn   string
}

// NewClickHouse NewClickHouse
func NewClickHouse(dsn string) *ClickHouse {
	return &ClickHouse{dsn}
}

// Connect Connect
func (c *ClickHouse) Connect() (*gorm.DB, error) {
	dsn := c.dsn
	loglevel := DefaultLogLevel
	orm, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(loglevel),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, err
	}
	sqlDB, _ := orm.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return orm, nil
}
