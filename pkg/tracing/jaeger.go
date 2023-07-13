package tracing

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport"
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

// TraceClose TraceClose
func TraceClose() {
	if TraceCloser != nil {
		_ = TraceCloser.Close()
	}
}
