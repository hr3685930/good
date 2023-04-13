package utils

import (
	"github.com/Shopify/sarama"
	"good/pkg/drive/queue"
)

// GetKafkaCli GetKafkaCli
func GetKafkaCli() sarama.Client {
	if queue.GetQueueDrive("kafka") != nil {
		kafkaCli := queue.GetQueueDrive("kafka").(*queue.Kafka)
		cli, _ := kafkaCli.GetCli()
		return cli
	}
	return nil
}
