package worker

import (
	"time"
	"errors"

	"github.com/bsm/redis-lock"
	"github.com/robfig/cron"
	"strings"
)

type Worker struct {
	opts Options
	cron *cron.Cron
}

type Options struct {
	Redis          RedisClient
	RedisNamespace string
}

type Task struct {
	Name        string
	Scheduler   string
	Handler     Handler
	Exclusive   bool
	LockTimeout time.Duration
	RetryCount  int
	RetryDelay  time.Duration
}

type Context struct {
	LockError error
	Lock      *lock.Locker
}

type Handler func(ctx Context) error
type RedisClient = lock.RedisClient

var (
	ErrEmptyName       = errors.New("empty task name")
	ErrEmptySchedule   = errors.New("empty scheduler string")
	ErrLockNotObtained = errors.New("error obtaining lock")
)

func New(opts Options) *Worker {
	c := cron.New()

	return &Worker{
		opts: opts,
		cron: c,
	}
}

func (w *Worker) Schedule(t Task) error {
	// Check all stuff
	if t.Name == "" {
		return ErrEmptyName
	}

	if t.Scheduler == "" {
		return ErrEmptySchedule
	}

	if t.Exclusive {
		if t.LockTimeout == 0 {
			t.LockTimeout = time.Minute
		}

		if t.RetryDelay == 0 {
			t.RetryDelay = time.Millisecond * 10
		}
	}

	wrapped := w.wrap(t)
	return w.cron.AddFunc(
		t.Scheduler,
		wrapped,
	)
}

func (w *Worker) wrap(t Task) func() {
	return func() {
		ctx := Context{}
		name := strings.Join([]string{w.opts.RedisNamespace, t.Name}, ".")

		if t.Exclusive {
			l, err := lock.Obtain(w.opts.Redis, name, &lock.Options{
				LockTimeout: t.LockTimeout,
				RetryCount:  t.RetryCount,
				RetryDelay:  t.RetryDelay,
			})

			if err != nil {
				if err == lock.ErrLockNotObtained {
					err = ErrLockNotObtained
				}

				ctx.LockError = err
				t.Handler(ctx)
				return
			}

			ctx.Lock = l
			defer l.Unlock()
		}

		t.Handler(ctx)
	}
}

func (w *Worker) Start() {
	w.cron.Start()
}

func (w *Worker) Stop() {
	w.cron.Stop()
}
