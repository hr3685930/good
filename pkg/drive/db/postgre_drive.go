package db

import "good/pkg/drive"

// PostgreDrive PostgreDrive
type PostgreDrive struct {
	Dsn string
}

// Connect Connect
func (p PostgreDrive) Connect(key string) error {
	if drive.IgnoreErr {
		return nil
	}
	postgreDB := NewPostgre(p.Dsn)
	orm, err := postgreDB.Connect()
	if err != nil {
		return err
	}
	ConnStore.Store(key, orm)
	return nil
}

// Register Register
func (PostgreDrive) Register(key string) {
	Orm = GetConnect(key)
}
