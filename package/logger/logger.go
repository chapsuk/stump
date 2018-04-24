package logger

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Options struct {
	Level Level
	Nop   bool
}

var (
	DefaultOptions = &Options{
		Level: Development,
		Nop:   false,
	}

	ErrUnknownLevel = errors.New("unknown level")
)

func Nop() *Logger {
	return zap.NewNop().Sugar()
}

func New(opts *Options) (*Logger, error) {
	var (
		logger *zap.Logger
		err    error
	)

	if opts == nil {
		opts = DefaultOptions
	}

	if opts.Nop == true {
		// Create Nop logger
		return Nop(), nil
	}

	switch opts.Level {
	case Development:
		logger, err = zap.NewDevelopment()
	case Production:
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
