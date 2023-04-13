package event

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/cloudevents/sdk-go/protocol/kafka_sarama/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	"github.com/pkg/errors"
)

// Kafka Kafka
type Kafka struct {
	Cli client.Client
	CloudEventID     string
	CloudEventSource string
	CloudEventType   string
	kafka_sarama.Sender
}

// KafkaClient KafkaClient
var KafkaClient sarama.Client

// NewKafka NewKafka
func NewKafka(topic string) (*Kafka, error) {
	if KafkaClient == nil {
		return nil, errors.New("kafka client is nil")
	}
	sender, err := kafka_sarama.NewSenderFromClient(KafkaClient, topic)
	if err != nil {
		return nil, err
	}

	c, err := cloudevents.NewClient(sender, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		return nil, err
	}
	return &Kafka{Cli: c}, nil
}

// SetCloudEventType SetCloudEventType
func (ke *Kafka) SetCloudEventType(eventName string) {
	ke.CloudEventType = eventName
}

// SetCloudEventID SetCloudEventID
func (ke *Kafka) SetCloudEventID(id string) {
	ke.CloudEventID = id
}

// SetCloudEventSource SetCloudEventSource
func (ke *Kafka) SetCloudEventSource(source string) {
	ke.CloudEventSource = source
}

// Send Send
func (ke *Kafka) Send(ctx context.Context, obj interface{}) error {
	e := cloudevents.NewEvent()
	e.SetID(ke.CloudEventID)
	e.SetType(ke.CloudEventType)
	e.SetSource(ke.CloudEventSource)
	err := e.SetData(cloudevents.ApplicationJSON, obj)
	if err != nil {
		return err
	}
	if result := ke.Cli.Send(
		// Set the producer message key
		kafka_sarama.WithMessageKey(context.Background(), sarama.StringEncoder(e.Type())),
		e,
	); cloudevents.IsUndelivered(result) {
		return errors.Errorf("%+v\n", result)
	}
	return nil
}
