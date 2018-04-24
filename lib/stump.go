package lib

import (
	"github.com/go-redis/redis"

	"github.com/m1ome/stump/package/logger"
	"github.com/m1ome/stump/package/db"
	"github.com/m1ome/stump/package/config"
	"github.com/m1ome/stump/package/web"
	"github.com/m1ome/stump/package/raven"
	"github.com/m1ome/stump/package/cli"
)

type Stump struct {
	logger       *logger.Logger
	db           *db.DB
	config       *config.Config
	web          *web.Web
	redis        *redis.Client
	raven        *raven.Raven
	cli          *cli.Cli

	opts *Options
}

func New(opts *Options) (*Stump, error) {
	if opts == nil {
		opts = DefaultOptions
	}
	opts.init()

	s := &Stump{
		opts: opts,
	}

	// Creating logger and reading configuration
	if err := initLogger(s); err != nil {
		return nil, err
	}

	// Loading configuraion
	if err := initConfig(s, opts.ConfigPath, opts.ConfigType); err != nil {
		return nil, err
	}

	// Loading CLI stuff
	if err := initCli(s, s.config); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Stump) Web() *web.Web {
	if s.web == nil {
		initWeb(s)
	}

	return s.web
}

func (s *Stump) ServeHTTP() error {
	address := s.config.GetString("web.address")
	if address == "" {
		s.logger.Info("Binding on default port: 8080")
		address = ":8080"
	} else {
		s.logger.Infof("Start listening on: %v", address)
	}

	return s.web.ListenAndServe(address)
}

func (s *Stump) Config() *config.Config {
	if s.config == nil {
		return &config.Config{}
	}

	return s.config
}

func (s *Stump) Logger() *logger.Logger {
	if s.logger == nil {
		return logger.Nop()
	}

	return s.logger
}

func (s *Stump) DB() *db.DB {
	if s.db == nil {
		return &db.DB{}
	}

	return s.db
}

func (s *Stump) Redis() *redis.Client {
	if s.redis == nil {
		return &redis.Client{}
	}

	return s.redis
}

func (s *Stump) Raven() *raven.Raven {
	if s.raven == nil {
		return &raven.Raven{}
	}

	return s.raven
}

func (s *Stump) Cli() *cli.Cli {
	return s.cli
}

func (s *Stump) ServeCommand(fn func () error) cli.Command {
	return cli.Command{
		Name:    "serve",
		Aliases: []string{"s"},
		Usage:   "start web server",
		Action: func(c *cli.Context) error {
			// Handling all stuff together
			if err := fn(); err != nil {
				return err
			}

			return s.ServeHTTP()
		},
	}
}

func (s *Stump) Init(redis, db bool) error {
	// Loading Raven error reporting
	if err := initRaven(s); err != nil {
		return err
	}

	// Loading redis
	if redis {
		if err := initRedis(s, s.config); err != nil {
			return err
		}
	}

	// Loading database
	if db {
		if err := initDatabase(s, s.config); err != nil {
			return err
		}
	}

	return nil
}

func (s *Stump) Run() error {
	s.cli.Add(cliMigrate(s))
	return s.cli.Run()
}
