package logger

import "go.uber.org/zap"

type LoggerLevel int

const (
	LoggerLevelDevelopment LoggerLevel = iota+1
	LoggerLevelProduction
)

type Logger = zap.SugaredLogger