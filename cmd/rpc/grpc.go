package rpc

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	proto "good/api/proto/pb"
	"good/configs"
	"good/internal/pkg/errs"
	"good/internal/pkg/errs/export"
	"good/pkg/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"time"
)

//NewGrpc NewGrpc
func NewGrpc() error {
	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(recoverFunc),
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			rpc.UnaryTimeoutInterceptor(time.Second*5),
			grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(opentracing.GlobalTracer())),
			rpc.CustomErrInterceptor(export.GRPCErrorReport),
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
		)),
	}
	grpcServer := rpc.NewGrpc(configs.ENV.App.Debug)
	err := grpcServer.Register(opts, func(s *grpc.Server) {
		healthServer := health.NewServer()
		healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
		healthpb.RegisterHealthServer(s, healthServer)

		proto.RegisterEventServer(s, NewEvent())
		reflection.Register(s)
	})

	err = grpcServer.Run(":8081")
	return err
}

// recoverFunc recover custom func
func recoverFunc(p interface{}) (err error) {
	return errs.InternalError(fmt.Sprintf("%v", p))
}
