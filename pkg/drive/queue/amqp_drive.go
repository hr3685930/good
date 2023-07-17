package queue

import "good/pkg/drive"

// RabbitMQDrive RabbitMQDrive
type RabbitMQDrive struct {
	AmqpURI string
}

// Connect Connect
func (m RabbitMQDrive) Connect(key string) error {
	if drive.IgnoreErr {
		return nil
	}
	rabbitMQ := NewAMQP(m.AmqpURI)
	QueueStore.Store(key, rabbitMQ)
	return nil
}

// Register Register
func (RabbitMQDrive) Register(key string) {
	MQ = GetQueueDrive(key)
}
