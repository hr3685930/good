package job

import (
	"context"
	"good/internal/logic/job"
)

// Listener Listener
type Listener func(ctx context.Context, msg []byte) error

// EventListeners EventListeners
var EventListeners = map[string][]Listener{
	"com.example": {
		job.Example,
	},
}

// TopicListener TopicListener
var TopicListener = map[string][]Listener{
	"example.topic": {
		job.Example,
		job.Examples,
	},
}
