package config

import (
	"good/pkg/drive/queue"
	"reflect"
)

// KafkaDrive KafkaDrive
type KafkaDrive struct {
	Addr string
	App  App
}

// Connect Connect
func (m KafkaDrive) Connect(key string, options interface{}, app interface{}) error {
	var typeInfo = reflect.TypeOf(options)
	var valInfo = reflect.ValueOf(options)
	for i := 0; i < typeInfo.NumField(); i++ {
		switch typeInfo.Field(i).Name {
		case "Addr":
			m.Addr = valInfo.Field(i).String()
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

	kafkaMQ := queue.NewKafka(m.Addr, m.App.Name)
	_, err := kafkaMQ.GetCli()
	if err != nil {
		return err
	}
	queue.QueueStore.Store(key, kafkaMQ)
	return nil
}

// Default Default
func (KafkaDrive) Default(key string) {
	queue.MQ = queue.GetQueueDrive(key)
}
