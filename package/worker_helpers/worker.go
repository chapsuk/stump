package worker_helpers

import (
	"time"

	"github.com/chapsuk/worker"
	"github.com/go-redis/redis"
	"github.com/m1ome/stump/package/logger"
)

type Options struct {
	Key        string
	TTL        time.Duration
	Retries    int
	RetryDelay time.Duration
	Logger     *logger.Logger
}

func (o *Options) init() {
	if o.Logger == nil {
		o.Logger = logger.Nop()
	}
}

func RedisLockOptions(cli *redis.Client, opts Options) worker.BsmRedisLockOptions {
	opts.init()

	o := worker.RedisLockOptions{
		LockKey:  opts.Key,
		LockTTL:  opts.TTL,
		RedisCLI: cli,
		Logger:   opts.Logger,
	}

	return worker.BsmRedisLockOptions{
		RedisLockOptions: o,
		RetryCount:       opts.Retries,
		RetryDelay:       opts.RetryDelay,
	}
}
