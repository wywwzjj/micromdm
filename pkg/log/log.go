package log

import (
	"github.com/go-kit/kit/log/level"
)

type Logger interface {
	Log(keyvals ...interface{}) error
}

func Debug(logger Logger) Logger { return level.Debug(logger) }
func Info(logger Logger) Logger  { return level.Info(logger) }
