package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

// Postgre Postgre
type Postgre struct {
	dsn   string
}

// NewPostgre NewPostgre
func NewPostgre(dsn string) *Postgre {
	return &Postgre{dsn}
}

// Connect Connect
func (m *Postgre) Connect() (*gorm.DB, error) {
	dsn := m.dsn
	loglevel := DefaultLogLevel
	orm, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
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
