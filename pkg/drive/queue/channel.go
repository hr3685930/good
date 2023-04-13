package queue

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

// Channel Channel
type Channel struct {
	*gochannel.GoChannel
	*message.Router
}

// NewChannel NewChannel
func NewChannel() *Channel {
	return &Channel{}
}

// NewPublisher NewPublisher
func (c *Channel) NewPublisher()  error {
	pubSub := gochannel.NewGoChannel(
		gochannel.Config{},
		watermill.NewStdLogger(false, false),
	)
	c.GoChannel = pubSub
	return nil
}

// NewRoute NewRoute
func (c *Channel) NewRoute() error {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return err
	}
	c.Router = router
	return  nil
}

// Publish Publish
func (c *Channel) Publish(topic string, msg []byte) error {
	msgBody := message.NewMessage(watermill.NewUUID(), msg)
	return c.GoChannel.Publish(topic, msgBody)
}

// AddHandler AddHandler
func (c *Channel) AddHandler(fn RouteFn) {
	fn(c.Router)
}

// Subscribe Subscribe
func (c *Channel) Subscribe(topic string, fn MsgFunc) {
	c.Router.AddNoPublisherHandler(
		topic,
		topic,
		c.GoChannel,
		func(msg *message.Message) error {
			return fn(msg.Context(), msg.Payload)
		},
	)
}

// RunRoute RunRoute
func (c *Channel) RunRoute(ctx context.Context) error {
	err := c.Router.Run(ctx)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

// Ping Ping
func (c *Channel) Ping() error {
	return nil
}
