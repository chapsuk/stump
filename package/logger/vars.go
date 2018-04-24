package logger

import "go.uber.org/zap"

type Level int

const (
	Development Level = iota + 1
	Production
)

type Logger = zap.SugaredLogger
