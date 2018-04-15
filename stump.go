package stump

import (
	"github.com/m1ome/stump/package/logger"
	"github.com/m1ome/stump/package/config"
	"github.com/m1ome/stump/package/raven"
	"github.com/m1ome/stump/package/db"
	"github.com/m1ome/stump/package/redis"
	"github.com/m1ome/stump/package/web"
	"github.com/m1ome/stump/package/cli"
)

type Stump struct {
	logger *logger.Logger
	config *config.Config
	raven  *raven.Raven
	db     *db.DB
	redis  *redis.Redis
	web    *web.Web
}

type Options struct {
	LogLevel   logger.LoggerLevel
	ConfigName string
	ConfigPath string
	ConfigType string
}

func (o *Options) init() {
	if o.LogLevel == 0 {
		o.LogLevel = logger.LoggerLevelDevelopment
	}

	if o.ConfigName == "" {
		o.ConfigName = "config"
	}

	if o.ConfigPath == "" {
		o.ConfigPath = "."
	}

	if o.ConfigType == "" {
		o.ConfigType = "yaml"
	}
}

func New(opts Options) (*Stump, error) {
	opts.init()
	stump := new(Stump)

	// Creating logger
	if err := initLogger(stump, opts); err != nil {
		return nil, err
	}

	// Init configuration
	if err := initConfig(stump, opts); err != nil {
		return nil, err
	}

	// Initializing Raven stuff
	if err := initRaven(stump); err != nil {
		return nil, err
	}

	// Initialize Web with all stuff
	if err := initWeb(stump, opts); err != nil {
		return nil, err
	}

	return stump, nil
}

type StorageOptions struct {
	Postgres bool
	Redis    bool
}

func (s *Stump) Start(name string, usage string) error {
	c := cli.New(&cli.Options{
		Name:  name,
		Usage: usage,
	})

	c.Add(s.serve())
	c.Add(s.migrate())

	return c.Run()
}

func (s *Stump) Storages(opts *StorageOptions) error {
	return initConnections(s, opts)
}

func (s *Stump) Web() *web.Web {
	return s.web
}

func (s *Stump) Config() *config.Config {
	return s.config
}

func (s *Stump) Logger() *logger.Logger {
	return s.logger
}

func (s *Stump) DB() *db.DB {
	return s.db
}

func (s *Stump) Redis() *redis.Redis {
	return s.redis
}

func (s *Stump) Raven() *raven.Raven {
	return s.raven
}
