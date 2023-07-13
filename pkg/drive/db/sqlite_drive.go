package db

// SQLiteDrive SQLiteDrive
type SQLiteDrive struct {}

// Connect Connect
func (s SQLiteDrive) Connect(key string) error {
	sqliteDB := NewSqlite()
	orm, err := sqliteDB.Connect()
	if err != nil {
		return err
	}
	ConnStore.Store(key, orm)
	return nil
}

// Default Default
func (SQLiteDrive) Default(key string) {
	Orm = GetConnect(key)
}
