package config

import (
	"good/pkg/drive"
	"good/pkg/drive/db"
	"reflect"
)

// SQLiteDrive SQLiteDrive
type SQLiteDrive struct {
	App App
}

// Connect Connect
func (s SQLiteDrive) Connect(key string, options interface{}, app interface{}) error {
	var appTypeInfo = reflect.TypeOf(app)
	var appValInfo = reflect.ValueOf(app)
	for i := 0; i < appTypeInfo.NumField(); i++ {
		switch appTypeInfo.Field(i).Name {
		case "Name":
			s.App.Name = appValInfo.Field(i).String()
			break
		case "Env":
			s.App.Env = appValInfo.Field(i).String()
			break
		case "Debug":
			s.App.Debug = appValInfo.Field(i).Bool()
			break
		}
	}

	sqliteDB := db.NewSqlite(s.App.Debug)
	orm, err := sqliteDB.Connect()
	if err != nil {
		return err
	}
	db.ConnStore.Store(key, orm)
	return nil
}

// Default Default
func (SQLiteDrive) Default(key string) {
	drive.Orm = db.GetConnect(key)
}
