package fslog

import (
	"go.uber.org/zap/zapcore"
)

type Level zapcore.Level

const (
	DebugLv = Level(zapcore.DebugLevel)
	InfoLv  = Level(zapcore.InfoLevel)
	WarnLv  = Level(zapcore.WarnLevel)
	ErrorLv = Level(zapcore.ErrorLevel)
)

func (l Level) zap() zapcore.Level {
	return zapcore.Level(l)
}
