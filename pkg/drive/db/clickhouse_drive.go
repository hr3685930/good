package db

import "good/pkg/drive"

// ClickhouseDrive ClickhouseDrive
type ClickhouseDrive struct {
	Dsn string
}

// Connect Connect
func (m ClickhouseDrive) Connect(key string) error {
	if drive.IgnoreErr {
		return nil
	}
	clickhouseDB := NewClickHouse(m.Dsn)
	orm, err := clickhouseDB.Connect()
	if err != nil {
		return err
	}
	ConnStore.Store(key, orm)
	return nil
}

// Register Register
func (m ClickhouseDrive) Register(key string) {
	Orm = GetConnect(key)
}
