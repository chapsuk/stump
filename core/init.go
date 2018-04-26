package core

import (
	"github.com/m1ome/stump/package/cli"
	"github.com/m1ome/stump/package/config"
	"github.com/m1ome/stump/package/db"
	"github.com/m1ome/stump/package/logger"
	"github.com/m1ome/stump/package/raven"
	"github.com/m1ome/stump/package/redis"
	"github.com/m1ome/stump/package/web"

	"os"
)

type Initer func(*Stump, *config.Config) error

func InitDatabase() Initer {
	return func(s *Stump, conf *config.Config) error {
		s.Logger().Infow(
			"Connecting to database",
			"addr", conf.GetString("database.addr"),
		)

		// Connecting to databases
		dbClient, err := db.New(db.Options{
			Addr:     conf.GetString("database.addr"),
			Username: conf.GetString("database.user"),
			Password: conf.GetString("database.password"),
			Database: conf.GetString("database.database"),
			PoolSize: conf.GetInt("database.pool_size"),
		})

		if err != nil {
			return err
		}

		s.SetDB(dbClient)
		return nil
	}

}

func InitCLI() Initer {
	return func(s *Stump, conf *config.Config) error {
		c := cli.New(&cli.Options{
			Name:    conf.GetString("cli.name"),
			Usage:   conf.GetString("cli.usage"),
			Version: conf.GetString("cli.version"),
		})

		s.SetCli(c)
		return nil
	}

}

func InitRedis() Initer {
	return func(s *Stump, conf *config.Config) error {
		s.Logger().Infow("Connecting to Redis", "addr", conf.GetString("redis.addr"))

		redisClient, err := redis.New(redis.Options{
			Addr:     conf.GetString("redis.addr"),
			Password: conf.GetString("redis.password"),
			Database: conf.GetInt("redis.database"),
			PoolSize: conf.GetInt("redis.pool_size"),
		})

		if err != nil {
			return err
		}
		s.redis = redisClient

		return nil
	}
}

func InitLogger() Initer {
	return func(s *Stump, conf *config.Config) error {
		// Creating logger
		log, err := logger.New(&logger.Options{
			Level: s.opts.LoggerLevel,
			Nop:   s.opts.LoggerNop,
		})

		if err != nil {
			return err
		}

		s.SetLogger(log)
		return nil
	}
}

func InitRaven() Initer {
	return func(s *Stump, conf *config.Config) error {
		if dsn := s.Config().GetString("sentry.dsn"); dsn != "" {
			s.Logger().Infow("Initializing Raven connection", "dsn", dsn)

			r, err := raven.New(&raven.Options{
				DSN: dsn,
			})

			if err != nil {
				return err
			}

			s.SetRaven(r)
		}

		return nil
	}
}

func InitConfig() Initer {
	return func(s *Stump, conf *config.Config) error {
		path := os.Getenv("CONFIG_PATH")
		t := os.Getenv("CONFIG_TYPE")

		if path == "" {
			path = DefaultConfigPath
		}

		if t == "" {
			t = DefaultConfigType
		}

		conf, err := config.New(&config.Options{
			Path:    path,
			Type:    t,
			AutoEnv: true,
		})

		if err != nil {
			return err
		}

		s.SetConfig(conf)
		return nil
	}
}

func initWeb(s *Stump) {
	s.web = web.New(&web.Options{
		Debug: s.Config().GetBool("web.debug"),
	})
}
