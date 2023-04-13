package app

import (
	"github.com/urfave/cli"
	"good/cmd/app/job"
	commands "good/internal/logic/script"
)

// Commands cmd
var Commands = []cli.Command{
	{
		Name:   "event",
		Usage:  "kafka事件监听",
		Action: job.KafkaEventListen,
	},
	{
		Name:  "queue",
		Usage: "队列job",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "topic",
			},
		},
		Action: job.Queue,
	},
	{
		Name:  "db",
		Usage: "db操作",
		Subcommands: []cli.Command{
			{
				Name:   "migrate",
				Usage:  "迁移数据表",
				Action: commands.Migrate,
			},
		},
	},
	{
		Name:        "once",
		Usage:       "一次性脚本",
		Subcommands: []cli.Command{},
	},
}
