package job

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	ce "github.com/cloudevents/sdk-go/v2/event"
	"github.com/cloudevents/sdk-go/v2/protocol"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/urfave/cli"
	proto "good/api/proto/pb"
	"good/configs"
	"good/internal/pkg/errs"
	"good/internal/pkg/errs/export"
	"good/internal/pkg/format"
	"good/pkg/drive/queue"
	"good/pkg/event"
	"good/pkg/goo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

// KafkaEventTopic KafkaEventTopic
var KafkaEventTopic = "event"

// KafkaEventListen KafkaEventListen
func KafkaEventListen(c *cli.Context) {
	err := event.NewKafkaReceiver(context.Background(), KafkaEventTopic, configs.ENV.App.Name, Bus)
	if err != nil {
		panic(err)
	}
}

// NewEventReceive NewEventReceive
func NewEventReceive() error {
	event.SendFn = RPCSend
	event.KafkaClient = GetKafkaCli()
	return event.NewChanReceive(Bus)
}

// RPCSend RpcSend
func RPCSend(ctx context.Context, obj interface{}, endpoint string, cloudevent ce.Event) error {
	conn, err := grpc.DialContext(ctx,
		endpoint,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			MinConnectTimeout: time.Second * 5,
		}),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			grpc_opentracing.UnaryClientInterceptor(),
		)),
	)

	if err != nil {
		return errs.InternalError("connect error")
	}
	defer conn.Close()
	msg, err := json.Marshal(obj)
	if err != nil {
		return errs.InternalError("Marshal error")
	}
	err = cloudevent.SetData(format.ContentTypeProtobuf, msg)
	if err != nil {
		return errs.InternalError("set event data error")
	}
	client := proto.NewEventClient(conn)
	req, err := format.ToProto(&cloudevent)
	if err != nil {
		return errs.InternalError("event to proto error")
	}
	_, err = client.Send(ctx, req)
	return err
}

// Bus event bus
func Bus(ctx context.Context, e cloudevents.Event) protocol.Result {
	if _, ok := EventListeners[e.Type()]; !ok {
		return errs.ResourceNotFound("event not found")
	}
	var g goo.Group
	for _, lis := range EventListeners[e.Type()] {
		listen := lis
		g.One(ctx, func(ctx context.Context) error {
			err := listen(ctx, e.DataEncoded)
			if err != nil {
				return err
			}
			return nil
		})
	}
	err := g.Wait()
	if err != nil {
		export.JobErrorReport(err, "event", e.DataEncoded)
		return err
	}
	return nil
}

// GetKafkaCli GetKafkaCli
func GetKafkaCli() sarama.Client {
	if queue.GetQueueDrive("kafka") != nil {
		KafkaDrive := queue.GetQueueDrive("kafka").(*queue.Kafka)
		KafkaCli, _ := KafkaDrive.GetCli()
		return KafkaCli
	}
	return nil
}
