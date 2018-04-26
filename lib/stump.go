package lib

import (
	"github.com/go-redis/redis"

	"github.com/m1ome/stump/package/cli"
	"github.com/m1ome/stump/package/config"
	"github.com/m1ome/stump/package/db"
	"github.com/m1ome/stump/package/logger"
	"github.com/m1ome/stump/package/raven"
	"github.com/m1ome/stump/package/web"
)

type Stump struct {
	logger  *logger.Logger
	db      *db.DB
	config  *config.Config
	web     *web.Web
	redis   *redis.Client
	raven   *raven.Raven
	cli     *cli.Cli
	initers []Initer

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
	defaultIniters := []Initer{InitLogger(), InitConfig(), InitCLI()}
	if err := s.init(defaultIniters); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Stump) SetWeb(w *web.Web) {
	s.web = w
}

func (s *Stump) Web() *web.Web {
	if s.web == nil {
		initWeb(s)
	}

	return s.web
}

func (s *Stump) SetConfig(c *config.Config) {
	s.config = c
}

func (s *Stump) Config() *config.Config {
	if s.config == nil {
		return &config.Config{}
	}

	return s.config
}

func (s *Stump) SetLogger(logger *logger.Logger) {
	s.logger = logger
}

func (s *Stump) Logger() *logger.Logger {
	if s.logger == nil {
		return logger.Nop()
	}

	return s.logger
}

func (s *Stump) SetDB(db *db.DB) {
	s.db = db
}

func (s *Stump) DB() *db.DB {
	if s.db == nil {
		return &db.DB{}
	}

	return s.db
}

func (s *Stump) SetRedis(r *redis.Client) {
	s.redis = r
}

func (s *Stump) Redis() *redis.Client {
	if s.redis == nil {
		return &redis.Client{}
	}

	return s.redis
}

func (s *Stump) SetRaven(r *raven.Raven) {
	s.raven = r
}

func (s *Stump) Raven() *raven.Raven {
	if s.raven == nil {
		return &raven.Raven{}
	}

	return s.raven
}

func (s *Stump) SetCli(c *cli.Cli) {
	s.cli = c
}

func (s *Stump) Cli() *cli.Cli {
	return s.cli
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

func (s *Stump) SetIniters(initers ...Initer) {
	s.initers = initers
}

func (s *Stump) init(initers []Initer) error {
	for _, init := range initers {
		if err := init(s, s.config); err != nil {
			return err
		}
	}

	return nil
}

func (s *Stump) Init() error {
	return s.init(s.initers)
}

func (s *Stump) Run() error {
	s.cli.Add(cliMigrate(s))
	return s.cli.Run()
}
