package raven

import (
	"github.com/getsentry/raven-go"
	"github.com/pkg/errors"
)

type Options struct {
	DSN string
}

type Raven = raven.Client

func New(opts *Options) (*Raven, error) {
	if opts == nil {
		return raven.DefaultClient, nil
	}

	c, err := raven.New(opts.DSN)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create Raven client")
	}

	return c, nil
}
