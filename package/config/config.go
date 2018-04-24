package config

import (
	"os"
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
	v.SetConfigType(opts.Type)

	if opts.AutoEnv {
		v.AutomaticEnv()
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	}

	// Reading file
	buf, err := os.Open(opts.Path)
	if err != nil {
		return nil, errors.Wrap(err, "error reading config file")
	}

	if err := v.ReadConfig(buf); err != nil {
		return nil, errors.Wrap(err, "error parsing config file")
	}

	return v, nil
}
