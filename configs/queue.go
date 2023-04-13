package configs

import "good/pkg/drive/config"

// Queue default
type Queue struct {
	Default string `default:"channel" mapstructure:"default"`
	Kafka   Kafka
	Channel Channel
}

// Kafka kafka config
type Kafka struct {
	config.KafkaDrive
	Addr string `default:"127.0.0.1:9092" mapstructure:"addr"`
}

// Rabbitmq Rabbitmq config
type Rabbitmq struct {
	config.RabbitMQDrive
	AmqpURI string `default:"amqp://guest:guest@rabbitmq:5672/" mapstructure:"amqp_uri"`
}

// Channel Channel
type Channel struct {
	config.ChannelDrive
}
