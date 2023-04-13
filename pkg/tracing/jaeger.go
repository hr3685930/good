package tracing

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport"
	"good/configs"
	"good/pkg/drive/cache"
	"io"
)

// TraceCloser TraceCloser
var TraceCloser io.Closer

// NewJaegerTracer NewJaegerTracer
func NewJaegerTracer(service, jaegerHostPort string) io.Closer {
	sender := transport.NewHTTPTransport(
		jaegerHostPort,
	)
	tracer, closer := jaeger.NewTracer(service,
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(sender))

	opentracing.SetGlobalTracer(tracer)
	return closer
}

//NewTrace NewTrace
func NewTrace() error {
	TraceCloser = NewJaegerTracer(configs.ENV.App.Name, configs.ENV.Trace.Endpoint)
	cache.AddTracingHook()
	return nil
}


// TraceClose TraceClose
func TraceClose() {
	if TraceCloser != nil {
		_ = TraceCloser.Close()
	}
}
