package http

import (
	"good/cmd/job"
	"good/internal/logic/http"
	"good/pkg/event"
)

// eventHandler eventHandler
func eventHandler(c *http.Context) error {
	ctx := c.Request.Context()
	e, err := event.NewHTTPReceive(ctx, job.Bus)
	if err != nil {
		return err
	}
	e.ServeHTTP(c.Writer, c.Request)
	return nil
}
