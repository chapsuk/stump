package db

import (
	"github.com/go-pg/pg"
	"github.com/pkg/errors"
)

type (
	DB = pg.DB
	PoolStats = pg.PoolStats
)

type Options struct {
	Addr     string
	Username string
	Password string
	Database string
	PoolSize int
}

const (
	DefaultPoolSize = 20
)

func (o *Options) init() {
	if o.PoolSize == 0 {
		o.PoolSize = DefaultPoolSize
	}
}

func New(opts Options) (*DB, error) {
	opts.init()

	db := pg.Connect(&pg.Options{
		Addr:     opts.Addr,
		User:     opts.Username,
		Password: opts.Password,
		Database: opts.Database,
		PoolSize: opts.PoolSize,
	})

	if _, err := db.Exec("SELECT 1"); err != nil {
		return nil, errors.Wrap(err, "error connecting to database")
	}

	return db, nil
}

func Stats(db *DB) *PoolStats {
	return db.PoolStats()
}
