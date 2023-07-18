package middleware

import (
	"good/internal/logic/http"
	"good/internal/pkg/errs"
)

// Auth 认证
func Auth(c *http.Context) error {
	openID := c.Request.Header.Get("X-Consumer-Custom-ID")
	if openID == "" {
		return errs.AuthorizationFailed("authorization failed")
	}
	c.Next()
	return nil
}
