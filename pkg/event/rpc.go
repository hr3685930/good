package event

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	ce "github.com/cloudevents/sdk-go/v2/event"
)

// SendFn SendFn
var SendFn = func(ctx context.Context, msg interface{}, endpoint string, event ce.Event) error {
	return nil
}

// RPC RPC
type RPC struct {
	ce.Event
	endpoint string
}

// NewRPC NewRPC
func NewRPC(endpoint string) *RPC {
	e := cloudevents.NewEvent()
	return &RPC{Event: e, endpoint: endpoint}
}

// SetCloudEventType SetCloudEventType
func (he *RPC) SetCloudEventType(topic string) {
	he.Event.SetType(topic)
}

// SetCloudEventID SetCloudEventID
func (he *RPC) SetCloudEventID(id string) {
	he.Event.SetID(id)
}

// SetCloudEventSource SetCloudEventSource
func (he *RPC) SetCloudEventSource(source string) {
	he.Event.SetSource(source)
}

// Send Send
func (he *RPC) Send(ctx context.Context, obj interface{}) error {
	return SendFn(ctx, obj, he.endpoint, he.Event)
}
