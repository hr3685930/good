package goo

import (
	"fmt"
	"github.com/pkg/errors"
)

// AsyncErr AsyncErr
var AsyncErr chan error
// AsyncErrFunc AsyncErrFunc
var AsyncErrFunc func(err error)
// AsyncFunc AsyncFunc
type AsyncFunc func() error

// New New
func New() {
	AsyncErr = make(chan error, 1)
	AsyncErrFunc = errHandler
	go func() {
		for {
			select {
			case err := <-AsyncErr:
				AsyncErrFunc(err)
			}
		}
	}()
}

// GO 无需等待处理完成
func GO(fns AsyncFunc) {
	go func(f AsyncFunc) {
		defer func() {
			if err := recover(); err != nil {
				AsyncErr <- errors.Errorf("%+v\n", err)
			}
		}()
		err := f()
		if err != nil {
			AsyncErr <- err
		}
	}(fns)
}

func errHandler(err error)  {
	fmt.Printf(err.Error())
}