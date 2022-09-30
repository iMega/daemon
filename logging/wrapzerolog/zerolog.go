package wrapzerolog

import (
	"github.com/rs/zerolog"
)

type ZLog struct {
	wrapped zerolog.Logger
}

func New(log zerolog.Logger) *ZLog {
	return &ZLog{wrapped: log}
}

func (l *ZLog) Infof(format string, args ...interface{}) {
	l.wrapped.Info().Msgf(format, args...)
}

func (l *ZLog) Errorf(format string, args ...interface{}) {
	l.wrapped.Error().Msgf(format, args...)
}

func (l *ZLog) Debugf(format string, args ...interface{}) {
	l.wrapped.Debug().Msgf(format, args...)
}
