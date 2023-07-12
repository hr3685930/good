package config

import (
	"good/pkg/drive"
	"good/pkg/drive/queue"
	"reflect"
)

// ChannelDrive ChannelDrive
type ChannelDrive struct {
	App App
}

// Connect Connect
func (m ChannelDrive) Connect(key string, options interface{}, app interface{}) error {
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
	c := queue.NewChannel()
	queue.QueueStore.Store(key, c)
	return nil
}

// Default Default
func (ChannelDrive) Default(key string) {
	drive.Queue = queue.GetQueueDrive(key)
}
