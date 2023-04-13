package main

import (
	"context"
	"github.com/aaronjan/hunch"
	"good/cmd/app"
	"good/configs"
)

func main() {
	ctx := context.Background()
	_, err := hunch.Waterfall(
		ctx,
		func(ctx context.Context, n interface{}) (interface{}, error) {
			return nil, app.Config()
		},
		func(ctx context.Context, n interface{}) (interface{}, error) {
			return hunch.All(
				ctx,
				func(ctx context.Context) (interface{}, error) {
					return nil, app.Log()
				},
				func(ctx context.Context) (interface{}, error) {
					return nil, app.Database(ctx, false, configs.ENV.App.Debug)
				},
				func(ctx context.Context) (interface{}, error) {
					return nil, app.Cache(ctx, false)
				},
				func(ctx context.Context) (interface{}, error) {
					return nil, app.Queue(ctx, false)
				},
			)
		},
		func(ctx context.Context, n interface{}) (interface{}, error) {
			return nil, app.APP()
		},
		func(ctx context.Context, n interface{}) (interface{}, error) {
			return nil, app.Signal()
		},
		func(ctx context.Context, n interface{}) (interface{}, error) {
			return nil, app.Command()
		},
		func(ctx context.Context, n interface{}) (interface{}, error) {
			return nil, app.HTTP()
		},
		func(ctx context.Context, n interface{}) (interface{}, error) {
			return nil, app.GRPC()
		},

	)
	if err != nil {
		panic(err)
	}
}
