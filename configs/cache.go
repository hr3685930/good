package configs

import "good/pkg/drive/config"

// Cache default once
type Cache struct {
	Sync    Sync
	Redis   Redis
	Default string `default:"sync" mapstructure:"default"`
}

// Redis Redis
type Redis struct {
	config.RedisDrive
	Dsn string `default:"" mapstructure:"dsn"`
}

// Sync Sync
type Sync struct {
	config.SyncDrive
}
