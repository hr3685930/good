package goo

import (
	"context"
	"github.com/pkg/errors"
	"sync"
)

// SyncErrFunc SyncErrFunc
type SyncErrFunc func(ctx context.Context) error

// Group Group
type Group struct {
	wg      sync.WaitGroup
	errOnce sync.Once
	err     error
}

// One One
func (g *Group) One(ctx context.Context, fn SyncErrFunc) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if err := fn(ctx); err != nil {
			g.errOnce.Do(func() {
				g.err = err
			})
		}
	}()
}

// Wait Wait
func (g *Group) Wait() error {
	g.wg.Wait()
	return g.err
}

// SyncFunc SyncFunc
type SyncFunc func(ctx context.Context) (interface{}, error)

// All 有序返回结果 func协程一次处理  error nil也返回
func All(ctx context.Context, fns ...SyncFunc) ([]interface{}, []error) {
	rs := make([]interface{}, len(fns))
	errs := make([]error, len(fns))

	var wg sync.WaitGroup
	wg.Add(len(fns))

	for i, fn := range fns {
		go func(index int, f SyncFunc) {
			defer func() {
				if err := recover(); err != nil {
					rs[index] = nil
					errs[index] = errors.Errorf("%+v\n", err)
				}
				wg.Done()
			}()

			res, err := f(ctx)
			rs[index] = res
			errs[index] = err
		}(i, fn)
	}

	wg.Wait()

	return rs, errs
}
