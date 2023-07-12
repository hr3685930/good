package drive

import (
	"context"
	"fmt"
	"good/pkg/drive/cache"
	"good/pkg/drive/db"
	"good/pkg/drive/queue"
	"gorm.io/gorm"
	"time"
)

var (
	// Queue instance
	Queue queue.Queue
	// Orm Orm
	Orm *gorm.DB
	// Cache Cached
	Cache cache.Cache
)


// ListenDriveErr ListenDriveErr
func ListenDriveErr(fn func())  {
	go func() {
		db.ConnStore.Range(func(key, value interface{}) bool {
			k := key.(string)
			d, _ := db.GetConnect(k).DB()
			go func() {
				for {
					err := d.PingContext(context.Background())
					if err != nil {
						fmt.Println(k + " connect error")
						fn()
					}
					time.Sleep(time.Second * 5)
				}
			}()
			return true
		})

		cache.CacheMap.Range(func(key, value interface{}) bool {
			k := key.(string)
			d := cache.GetCache(k)
			go func() {
				for {
					if d.Ping() != nil {
						fmt.Println(k + " connect error")
						fn()
					}
					time.Sleep(time.Second * 5)
				}
			}()
			return true
		})

		queue.QueueStore.Range(func(key, value interface{}) bool {
			k := key.(string)
			d := queue.GetQueueDrive(k)
			go func() {
				for {
					if d.Ping() != nil {
						fmt.Println(k + " connect error")
						fn()
					}
					time.Sleep(time.Second * 5)
				}
			}()
			return true
		})
	}()
}
