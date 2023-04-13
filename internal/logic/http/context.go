package http

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"good/internal/pkg/errs"
	"good/internal/service"
	"net/http"
)

// Context Context
type Context struct {
	*gin.Context
	svc *service.Context
}

// Ping Ping
func (c *Context) Ping() error {
	c.Abort()
	c.Context.String(http.StatusOK, "pong")
	return nil
}

// Response implement your context method like below
func (c *Context) Response(body interface{}) error {
	c.AbortWithStatusJSON(http.StatusOK, body)
	return nil
}

// Bind Bind
func (c *Context) Bind(d interface{}, bindings ...binding.Binding) error {
	var err error
	for i := range bindings {
		switch bindings[i] {
		case binding.JSON:
			err = c.ShouldBindWith(d, binding.JSON)
		case binding.XML:
			err = c.ShouldBindWith(d, binding.XML)
		case binding.Form:
			err = c.ShouldBindWith(d, binding.Form)
		case binding.Query:
			err = c.ShouldBindWith(d, binding.Query)
		case binding.FormPost:
			err = c.ShouldBindWith(d, binding.FormPost)
		case binding.FormMultipart:
			err = c.ShouldBindWith(d, binding.FormMultipart)
		case binding.ProtoBuf:
			err = c.ShouldBindWith(d, binding.ProtoBuf)
		case binding.MsgPack:
			err = c.ShouldBindWith(d, binding.MsgPack)
		case binding.YAML:
			err = c.ShouldBindWith(d, binding.YAML)
		case binding.Header:
			err = c.ShouldBindWith(d, binding.Header)
		default:
			err = c.ShouldBindUri(d)
		}
		if err != nil {
			return errs.ValidationFailed("参数错误:" + err.Error())
		}
	}
	return nil
}


// HandlerFunc HandlerFunc
type HandlerFunc func(c *Context) error

// Router Router
type Router struct {
	*gin.RouterGroup
}

// NewRouter NewRouter
func NewRouter(e *gin.Engine) *Router {
	return &Router{
		RouterGroup: &e.RouterGroup,
	}
}

// GET GET
func (r *Router) GET(relativePath string, handlers ...HandlerFunc) *Router {
	r.wrapRoute(http.MethodGet, relativePath, handlers...)
	return r
}

// POST POST
func (r *Router) POST(relativePath string, handlers ...HandlerFunc) *Router {
	r.wrapRoute(http.MethodPost, relativePath, handlers...)
	return r
}

// PUT PUT
func (r *Router) PUT(relativePath string, handlers ...HandlerFunc) *Router {
	r.wrapRoute(http.MethodPut, relativePath, handlers...)
	return r
}

// PATCH PATCH
func (r *Router) PATCH(relativePath string, handlers ...HandlerFunc) *Router {
	r.wrapRoute(http.MethodPatch, relativePath, handlers...)
	return r
}

// HEAD HEAD
func (r *Router) HEAD(relativePath string, handlers ...HandlerFunc) *Router {
	r.wrapRoute(http.MethodHead, relativePath, handlers...)
	return r
}

// OPTIONS OPTIONS
func (r *Router) OPTIONS(relativePath string, handlers ...HandlerFunc) *Router {
	r.wrapRoute(http.MethodOptions, relativePath, handlers...)
	return r
}

// DELETE DELETE
func (r *Router) DELETE(relativePath string, handlers ...HandlerFunc) *Router {
	r.wrapRoute(http.MethodDelete, relativePath, handlers...)
	return r
}

// CONNECT CONNECT
func (r *Router) CONNECT(relativePath string, handlers ...HandlerFunc) *Router {
	r.wrapRoute(http.MethodConnect, relativePath, handlers...)
	return r
}

// TRACE TRACE
func (r *Router) TRACE(relativePath string, handlers ...HandlerFunc) *Router {
	r.wrapRoute(http.MethodTrace, relativePath, handlers...)
	return r
}

// Use Use
func (r *Router) Use(handlers ...HandlerFunc) *Router {
	r.wrapRoute("use", "", handlers...)
	return r
}

// Group Group
func (r *Router) Group(relativePath string, handlers ...HandlerFunc) *Router {
	g := r.wrapRoute("group", relativePath, handlers...).(*gin.RouterGroup)
	return &Router{
		RouterGroup: g,
	}
}

// wrapRoute wrapRoute
func (r *Router) wrapRoute(method string, relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	hds := make([]gin.HandlerFunc, 0, len(handlers))
	for _, hd := range handlers {
		hds = append(hds, wrapHandler(hd))
	}
	switch method {
	case http.MethodGet:
		return r.RouterGroup.GET(relativePath, hds...)
	case http.MethodPost:
		return r.RouterGroup.POST(relativePath, hds...)
	case http.MethodPut:
		return r.RouterGroup.PUT(relativePath, hds...)
	case http.MethodPatch:
		return r.RouterGroup.PATCH(relativePath, hds...)
	case http.MethodHead:
		return r.RouterGroup.HEAD(relativePath, hds...)
	case http.MethodOptions:
		return r.RouterGroup.OPTIONS(relativePath, hds...)
	case http.MethodDelete:
		return r.RouterGroup.DELETE(relativePath, hds...)
	case http.MethodConnect:
		return r.RouterGroup.Handle(http.MethodConnect, relativePath, hds...)
	case "use":
		return r.RouterGroup.Use(hds...)
	case "group":
		return r.RouterGroup.Group(relativePath, hds...)
	}
	return r.RouterGroup.Handle(http.MethodTrace, relativePath, hds...)
}

// wrapHandler wrapHandler
func wrapHandler(hd HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := hd(&Context{Context: ctx, svc: service.NewContext(ctx)})
		if err != nil {
			ctx.Abort()
			_ = ctx.Error(err)
		}
	}
}
