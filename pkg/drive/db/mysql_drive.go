package db

import "good/pkg/drive"

// MYSQLDrive MYSQLDrive
type MYSQLDrive struct {
	Dsn string `default:"" mapstructure:"dsn"`
}

// Connect Connect
func (m MYSQLDrive) Connect(key string) error {
	if drive.IgnoreErr {
		return nil
	}
	mysqlDB := NewMySQL(m.Dsn)
	orm, err := mysqlDB.Connect()
	if err != nil {
		return err
	}
	ConnStore.Store(key, orm)
	return nil
}

// Register Register
func (MYSQLDrive) Register(key string) {
	Orm = GetConnect(key)
}
