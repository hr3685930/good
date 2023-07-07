package app

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronjan/hunch"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	cmdHTTP "good/cmd/app/http"
	"good/cmd/app/job"
	"good/cmd/app/rpc"
	"good/configs"
	"good/internal/logic/http"
	"good/internal/pkg/errs/export"
	"good/pkg/drive"
	"good/pkg/drive/cache"
	"good/pkg/drive/config"
	"good/pkg/drive/db"
	"good/pkg/drive/queue"
	"good/pkg/goo"
	httpPkg "good/pkg/http"
	zap "good/pkg/log"
	"good/pkg/tracing"
	"gorm.io/gorm/logger"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"
)

//Config config
func Config() error {
	return config.Load(&configs.ENV)
}

//Log log
func Log() error {
	path := "./storage/log/"
	filename := "app.log"
	return zap.NewLog(path, filename).Init()
}

//Database db
func Database(ctx context.Context, ignoreErr bool) error {
	db.DefaultLogLevel = logger.Silent
	if configs.ENV.App.Debug {
		db.DefaultLogLevel = logger.Info
	}
	_, err := hunch.Retry(ctx, 0, func(c context.Context) (interface{}, error) {
		err := config.Drive(configs.ENV.Database, configs.ENV.App, ignoreErr)
		if err != nil {
			fmt.Println("数据库重连中...", err)
			time.Sleep(time.Second * 2)
		}
		return nil, err
	})
	if !ignoreErr {
		go func() {
			db.ConnStore.Range(func(key, value interface{}) bool {
				k := key.(string)
				d, _ := db.GetConnect(k).DB()
				go func() {
					for {
						err := d.PingContext(context.Background())
						if err != nil {
							fmt.Println(k + " connect error")
							os.Exit(0)
						}
						time.Sleep(time.Second * 5)
					}
				}()
				return true
			})
		}()
	}

	return err
}

//Cache cache
func Cache(ctx context.Context, ignoreErr bool) error {
	_, err := hunch.Retry(ctx, 0, func(c context.Context) (interface{}, error) {
		err := config.Drive(configs.ENV.Cache, configs.ENV.App, ignoreErr)
		if err != nil {
			fmt.Println("缓存重连中...", err)
			time.Sleep(time.Second * 2)
		}
		return nil, err
	})
	if !ignoreErr {
		go func() {
			cache.CacheMap.Range(func(key, value interface{}) bool {
				k := key.(string)
				d := cache.GetCache(k)
				go func() {
					for {
						if d.Ping() != nil {
							fmt.Println(k + " connect error")
							os.Exit(0)
						}
						time.Sleep(time.Second * 5)
					}
				}()
				return true
			})
		}()
	}

	return err
}

//Queue queue
func Queue(ctx context.Context, ignoreErr bool) error {
	_, errs := hunch.Retry(ctx, 0, func(c context.Context) (interface{}, error) {
		err := config.Drive(configs.ENV.Queue, configs.ENV.App, ignoreErr)
		if err != nil {
			fmt.Println("队列重连中...", err)
			time.Sleep(time.Second * 2)
		}
		return nil, err
	})

	if !ignoreErr {
		go func() {
			queue.QueueStore.Range(func(key, value interface{}) bool {
				k := key.(string)
				d := queue.GetQueueDrive(k)
				go func() {
					for {
						if d.Ping() != nil {
							fmt.Println(k + " connect error")
							os.Exit(0)
						}
						time.Sleep(time.Second * 5)
					}
				}()
				return true
			})
		}()
	}

	return errs
}

// APP APP
func APP() error {
	goo.New()
	goo.AsyncErrFunc = export.GoroutineErr

	drive.InitFacade()
	err := job.EventReceive()
	if err != nil {
		return err
	}

	err = queue.MQ.NewPublisher()
	if err != nil {
		return err
	}

	if configs.ENV.Queue.Default == "channel" {
		go job.SubscribeAll()
	}
	return tracing.NewTrace()
}

// Signal 信号量处理关闭
func Signal() error {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM:
				tracing.TraceClose()
				fmt.Println("退出:", s)
				os.Exit(0)
			}
		}
	}()
	return nil
}

// Command Command
func Command() error {
	flag.Parse()
	if len(flag.Args()) == 0 {
		return nil
	}
	app := cli.NewApp()
	app.Commands = Commands
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	_ = app.Run(os.Args)
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
