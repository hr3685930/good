package commands

import (
	"fmt"
	"github.com/urfave/cli"
	"good/configs"
	"good/pkg/drive/db"
)

// Migrate migrate
func Migrate(c *cli.Context) {
	dbs := db.Orm.Debug()
	if configs.ENV.App.Env != "testing" {
		dbs = dbs.Set("gorm:table_options", "CHARSET=utf8mb4")
	}
	err := dbs.AutoMigrate()

	if err != nil {
		fmt.Print(err)
		return
	}
}
