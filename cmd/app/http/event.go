package http

import (
	"good/cmd/app/job"
	"good/internal/logic/http"
	"good/internal/pkg/errs"
	"good/pkg/event"
)

func eventHandler(c *http.Context) error {
	ctx := c.Request.Context()
	e, err := event.NewHTTPReceive(ctx, job.Bus)
	if err != nil {
		return errs.InternalError("event err:" + err.Error())
	}
	e.ServeHTTP(c.Writer, c.Request)
	return nil
}
