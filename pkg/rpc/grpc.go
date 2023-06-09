package rpc

import (
	"context"
	"fmt"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/pkg/errors"
	sentinel "github.com/sentinel-group/sentinel-go-adapters/grpc"
	"good/pkg/errs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

// Grpc Grpc
type Grpc struct {
	G     *grpc.Server
	Debug bool
}

// NewGrpc NewGrpc
func NewGrpc(debug bool) *Grpc {
	return &Grpc{Debug: debug}
}

// Register Register
func (g *Grpc) Register(opts []grpc.ServerOption, register func(s *grpc.Server)) error {
	s := grpc.NewServer(opts...)
	register(s)
	g.G = s
	return nil
}

// Run Run
func (g *Grpc) Run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	fmt.Println("grpc connect success! listen address: " + addr)
	return g.G.Serve(lis)
}

// GrpcError GrpcError
type GrpcError struct {
	Err   error
	Stack []byte `json:"-"`
}

// Error Error
func (h *GrpcError) Error() string {
	return h.Err.Error()
}

// GetStack GetStack
func (h *GrpcError) GetStack() string {
	return string(h.Stack)
}

// ErrorReport ErrorReport
type ErrorReport func(md metadata.MD, req interface{}, stack string, resp *status.Status)

// CustomErrInterceptor CustomErrInterceptor
func CustomErrInterceptor(errReport ErrorReport) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		if err != nil {
			md, _ := metadata.FromIncomingContext(ctx)
			err = ErrorHandler(md, req, err, errReport)
		}
		return
	}
}

// ErrorHandler ErrorHandler
func ErrorHandler(md metadata.MD, req interface{}, err error, errReport ErrorReport) error {
	var stack string
	if e, ok := err.(*GrpcError); ok {
		stack = string(e.Stack)
		err = e.Err
	} else if errPkg, okk := err.(interface{ errs.Error }); okk {
		stack = errPkg.GetStack()
	} else {
		stack = string(debug.Stack())
	}

	s := status.Convert(err)
	errReport(md, req, stack, s)
	return status.Error(s.Code(), s.Message())
}

// Err Err
func Err(code codes.Code, msg string) *GrpcError {
	return &GrpcError{
		Err:   status.New(code, msg).Err(),
		Stack: []byte(fmt.Sprintf("%+v\n", errors.New(msg))),
	}
}

// UnaryTimeoutInterceptor returns a func that sets timeout to incoming unary requests.
func UnaryTimeoutInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var resp interface{}
		var err error
		var lock sync.Mutex
		done := make(chan struct{})
		// create channel with buffer size 1 to avoid goroutine leak
		panicChan := make(chan interface{}, 1)
		go func() {
			defer func() {
				if p := recover(); p != nil {
					// attach call stack to avoid missing in different goroutine
					panicChan <- fmt.Sprintf("%+v\n\n%s", p, strings.TrimSpace(string(debug.Stack())))
				}
			}()

			lock.Lock()
			defer lock.Unlock()
			resp, err = handler(ctx, req)
			close(done)
		}()

		select {
		case p := <-panicChan:
			panic(p)
		case <-done:
			lock.Lock()
			defer lock.Unlock()
			return resp, err
		case <-ctx.Done():
			err := ctx.Err()

			if err == context.Canceled {
				err = status.Error(codes.Canceled, err.Error())
			} else if err == context.DeadlineExceeded {
				err = status.Error(codes.DeadlineExceeded, err.Error())
			}
			return nil, err
		}
	}
}

// UnaryGovernanceClientInterceptor UnaryGovernanceClientInterceptor
func UnaryGovernanceClientInterceptor(err error) grpc.UnaryClientInterceptor {
	return sentinel.NewUnaryClientInterceptor(
		sentinel.WithUnaryClientBlockFallback(func(ctx context.Context, s string, i interface{}, conn *grpc.ClientConn, blockError *base.BlockError) error {
			return err
		}),
	)
}

// UnaryGovernanceServerInterceptor UnaryGovernanceServerInterceptor
func UnaryGovernanceServerInterceptor(err error) grpc.UnaryServerInterceptor {
	return sentinel.NewUnaryServerInterceptor(
		sentinel.WithUnaryServerBlockFallback(func(ctx context.Context, i interface{}, info *grpc.UnaryServerInfo, blockError *base.BlockError) (interface{}, error) {
			return nil, err
		}),
	)
}
