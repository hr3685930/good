package queue

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"sync"
)

var (
	// MQ instance
	MQ Queue
	// QueueStore QueueStore
	QueueStore sync.Map
	// logger logger
	logger = watermill.NewStdLogger(false, true)
)

type (
	// MsgFunc MsgFunc
	MsgFunc func(ctx context.Context, msg []byte) error
	// RouteFn RouteFn
	RouteFn func(route *message.Router)
)

// GetQueueDrive GetQueueDrive
func GetQueueDrive(c string) Queue {
	v, ok := QueueStore.Load(c)
	if ok {
		return v.(Queue)
	}
	return nil
}

// Queue Queue
type Queue interface {
	// NewPublisher NewPublisher
	NewPublisher() error
	// NewRoute NewRoute
	NewRoute() error
	// Publish pub
	Publish(topic string, message []byte) error
	// Subscribe Subscribe
	Subscribe(topic string, fn MsgFunc)
	// AddHandler AddHandler
	AddHandler(fn RouteFn)
	// RunRoute Run route
	RunRoute(ctx context.Context) error
	// Ping Ping
	Ping() error
}
