package config

import (
	"good/pkg/drive"
	"good/pkg/drive/queue"
	"reflect"
)

// RabbitMQDrive RabbitMQDrive
type RabbitMQDrive struct {
	AmqpURI string
	App     App
}

// Connect Connect
func (m RabbitMQDrive) Connect(key string, options interface{}, app interface{}) error {
	var typeInfo = reflect.TypeOf(options)
	var valInfo = reflect.ValueOf(options)
	for i := 0; i < typeInfo.NumField(); i++ {
		switch typeInfo.Field(i).Name {
		case "AmqpURI":
			m.AmqpURI = valInfo.Field(i).String()
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
	rabbitMQ := queue.NewAMQP(m.AmqpURI)
	queue.QueueStore.Store(key, rabbitMQ)
	return nil
}

// Default Default
func (RabbitMQDrive) Default(key string) {
	drive.Queue = queue.GetQueueDrive(key)
}
