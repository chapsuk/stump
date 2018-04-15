package logger

import (
	"go.uber.org/zap"
	"github.com/pkg/errors"
)

type Options struct {
	Level LoggerLevel
}

var (
	DefaultOptions = &Options{
		Level: LoggerLevelDevelopment,
	}

	ErrUnknownLevel = errors.New("unknown level")
)

func New(opts *Options) (*Logger, error) {
	var (
		logger *zap.Logger
		err error
	)

	if opts == nil {
		opts = DefaultOptions
	}

	switch opts.Level {
	case LoggerLevelDevelopment:
		logger, err = zap.NewDevelopment()
	case LoggerLevelProduction:
		logger, err = zap.NewProduction()
	default:
		return nil, ErrUnknownLevel
	}

	if err != nil {
		return nil, errors.Wrap(err, "cannot create logger")
	}

	sugar := logger.Sugar()
	return sugar, nil
}