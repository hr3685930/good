package rpc

import (
	"context"
	proto "good/api/proto/pb"
	"good/cmd/job"
	"good/internal/pkg/format"
)

// Event Event
type Event struct {
}

// NewEvent NewEvent
func NewEvent() *Event {
	return &Event{}
}

// Send Send
func (ev *Event) Send(c context.Context, in *proto.CloudEvent) (*proto.Empty, error) {
	e, err := format.FromProto(in)
	if err != nil {
		return nil, err
	}
	err = job.Bus(c, *e)
	if err != nil {
		return nil, err
	}
	return &proto.Empty{}, err
}
