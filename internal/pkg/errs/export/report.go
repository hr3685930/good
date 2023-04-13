package export

import (
	"github.com/ddliu/go-httpclient"
	"good/configs"
)

// Report Report
func Report(option map[string]interface{}, msg string) {
	go func() {
		if configs.ENV.App.ErrReport != "" {
			_, _ = httpclient.Begin().PostJson(configs.ENV.App.ErrReport, option)
		}
	}()
}
