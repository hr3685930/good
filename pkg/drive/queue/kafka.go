package queue

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"net"
)

// Kafka Kafka
type Kafka struct {
	Addr            string
	PrefixGroupName string
	Pub             *kafka.Publisher
	Route           *message.Router
}

// NewKafka NewKafka
func NewKafka(addr, prefixGroupName string) *Kafka {
	return &Kafka{Addr: addr, PrefixGroupName: prefixGroupName}
}

// GetCli GetCli
func (k *Kafka) GetCli() (sarama.Client,error) {
	config := sarama.NewConfig()
	config.Version = sarama.V1_1_1_0
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	return sarama.NewClient([]string{k.Addr}, config)
}

// NewPublisher NewPublisher
func (k *Kafka) NewPublisher() error {
	publisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   []string{k.Addr},
			Marshaler: kafka.DefaultMarshaler{},
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	k.Pub = publisher
	return nil
}

// NewRoute NewRoute
func (k *Kafka) NewRoute() error {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return err
	}
	k.Route = router
	return nil
}

// Publish Publish
func (k *Kafka) Publish(topic string, msg []byte) error {
	msgBody := message.NewMessage(watermill.NewUUID(), msg)
	return k.Pub.Publish(topic, msgBody)
}

// Subscribe Subscribe
func (k *Kafka) Subscribe(topic string, fn MsgFunc) {
	saramaSubscriberConfig := kafka.DefaultSaramaSubscriberConfig()
	saramaSubscriberConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	subscriber, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:               []string{k.Addr},
			Unmarshaler:           kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: saramaSubscriberConfig,
			ConsumerGroup:         k.PrefixGroupName + "." + topic,
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		panic(err)
	}

	k.Route.AddNoPublisherHandler(
		topic,
		topic,
		subscriber,
		func(msg *message.Message) error {
			return fn(msg.Context(), msg.Payload)
		},
	)
}

// AddHandler AddHandler
func (k *Kafka) AddHandler(fn RouteFn) {
	fn(k.Route)
}

// RunRoute RunRoute
func (k *Kafka) RunRoute(ctx context.Context) error {
	err := k.Route.Run(ctx)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

// Ping Ping
func (k *Kafka) Ping() error {
	conn, err := net.Dial("tcp", k.Addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}
