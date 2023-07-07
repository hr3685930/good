package configs

import (
	"good/pkg/drive/config"
)

// Database default
type Database struct {
	Sqlite  SQLite
	Mysql   Mysql
	Default string `default:"sqlite" mapstructure:"default"`
}

// Mysql Mysql
type Mysql struct {
	config.MYSQLDrive
	Dsn string `default:"" mapstructure:"dsn"`
}

// Postgre Postgre
type Postgre struct {
	config.PostgreDrive
	Dsn string `default:"" mapstructure:"dsn"`
}

// Clickhouse Clickhouse
type Clickhouse struct {
	config.ClickhouseDrive
	Dsn string `default:"" mapstructure:"dsn"`
}

// SQLite SQLite
type SQLite struct {
	config.SQLiteDrive
}
