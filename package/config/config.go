package config

import (
	"github.com/spf13/viper"
	"github.com/pkg/errors"
)

type Options struct {
	Path    string
	Type    string
	Name    string
	AutoEnv bool
}

var (
	DefaultConfig *Config
	DefaultOptions = &Options{
		Path:    ".",
		Type:    "yaml",
		Name:    "config",
		AutoEnv: true,
	}
)

type Config = viper.Viper

func New(opts *Options) (*Config, error) {
	if opts == nil {
		opts = DefaultOptions
	}

	v := viper.New()

	v.SetConfigName(opts.Name)
	v.SetConfigType(opts.Type)
	v.AddConfigPath(opts.Path)

	if opts.AutoEnv {
		v.AutomaticEnv()
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "error reading config")
	}

	return v, nil
}

func SetDefault(conf *Config) {
	DefaultConfig = conf
}