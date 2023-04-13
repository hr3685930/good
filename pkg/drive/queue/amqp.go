package queue

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	amqpNet "github.com/streadway/amqp"
)

// AMQP AMQP
type AMQP struct {
	URI   string
	Pub   *amqp.Publisher
	Route *message.Router
}

// NewAMQP NewAMQP
func NewAMQP(URI string) *AMQP {
	return &AMQP{URI: URI}
}

// NewPublisher NewPublisher
func (A *AMQP) NewPublisher() error {
	amqpConfig := amqp.NewDurableQueueConfig(A.URI)
	publisher, err := amqp.NewPublisher(
		amqpConfig,
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		panic(err)
	}

	A.Pub = publisher
	return nil
}

// NewRoute NewRoute
func (A *AMQP) NewRoute() error {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return err
	}
	A.Route = router
	return nil
}

// Publish Publish
func (A *AMQP) Publish(topic string, msg []byte) error {
	msgBody := message.NewMessage(watermill.NewUUID(), msg)
	return A.Pub.Publish(topic, msgBody)
}

// Subscribe Subscribe
func (A *AMQP) Subscribe(topic string, fn MsgFunc) {
	amqpConfig := amqp.NewDurableQueueConfig(A.URI)
	subscriber, err := amqp.NewSubscriber(
		amqpConfig,
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		panic(err)
	}
	A.Route.AddNoPublisherHandler(
		topic,
		topic,
		subscriber,
		func(msg *message.Message) error {
			return fn(msg.Context(), msg.Payload)
		},
	)
}

// AddHandler AddHandler
func (A *AMQP) AddHandler(fn RouteFn) {
	fn(A.Route)
}

// RunRoute RunRoute
func (A *AMQP) RunRoute(ctx context.Context) error {
	err := A.Route.Run(ctx)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

// Ping Ping
func (A *AMQP) Ping() error {
	conn, err := amqpNet.Dial(A.URI)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}
