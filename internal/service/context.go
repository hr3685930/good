package service

import (
	"context"
	"good/configs"
)

// Context Context
type Context struct {
	Ctx  context.Context
	Conf configs.Conf
}

// NewContext NewContext
func NewContext(ctx context.Context) *Context {
	return &Context{
		Ctx:  ctx,
		Conf: configs.ENV,
	}
}