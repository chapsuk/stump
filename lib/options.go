package lib

import "github.com/m1ome/stump/package/logger"

type Options struct {
	LoggerLevel logger.Level
	LoggerNop   bool
	ConfigPath  string
	ConfigType  string
	Name        string
}

const (
	DefaultConfigPath = "./config.yml"
	DefaultConfigType = "yaml"
)

var (
	DefaultOptions = &Options{
		LoggerLevel: logger.Development,
		LoggerNop:   false,
		ConfigType:  DefaultConfigType,
		ConfigPath:  DefaultConfigPath,
	}
)

func (o *Options) init() {
	if o.LoggerLevel == 0 {
		o.LoggerLevel = DefaultOptions.LoggerLevel
	}

	if o.ConfigPath == "" {
		o.ConfigPath = DefaultOptions.ConfigPath
	}

	if o.ConfigType == "" {
		o.ConfigType = DefaultOptions.ConfigType
	}
}
