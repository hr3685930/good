package main

import (
	"context"
	"github.com/aaronjan/hunch"
	"good/cmd"
)

func main() {
	ctx := context.Background()
	_, err := hunch.Waterfall(
		ctx,
		func(ctx context.Context, n interface{}) (interface{}, error) {
			return nil, cmd.Drive(ctx)
		},
		func(ctx context.Context, n interface{}) (interface{}, error) {
			return nil, cmd.APP()
		},
		func(ctx context.Context, n interface{}) (interface{}, error) {
			return nil, cmd.Command()
		},
		func(ctx context.Context, n interface{}) (interface{}, error) {
			return nil, cmd.HTTP()
		},
		func(ctx context.Context, n interface{}) (interface{}, error) {
			return nil, cmd.GRPC()
		},

	)
	if err != nil {
		panic(err)
	}
}
