package event

import (
	"context"
	"github.com/cloudevents/sdk-go/protocol/kafka_sarama/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	"github.com/cloudevents/sdk-go/v2/protocol"
	"github.com/hashicorp/go-uuid"
	"github.com/pkg/errors"
	"log"
)

// DefaultSource DefaultSource
const DefaultSource = "https://good/pkg/event/sender"

// CEfn CEfn
type CEfn func(ctx context.Context, event cloudevents.Event) protocol.Result

// Event Event
type Event struct {
	CloudEvent
}

// CloudEvent CloudEvent
type CloudEvent interface {
	SetCloudEventID(id string)
	SetCloudEventType(topic string)
	SetCloudEventSource(source string)
	Send(ctx context.Context, msg interface{}) error
}

// NewHTTPEvent NewHTTPEvent
func NewHTTPEvent(endpoint string, eventName string) *Event {
	httpEvent, err := NewHTTP(endpoint)
	if err != nil {
		panic(err)
	}
	UUID, err := uuid.GenerateUUID()
	if err != nil {
		panic(err)
	}
	httpEvent.SetCloudEventID(UUID)
	httpEvent.SetCloudEventType(eventName)
	httpEvent.SetCloudEventSource(DefaultSource)
	return &Event{CloudEvent: httpEvent}
}

// NewHTTPReceive NewHTTPReceive
func NewHTTPReceive(ctx context.Context, fn CEfn) (*client.EventReceiver, error) {
	p, err := cloudevents.NewHTTP()
	if err != nil {
		return nil, errors.Errorf("%+v\n", err)
	}

	h, err := cloudevents.NewHTTPReceiveHandler(ctx, p, fn)
	if err != nil {
		return nil, errors.Errorf("%+v\n", err)
	}

	return h, nil
}

// NewKafkaEvent NewKafkaEvent
func NewKafkaEvent(topic string, eventName string) *Event {
	kafkaEvent, err := NewKafka(topic)
	if err != nil {
		panic(err)
	}
	UUID, err := uuid.GenerateUUID()
	if err != nil {
		panic(err)
	}
	kafkaEvent.SetCloudEventID(UUID)
	kafkaEvent.SetCloudEventType(eventName)
	kafkaEvent.SetCloudEventSource(DefaultSource)
	return &Event{CloudEvent: kafkaEvent}
}

// NewKafkaReceiver NewKafkaReceiver
func NewKafkaReceiver(ctx context.Context, topic, group string, fn CEfn) error {
	consumer := kafka_sarama.NewConsumerFromClient(KafkaClient, group, topic)
	c, err := cloudevents.NewClient(consumer)
	if err != nil {
		return err
	}

	err = c.StartReceiver(ctx, fn)
	if err != nil {
		return err
	}
	return nil
}

// NewChannelEvent NewChannelEvent
func NewChannelEvent(eventName string) *Event {
	ch, err := NewChannel()
	if err != nil {
		panic(err)
	}
	UUID, err := uuid.GenerateUUID()
	if err != nil {
		panic(err)
	}
	ch.SetCloudEventID(UUID)
	ch.SetCloudEventType(eventName)
	ch.SetCloudEventSource(DefaultSource)
	return &Event{CloudEvent: ch}
}

// NewChanReceive NewChanReceive
func NewChanReceive(fn CEfn) error {
	ch, err := NewChannel()
	if err != nil {
		return err
	}
	// Start the receiver
	go func() {
		if err := ch.Client.StartReceiver(ch.Context, fn); err != nil && err.Error() != "context deadline exceeded" {
			log.Fatalf("[receiver] channel event listen stop error: %s", err)
		}
		log.Println("channel event listen stop")
	}()

	return nil
}

// NewRPCEvent NewRPCEvent
func NewRPCEvent(endpoint, eventName string) *Event {
	r := NewRPC(endpoint)
	UUID, err := uuid.GenerateUUID()
	if err != nil {
		panic(err)
	}
	r.SetCloudEventID(UUID)
	r.SetCloudEventType(eventName)
	r.SetCloudEventSource(DefaultSource)
	return &Event{CloudEvent: r}
}
