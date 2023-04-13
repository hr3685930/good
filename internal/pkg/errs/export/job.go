package export

import (
	"fmt"
	"github.com/pkg/errors"
	"good/configs"
	"good/pkg/goo"
)

// JobErrorReport JobErrorReport
func JobErrorReport(err error, errorType string, data interface{}) {
	e, ok := err.(interface{ goo.Error })
	var stack string
	if ok {
		stack = e.GetStack()
	}else{
		stack = fmt.Sprintf("%+v\n", errors.New(err.Error()))
	}
	if configs.ENV.App.Env == "local" {
		fmt.Println(stack)
		return
	}
	app := map[string]string{
		"name":        configs.ENV.App.Name,
		"environment": configs.ENV.App.Env,
	}
	option := map[string]interface{}{
		"error_type": errorType,
		"app":        app,
		"exception":  stack,
		"data":       fmt.Sprintf("%s", data),
	}
	Report(option, "job error")
}
