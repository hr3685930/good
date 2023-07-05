package export

import (
	"fmt"
	"github.com/pkg/errors"
	"good/configs"
	"good/pkg/errs"
)

// GoroutineErr Custom Goroutine Err
func GoroutineErr(err error) {
	e, ok := err.(interface{ errs.Error })
	var stack string
	if ok {
		stack = e.GetStack()
	} else {
		stack = fmt.Sprintf("%+v\n", errors.New(err.Error()))
	}
	if configs.ENV.App.Debug {
		fmt.Println(stack)
		return
	}
	app := map[string]string{
		"name":        configs.ENV.App.Name,
		"environment": configs.ENV.App.Env,
	}
	option := map[string]interface{}{
		"error_type": "goroutine_error",
		"app":        app,
		"exception":  stack,
	}
	Report(option, "goroutine error")
}
