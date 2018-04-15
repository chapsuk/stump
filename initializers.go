package stump

import (
	"github.com/m1ome/stump/package/db"
	"github.com/m1ome/stump/package/redis"
	"github.com/m1ome/stump/package/logger"
	"github.com/m1ome/stump/package/raven"
	"github.com/m1ome/stump/package/config"
	"github.com/m1ome/stump/package/web"
)

func initConnections(s *Stump, opts *StorageOptions) error {
	// Connecting to databases
	if opts.Postgres {
		dbClient, err := db.New(db.Options{
			Addr:     s.Config().GetString("database.addr"),
			Username: s.Config().GetString("database.user"),
			Password: s.Config().GetString("database.password"),
			Database: s.Config().GetString("database.database"),
			PoolSize: s.Config().GetInt("database.pool_size"),
		})

		if err != nil {
			return err
		}

		s.db = dbClient
	}

	if opts.Redis {
		redisClient, err := redis.New(redis.Options{
			Addr:     s.Config().GetString("redis.addr"),
			Password: s.Config().GetString("redis.password"),
			Database: s.Config().GetInt("redis.database"),
			PoolSize: s.Config().GetInt("redis.pool_size"),
		})

		if err != nil {
			return err
		}
		s.redis = redisClient
	}

	return nil
}

func initLogger(s *Stump, opts Options) error {
	// Creating logger
	log, err := logger.New(&logger.Options{
		Level: opts.LogLevel,
	})

	if err != nil {
		return err
	}

	s.logger = log
	return nil
}

func initRaven(s *Stump) error {
	if dsn := s.Config().GetString("sentry.dsn"); dsn != "" {
		r, err := raven.New(&raven.Options{
			DSN: dsn,
		})

		if err != nil {
			return err
		}

		s.raven = r
	}

	return nil
}

func initConfig(s *Stump, opts Options) error {
	conf, err := config.New(&config.Options{
		Path:    opts.ConfigPath,
		Type:    opts.ConfigType,
		Name:    opts.ConfigName,
		AutoEnv: true,
	})

	if err != nil {
		return err
	}

	s.config = conf
	return nil
}

func initWeb(s *Stump, opts Options) error {
	e := web.New(&web.Options{
		Debug: s.Config().GetBool("web.debug"),
	})

	s.web = e
	return nil
}
