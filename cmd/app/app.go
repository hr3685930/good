package app

import (
	"fmt"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"good/configs"
	"good/internal/pkg/errs/export"
	"good/pkg/drive/cache"
	"good/pkg/drive/cache/redis"
	"good/pkg/drive/db"
	"good/pkg/drive/queue"
	"good/pkg/goo"
	"good/pkg/log"
	"good/pkg/tracing"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"
)

var (
	// Logger Logger
	Logger *zap.Logger
	// Queue Queue
	Queue queue.Queue
	// Orm Orm
	Orm *gorm.DB
	// Cache Cache
	Cache cache.Cache
	// Kafka Kafka
	Kafka sarama.Client
	// Redis Redis
	Redis *redis.Redis
)

// NewAPP NewAPP
func NewAPP() error {
	err := Log()
	if err != nil {
		return err
	}
	InitFacade()
	goo.New()
	goo.AsyncErrFunc = export.GoroutineErr
	err = NewTrace()
	if err != nil {
		return err
	}
	return Signal()
}


//NewTrace NewTrace
func NewTrace() error {
	tracing.TraceCloser = tracing.NewJaegerTracer(configs.ENV.App.Name, configs.ENV.Trace.Endpoint)
	Cache.AddTracingHook()
	return nil
}


// InitFacade InitFacade
func InitFacade() {
	Queue = queue.MQ
	Orm = db.Orm
	Cache = cache.Cached
	if queue.GetQueueDrive("kafka") != nil {
		KafkaDrive := queue.GetQueueDrive("kafka").(*queue.Kafka)
		Kafka, _ = KafkaDrive.GetCli()
	}

	if cache.GetCache("redis") != nil {
		Redis = cache.GetCache("redis").(*redis.Redis)
	}
}

// Signal Signal
func Signal() error {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM:
				tracing.TraceClose()
				fmt.Println("exit:", s)
				os.Exit(0)
			}
		}
	}()
	return nil
}

//Log log
func Log() error {
	path := "./storage/log/"
	filename := "app.log"
	l, err := log.NewLog(path, filename).InitZap()
	if err != nil {
		return err
	}
	Logger = l
	return nil
}
