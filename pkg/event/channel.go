package event

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	ce "github.com/cloudevents/sdk-go/v2/event"
	"github.com/cloudevents/sdk-go/v2/protocol/gochan"
	"github.com/pkg/errors"
)

var channelCE *Channel

// Channel Channel
type Channel struct {
	ce.Event
	client.Client
	context.Context
}

// NewChannel NewChannel
func NewChannel() (*Channel, error) {
	if channelCE != nil {
		return channelCE, nil
	}
	ctx := context.Background()
	c, err := cloudevents.NewClient(gochan.New(), cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		return nil, err
	}
	e := cloudevents.NewEvent()
	channelCE = &Channel{Event: e, Client: c, Context: ctx}
	return channelCE, nil
}

// SetCloudEventType SetCloudEventType
func (ce *Channel) SetCloudEventType(topic string) {
	ce.Event.SetType(topic)
}

// SetCloudEventID SetCloudEventID
func (ce *Channel) SetCloudEventID(id string) {
	ce.Event.SetID(id)
}

// SetCloudEventSource SetCloudEventSource
func (ce *Channel) SetCloudEventSource(source string) {
	ce.Event.SetSource(source)
}

// Send Send
func (ce *Channel) Send(ctx context.Context, obj interface{}) error {
	err := ce.Event.SetData(cloudevents.ApplicationJSON, obj)
	if err != nil {
		return errors.Errorf("%+v\n", err)
	}
	if res := ce.Client.Send(ce.Context, ce.Event); cloudevents.IsUndelivered(res) {
		return errors.Errorf("%+v\n", res)
	}
	return nil
}
