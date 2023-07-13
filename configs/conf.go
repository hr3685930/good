package configs

import (
	"good/pkg/drive/cache/redis"
	"good/pkg/drive/cache/sync"
	"good/pkg/drive/db"
	"good/pkg/drive/queue"
)

// ENV ENV
var ENV Conf

//Conf Conf
type Conf struct {
	App      App
	Database Database
	Queue    Queue
	Cache    Cache
	Trace    Trace
}

// App App
type App struct {
	Name      string `default:"demo" mapstructure:"name"`  //应用名
	Env       string `default:"local" mapstructure:"env"`  //环境
	Debug     bool   `default:"true" mapstructure:"debug"` //开启debug
	ErrReport string `default:"" mapstructure:"err_report"`
}

// Cache default once
type Cache struct {
	Sync    sync.Drive
	Redis   redis.Drive
	Default string `default:"sync" mapstructure:"default"`
}

// Database default
type Database struct {
	Sqlite  db.SQLiteDrive
	Mysql   db.MYSQLDrive
	Default string `default:"sqlite" mapstructure:"default"`
}

// Queue default
type Queue struct {
	Default string `default:"channel" mapstructure:"default"`
	Kafka   queue.KafkaDrive
	Channel queue.ChannelDrive
}

//Trace Trace
type Trace struct {
	Endpoint string `default:"" mapstructure:"endpoint"`
}
