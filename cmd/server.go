package cmd

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronjan/hunch"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"good/cmd/app"
	cmdHTTP "good/cmd/http"
	"good/cmd/job"
	"good/cmd/rpc"
	"good/configs"
	"good/internal/logic/http"
	"good/pkg/drive/config"
	"good/pkg/drive/db"
	httpPkg "good/pkg/http"
	"gorm.io/gorm/logger"
	"os"
	"sort"
	"time"
)

//Config config
func Config() error {
	return config.Load(&configs.ENV)
}

// Drive Drive
func Drive(ctx context.Context) error {
	_, err := hunch.All(
		ctx,
		func(ctx context.Context) (interface{}, error) {
			return nil, NewDatabase(ctx)
		},
		func(ctx context.Context) (interface{}, error) {
			return nil, NewCache(ctx)
		},
		func(ctx context.Context) (interface{}, error) {
			return nil, NewQueue(ctx)
		},
	)

	err = job.NewEventReceive()
	if err != nil {
		return err
	}
	err = job.NewQueueJob()
	if err != nil {
		return err
	}
	return nil
}

//NewDatabase db
func NewDatabase(ctx context.Context) error {
	db.DefaultLogLevel = logger.Silent
	if configs.ENV.App.Debug {
		db.DefaultLogLevel = logger.Info
	}
	_, err := hunch.Retry(ctx, 0, func(c context.Context) (interface{}, error) {
		err := config.Drive(configs.ENV.Database, configs.ENV.App)
		if err != nil {
			fmt.Println("db reconnect...", err)
			time.Sleep(time.Second * 2)
		}
		return nil, err
	})

	return err
}

//NewCache NewCache
func NewCache(ctx context.Context) error {
	_, err := hunch.Retry(ctx, 0, func(c context.Context) (interface{}, error) {
		err := config.Drive(configs.ENV.Cache, configs.ENV.App)
		if err != nil {
			fmt.Println("cache reconnect...", err)
			time.Sleep(time.Second * 2)
		}
		return nil, err
	})

	return err
}

//NewQueue NewQueue
func NewQueue(ctx context.Context) error {
	_, err := hunch.Retry(ctx, 0, func(c context.Context) (interface{}, error) {
		err := config.Drive(configs.ENV.Queue, configs.ENV.App)
		if err != nil {
			fmt.Println("queue reconnect...", err)
			time.Sleep(time.Second * 2)
		}
		return nil, err
	})
	return err
}

// APP APP
func APP() error {
	return app.NewAPP()
}


// Command Command
func Command() error {
	flag.Parse()
	if len(flag.Args()) == 0 {
		return nil
	}
	appCli := cli.NewApp()
	appCli.Commands = Commands
	sort.Sort(cli.FlagsByName(appCli.Flags))
	sort.Sort(cli.CommandsByName(appCli.Commands))
	_ = appCli.Run(os.Args)
	os.Exit(0)
	return nil
}

//HTTP HTTP
func HTTP() error {
	if configs.ENV.App.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	if err := httpPkg.LoadValidatorLocal("zh"); err != nil {
		return err
	}
	gServer := gin.New()
	cmdHTTP.Routes(http.NewRouter(gServer))
	return gServer.Run(":8080")
}

//GRPC GRPC
func GRPC() error {
	return rpc.NewGrpc()
}
