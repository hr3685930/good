package queue

import "good/pkg/drive"

// KafkaPrefixGroupName KafkaPrefixGroupName
var KafkaPrefixGroupName = ""

// KafkaDrive KafkaDrive
type KafkaDrive struct {
	Addr string `default:"127.0.0.1:9092" mapstructure:"addr"`
}

// Connect Connect
func (m KafkaDrive) Connect(key string) error {
	if drive.IgnoreErr {
		return nil
	}
	kafkaMQ := NewKafka(m.Addr, KafkaPrefixGroupName)
	_, err := kafkaMQ.GetCli()
	if err != nil {
		return err
	}
	QueueStore.Store(key, kafkaMQ)
	return nil
}

// Default Default
func (KafkaDrive) Default(key string) {
	MQ = GetQueueDrive(key)
}
