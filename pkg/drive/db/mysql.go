package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

// MySQL MySQL
type MySQL struct {
	dsn   string
}

// NewMySQL NewMySQL
func NewMySQL(dsn string) *MySQL {
	return &MySQL{dsn}
}

// Connect Connect
func (m *MySQL) Connect() (*gorm.DB, error) {
	dsn := m.dsn
	loglevel := DefaultLogLevel
	orm, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(loglevel),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := orm.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return  orm, nil
}