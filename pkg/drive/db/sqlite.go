package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

// SQLite SQLite
type SQLite struct{}

// NewSqlite NewSqlite
func NewSqlite() *SQLite {
	return &SQLite{}
}

// Connect Connect
func (m *SQLite) Connect() (*gorm.DB, error) {
	loglevel := DefaultLogLevel
	orm, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{
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
