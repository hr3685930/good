package event

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	ce "github.com/cloudevents/sdk-go/v2/event"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/pkg/errors"
)

// HTTP HTTP
type HTTP struct {
	ce.Event
	client.Client
	endpoint string
}

// NewHTTP NewHTTP
func NewHTTP(endpoint string) (*HTTP, error) {
	p, err := cloudevents.NewHTTP()
	if err != nil {
		return nil, err
	}

	c, err := cloudevents.NewClient(p, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		return nil, err
	}

	e := cloudevents.NewEvent()
	return &HTTP{Event: e, Client: c, endpoint: endpoint}, nil
}

// SetCloudEventType SetCloudEventType
func (he *HTTP) SetCloudEventType(topic string) {
	he.Event.SetType(topic)
}

// SetCloudEventID SetCloudEventID
func (he *HTTP) SetCloudEventID(id string) {
	he.Event.SetID(id)
}

// SetCloudEventSource SetCloudEventSource
func (he *HTTP) SetCloudEventSource(source string) {
	he.Event.SetSource(source)
}

// Send Send
func (he *HTTP) Send(ctx context.Context, obj interface{}) error {
	_ = he.Event.SetData(cloudevents.ApplicationJSON, obj)
	ceCTX := cloudevents.ContextWithTarget(ctx, he.endpoint)
	res := he.Client.Send(ceCTX, he.Event)
	if cloudevents.IsUndelivered(res) {
		return errors.Errorf("%+v\n", res)
	}
	var httpResult *cehttp.Result
	cloudevents.ResultAs(res, &httpResult)
	if httpResult.StatusCode != 200 {
		return errors.Errorf("状态码错误: %+v\n", res)
	}
	return nil
}
