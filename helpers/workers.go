package helpers

import (
	"time"
	"github.com/chapsuk/worker"
	"context"
	"github.com/m1ome/stump/package/logger"
	"github.com/go-redis/redis"
)

func Schedule(job func(ctx context.Context), every time.Duration) *worker.Worker {
	return worker.New(job).ByTicker(every)
}

type LockOptions struct {
	Key        string
	TTL        time.Duration
	Logger     *logger.Logger
	Retries    int
	RetryDelay time.Duration
}

func ScheduleWithLock(s *redis.Client, job func(ctx context.Context), every time.Duration, opts LockOptions) *worker.Worker {
	if opts.Logger == nil {
		opts.Logger = logger.Nop()
	}

	o := worker.RedisLockOptions{
		LockKey:  opts.Key,
		LockTTL:  opts.TTL,
		RedisCLI: s,
		Logger:   opts.Logger,
	}

	bsm := worker.BsmRedisLockOptions{
		RedisLockOptions: o,
		RetryCount:       opts.Retries,
		RetryDelay:       opts.RetryDelay,
	}

	return worker.New(job).WithBsmRedisLock(bsm).ByTicker(every)
}
