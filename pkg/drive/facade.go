package drive

import (
	"github.com/Shopify/sarama"
	"good/pkg/drive/cache"
	"good/pkg/drive/cache/redis"
	"good/pkg/drive/queue"
)

// Kafka Kafka
var Kafka sarama.Client

// Redis Redis
var Redis *redis.Redis

// InitFacade InitFacade
func InitFacade() {
	if queue.GetQueueDrive("kafka") != nil {
		KafkaDrive := queue.GetQueueDrive("kafka").(*queue.Kafka)
		Kafka, _ = KafkaDrive.GetCli()
	}

	if cache.GetCache("redis") != nil {
		Redis = cache.GetCache("redis").(*redis.Redis)
	}
}
