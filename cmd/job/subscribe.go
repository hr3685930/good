package job

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/urfave/cli"
	"good/configs"
	"good/internal/pkg/errs/export"
	"good/pkg/drive/queue"
	"good/pkg/goo"
)

// NewQueueJob NewQueueJob
func NewQueueJob() error  {
	err := queue.MQ.NewPublisher()
	if err != nil {
		return err
	}

	if configs.ENV.Queue.Default == "channel" {
		go SubscribeAll()
	}
	return nil
}

// Queue Queue
func Queue(c *cli.Context) {
	topic := c.String("topic")
	if _, ok := TopicListener[topic]; !ok {
		panic("topic not found")
	}
	Subscribe(topic)
}

// Subscribe Subscribe
func Subscribe(topic string) {
	err := queue.MQ.NewRoute()
	if err != nil {
		panic(err.Error())
	}
	queue.MQ.Subscribe(topic, func(ctx context.Context, msg []byte) error {
		var g goo.Group
		for _, lis := range TopicListener[topic] {
			listen := lis
			g.One(ctx, func(ctx context.Context) error {
				if err = listen(ctx, msg); err != nil {
					return err
				}
				return nil
			})
		}
		err := g.Wait()
		if err != nil {
			export.JobErrorReport(err, "queue", msg)
		}
		return nil
	})

	if err = queue.MQ.RunRoute(context.Background()); err != nil {
		panic(err)
	}
}

// SubscribeAll SubscribeAll
func SubscribeAll() {
	err := queue.MQ.NewRoute()
	if err != nil {
		panic(err.Error())
	}
	queue.MQ.AddHandler(func(route *message.Router) {
		route.AddMiddleware(
			middleware.CorrelationID,
			middleware.Recoverer,
			middleware.InstantAck,
		)
	})

	for topic, listeners := range TopicListener {
		queue.MQ.Subscribe(topic, func(ctx context.Context, msg []byte) error {
			var g goo.Group
			for _, lis := range listeners {
				listen := lis
				g.One(ctx, func(ctx context.Context) error {
					if err = listen(ctx, msg); err != nil {
						return err
					}
					return  nil
				})
			}
			err := g.Wait()
			if err != nil {
				export.JobErrorReport(err, "queue", msg)
			}
			return nil
		})
	}

	if err = queue.MQ.RunRoute(context.Background()); err != nil {
		panic(err)
	}
}
