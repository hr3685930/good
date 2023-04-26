package redis

import (
	"context"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/pkg/errors"
	"net"
	"time"

	rd "github.com/go-redis/redis/v8"
)

// Redis Redis
type Redis struct {
	address string
	*rd.Client
	*redsync.Redsync
}

// New creates an instance of Redis cache driver
func New(dsn string) (*Redis, error) {
	opt, err := rd.ParseURL(dsn)
	if err != nil {
		return nil, errors.New("解析redis dsn失败")
	}
	conn := rd.NewClient(opt)
	if _, err := net.Dial("tcp", opt.Addr); err != nil {
		return nil, err
	}
	// new redis lock
	pool := goredis.NewPool(conn)
	rs := redsync.New(pool)
	return &Redis{address: opt.Addr, Client: conn, Redsync: rs}, nil
}

// Contains checks if cached key exists in Redis storage
func (r *Redis) Contains(ctx context.Context, key string) bool {
	status, _ := r.Exists(ctx, key).Result()
	if status > 0 {
		return true
	}
	return false
}

// Delete the cached key from Redis storage
func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.Del(ctx, key).Err()
}

// Fetch retrieves the cached value from key of the Redis storage
func (r *Redis) Fetch(ctx context.Context, key string) (string, error) {
	return r.Get(ctx, key).Result()
}

// FetchMulti retrieves multiple cached value from keys of the Redis storage
func (r *Redis) FetchMulti(ctx context.Context, keys []string) map[string]string {
	result := make(map[string]string)

	items, err := r.MGet(ctx, keys...).Result()
	if err != nil {
		return result
	}

	for i := 0; i < len(keys); i++ {
		if items[i] != nil {
			result[keys[i]] = items[i].(string)
		}
	}

	return result
}

// Flush removes all cached keys of the Redis storage
func (r *Redis) Flush(ctx context.Context) error {
	return r.FlushAll(ctx).Err()
}

// Save a value in Redis storage by key
func (r *Redis) Save(ctx context.Context, key string, value string, lifeTime time.Duration) error {
	return r.Set(ctx, key, value, lifeTime).Err()
}

// AddTracingHook AddTracingHook
func (r *Redis) AddTracingHook() {
	r.AddHook(NewTraceHook())
}

// AddMetricHook AddMetricHook
func (r *Redis) AddMetricHook() {
	r.AddHook(NewMetricHook(
		WithInstanceName("cache"),
		WithDurationBuckets([]float64{.001, .005, .01}),
	))
}

// Ping Ping
func (r *Redis) Ping() error {
	conn, err := net.Dial("tcp", r.address)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}
