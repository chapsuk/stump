package logger

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Options struct {
	Level LoggerLevel
	Nop   bool
}

var (
	DefaultOptions = &Options{
		Level: LoggerLevelDevelopment,
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
