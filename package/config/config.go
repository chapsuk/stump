package config

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Options struct {
	Path    string
	Type    string
	AutoEnv bool
}

var (
	DefaultConfig  *Config
	DefaultOptions = &Options{
		Path:    "./config.yml",
		Type:    "yaml",
		AutoEnv: true,
	}
)

type Config = viper.Viper

func New(opts *Options) (*Config, error) {
	if opts == nil {
		opts = DefaultOptions
	}

	v := viper.New()
	if len(opts.Type) > 0 {
		v.SetConfigType(opts.Type)
	}

	if opts.AutoEnv {
		v.AutomaticEnv()
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	}

	// Reading file
	if len(opts.Path) > 0 {
		v.SetConfigFile(opts.Path)

		if err := v.ReadInConfig(); err != nil {
			return nil, errors.Wrap(err, "error reading config file")
		}
	}

	return v, nil
}
