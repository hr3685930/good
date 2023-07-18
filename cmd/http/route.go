package http

import (
	"github.com/gin-gonic/gin"
	"good/internal/logic/http"
	"good/internal/logic/http/middleware"
	"good/internal/pkg/errs"
	"good/internal/pkg/errs/export"
	httpPkg "good/pkg/http"
	"time"
)

// Routes Routes
func Routes(e *http.Router) {
	e.RouterGroup.Use(
		httpPkg.TimeoutMiddleware(time.Second*10),
		httpPkg.ErrHandler(export.HTTPErrorReport),
		gin.CustomRecovery(func(c *gin.Context, err interface{}) {
			_ = c.Error(errs.InternalError("internal error"))
		}),
	)

	e.GET("/ping", func(c *http.Context) error {
		return c.Ping()
	})

	e.POST("/event", eventHandler)
	api := e.Group("/api", middleware.Auth)
	{
		api.GET("/user", http.GetUser)
	}
}