package config

import (
	"good/pkg/drive"
	"good/pkg/drive/db"
	"reflect"
)

// PostgreDrive PostgreDrive
type PostgreDrive struct {
	Dsn string
	App App
}

// Connect Connect
func (p PostgreDrive) Connect(key string, options interface{}, app interface{}) error {
	var typeInfo = reflect.TypeOf(options)
	var valInfo = reflect.ValueOf(options)
	num := typeInfo.NumField()
	for i := 0; i < num; i++ {
		switch typeInfo.Field(i).Name {
		case "Dsn":
			p.Dsn = valInfo.Field(i).String()
			break
		}
	}

	var appTypeInfo = reflect.TypeOf(app)
	var appValInfo = reflect.ValueOf(app)
	for i := 0; i < appTypeInfo.NumField(); i++ {
		switch appTypeInfo.Field(i).Name {
		case "Name":
			p.App.Name = appValInfo.Field(i).String()
			break
		case "Env":
			p.App.Env = appValInfo.Field(i).String()
			break
		case "Debug":
			p.App.Debug = appValInfo.Field(i).Bool()
			break
		}
	}
	if IgnoreErr {
		return nil
	}
	postgreDB := db.NewPostgre(p.Dsn, p.App.Debug)
	orm, err := postgreDB.Connect()
	if err != nil {
		return err
	}
	db.ConnStore.Store(key, orm)
	return nil
}

// Default Default
func (PostgreDrive) Default(key string) {
	drive.Orm = db.GetConnect(key)
}
