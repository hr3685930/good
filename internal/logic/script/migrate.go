package commands

import (
	"fmt"
	"github.com/urfave/cli"
	"good/cmd/app"
	"good/configs"
)

// Migrate migrate
func Migrate(c *cli.Context) {
	dbs := app.Orm.Debug()
	if configs.ENV.App.Env != "testing" {
		dbs = dbs.Set("gorm:table_options", "CHARSET=utf8mb4")
	}
	err := dbs.AutoMigrate()
	if err != nil {
		fmt.Print(err)
		return
	}
}
