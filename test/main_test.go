package test

import (
	"context"
	"good/cmd"
	"good/cmd/app"
	"good/configs"
	"good/internal/logic/script"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"testing"
)

func TestMain(m *testing.M) {
	config()
	commands.Migrate(nil)
	RegisterGrpc(func(server *grpc.Server) {
		healthServer := health.NewServer()
		healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
		healthpb.RegisterHealthServer(server, healthServer)
	})
	m.Run()
}

func config() {
	ctx := context.Background()
	_ = cmd.Config()
	configs.ENV.App.Name = "test"
	configs.ENV.App.Debug = true
	configs.ENV.App.ErrReport = ""
	configs.ENV.App.Env = "testing"
	configs.ENV.Queue.Default = "channel"
	configs.ENV.Database.Default = "sqlite"
	configs.ENV.Cache.Default = "sync"
	_ = app.Log()
	_ = cmd.Drive(ctx)
	_ = cmd.APP()
}
