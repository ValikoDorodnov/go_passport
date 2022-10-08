package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct {
	error *zerolog.Event
}

func NewLogger() *Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	return &Logger{
		error: log.Error(),
	}
}

func (l Logger) Error(msg string) {
	l.error.Msg(msg)
}
