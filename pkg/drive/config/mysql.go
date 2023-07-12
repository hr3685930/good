package config

import (
	"good/pkg/drive"
	"good/pkg/drive/db"
	"reflect"
)

// MYSQLDrive MYSQLDrive
type MYSQLDrive struct {
	Dsn string
	App App
}

// Connect Connect
func (m MYSQLDrive) Connect(key string, options interface{}, app interface{}) error {
	var typeInfo = reflect.TypeOf(options)
	var valInfo = reflect.ValueOf(options)
	num := typeInfo.NumField()
	for i := 0; i < num; i++ {
		switch typeInfo.Field(i).Name {
		case "Dsn":
			m.Dsn = valInfo.Field(i).String()
			break
		}
	}

	var appTypeInfo = reflect.TypeOf(app)
	var appValInfo = reflect.ValueOf(app)
	for i := 0; i < appTypeInfo.NumField(); i++ {
		switch appTypeInfo.Field(i).Name {
		case "Name":
			m.App.Name = appValInfo.Field(i).String()
			break
		case "Env":
			m.App.Env = appValInfo.Field(i).String()
			break
		case "Debug":
			m.App.Debug = appValInfo.Field(i).Bool()
			break
		}
	}

	if m.App.Env == "testing" {
		return nil
	}
	mysqlDB := db.NewMySQL(m.Dsn, m.App.Debug)
	orm, err := mysqlDB.Connect()
	if err != nil {
		return err
	}
	db.ConnStore.Store(key, orm)
	return nil
}

// Default Default
func (MYSQLDrive) Default(key string) {
	drive.Orm = db.GetConnect(key)
}
