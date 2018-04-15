package redis

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type Options struct {
	Addr     string
	Password string
	Database int
	PoolSize int
}

type Redis = redis.Client

func New(opts Options) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     opts.Addr,
		Password: opts.Password,
		DB:       opts.Database,
		PoolSize: opts.PoolSize,
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, errors.Wrap(err, "error connecting to redis")
	}

	return client, nil
}
