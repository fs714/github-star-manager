package config

import "github.com/spf13/viper"

var (
	Viper         *viper.Viper
	DefaultConfig Configuration = GetDefaultConfig()
	Config        Configuration = DefaultConfig
)

func GetDefaultConfig() Configuration {
	return Configuration{
		Common: Common{
			RunMode:   "release",
			Profiling: false,
		},
		Database: Database{
			Path: "./db.json",
		},
		Logging: Logging{
			File:       "",
			Level:      "info",
			Format:     "console",
			MaxSize:    10,
			MaxAge:     16,
			MaxBackups: 16,
			Compress:   true,
		},
		HttpServer: HttpServer{
			Host:         "0.0.0.0",
			Port:         "9500",
			ReadTimeout:  60,
			WriteTimeout: 60,
		},
	}
}

type Common struct {
	RunMode   string `mapstructure:"run_mode"`
	Profiling bool   `mapstructure:"profiling"`
}

type Database struct {
	Path string `mapstructure:"path"`
}

type Logging struct {
	File       string `mapstructure:"file"`
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
	Compress   bool   `mapstructure:"compress"`
}

type HttpServer struct {
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

type Configuration struct {
	Common     Common     `mapstructure:"common"`
	Database   Database   `mapstructure:"database"`
	Logging    Logging    `mapstructure:"logging"`
	HttpServer HttpServer `mapstructure:"http_server"`
}
