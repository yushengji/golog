package fslog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

var Log *Logger

func init() {
	consoleCfg := zap.NewDevelopmentEncoderConfig()
	consoleCfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	jsonCfg := zap.NewProductionEncoderConfig()
	jsonCfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")

	Log = &Logger{
		lv:           zap.NewAtomicLevelAt(zap.DebugLevel),
		loggerConfig: &loggerConfig{SkipCaller: 2},
	}
	Log.cores = append(Log.cores,
		zapcore.NewCore(zapcore.NewConsoleEncoder(consoleCfg), zapcore.AddSync(os.Stdout), Log.lv),
	)
	Log.flushLogger()
}

// With 见 Logger.With
func With(k string, v any) *Logger {
	return Log.With(k, v)
}

// Info 见 Logger.Info
func Info(msg string, vs ...any) *Logger {
	Log.Info(msg, vs...)
	return Log
}

// Debug 见 Logger.Debug
func Debug(msg string, vs ...any) *Logger {
	Log.Debug(msg, vs...)
	return Log
}

// Warn 见 Logger.Warn
func Warn(msg string, vs ...any) *Logger {
	Log.Warn(msg, vs...)
	return Log
}

// Error 见 Logger.Error
func Error(vs ...any) *Logger {
	Log.Error(vs...)
	return Log
}

// Lv 见 Logger.Lv
func Lv() Level {
	return Log.Lv()
}

// SetLv 见 Logger.SetLv
func SetLv(lv Level) *Logger {
	Log.SetLv(lv)
	return Log
}

// NewFileOutput 见 Logger.NewFileOutput
func NewFileOutput(opts ...FileOutputOpt) *Logger {
	Log.NewFileOutput(opts...)
	return Log
}

// NewOutput 见 Logger.NewOutput
func NewOutput(writer io.Writer) *Logger {
	Log.NewOutput(writer)
	return Log
}

// SetLogger 替换fslog包默认的log对象
// Logger 支持了zapcore级别的底层配置，见 Logger
func SetLogger(logger *Logger) *Logger {
	Log = logger
	return Log
}

// Flush 见 Logger.Flush
func Flush() {
	Log.Flush()
}
