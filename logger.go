package fslog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"sync"
)

// Logger 是对 zap.SugaredLogger 的封装，简化其对外使用的接口，
// 并且针对于 Logger 所有配置都是支持运行时动态修改，且并发安全。
// Logger 也支持了针对于 zapcore.Core 的扩展，调用 Logger.AddCore 即可，
//
// 构建该对象请使用 New 进行创建，在该操作中会对部分必要属性进行初始化，
// 直接使用结构体创建会导致结构体不可用（甚至panic）。
//
// 如非深度定制化扩展，非必要不建议使用 Logger.AddCore 进行扩展，该操作会
// 导致客户端应用程序对zap包编译依赖，不保证fslog切换内部日志实现。
type Logger struct {
	// 日志打印用
	Logger zap.SugaredLogger

	// 互斥量
	// 用于内部不可并发逻辑使用
	lock sync.Mutex

	// 当前日志打印级别
	// 借助zap内部级别设置机制，该机制
	// 内部使用乐观锁，协程安全
	lv zap.AtomicLevel

	// 所有的core
	cores []zapcore.Core

	writers []io.Writer
	*loggerConfig
}

func New(opts ...LoggerOption) *Logger {
	logger := new(Logger)
	logger.loggerConfig = new(loggerConfig)
	if len(opts) == 0 {
		opts = []LoggerOption{WithSkipCaller(1)}
	}
	for _, opt := range opts {
		opt(logger.loggerConfig)
	}

	logger.lv = zap.NewAtomicLevelAt(zap.DebugLevel)
	logger.flushLogger()
	return logger
}

// With 拼接自定义信息
// 主要用于打印当前环境快照信息(变量或其他自定义信息)
// 打印后，该信息会跟随日志一起打印
func (l *Logger) With(k string, v any) *Logger {
	newL := l.clone()
	newL.Logger = *newL.Logger.With(zap.Any(k, v))
	return newL
}

// Debug 格式化打印调试级别日志
// 不同于zap内部可变参数逻辑，该可变参数是用于，字符串格式化的
func (l *Logger) Debug(msg string, vs ...any) {
	if len(vs) == 0 {
		l.Logger.Debug(msg)
		return
	}
	l.Logger.Debugf(msg, vs...)
}

// Info 格式化打印信息级别日志
// 不同于zap内部可变参数逻辑，该可变参数是用于，字符串格式化的
func (l *Logger) Info(msg string, vs ...any) {
	if len(vs) == 0 {
		l.Logger.Info(msg)
		return
	}
	l.Logger.Infof(msg, vs...)
}

// Warn 格式化打印警告级别日志
// 不同于zap内部可变参数逻辑，该可变参数是用于，字符串格式化的
func (l *Logger) Warn(msg string, vs ...any) {
	if len(vs) == 0 {
		l.Logger.Warn(msg)
		return
	}
	l.Logger.Warnf(msg, vs...)
}

// Error 打印错误级别日志
// 该方法具有两种传参形式：
//  1. error类型：会直接格式化打印%+v日志
//  2. 信息(格式化)
//  3. error+格式化信息：error会作为 With 格式存在，且依旧以%+v格式输出
func (l *Logger) Error(vs ...any) {
	if len(vs) == 0 {
		return
	}

	err, ok := vs[0].(error)
	if ok {
		if len(vs) == 1 {
			l.Logger.Errorf("%+v", err)
			return
		}

		withed := l.Logger.With(zap.String("err", fmt.Sprintf("%+v", err)))
		msg, ok := vs[1].(string)
		if ok {
			withed.Errorf(msg, vs[2:])
			return
		}

		withed.Error(vs[1:])
		return
	}

	msg, ok := vs[0].(string)
	if ok {
		if len(vs) > 1 {
			l.Logger.Errorf(msg, vs[1:]...)
			return
		}

		if len(vs) == 1 {
			l.Logger.Error(vs[0])
			return
		}
	}
}

// Lv 获取当前日志打印级别
func (l *Logger) Lv() Level {
	return Level(l.lv.Level())
}

// SetLv 设置当前日志打印级别
func (l *Logger) SetLv(lv Level) {
	l.lv.SetLevel(lv.zap())
}

// NewFileOutput 新增日志输出文件配置
func (l *Logger) NewFileOutput(opts ...FileOutputOpt) {
	cfg := new(outFileConfig)
	for _, opt := range opts {
		opt(cfg)
	}
	l.NewOutput(&lumberjack.Logger{
		Filename:   cfg.filename,
		MaxSize:    cfg.maxSize,
		MaxAge:     cfg.maxAge,
		MaxBackups: cfg.maxBackups,
		LocalTime:  cfg.localTime,
		Compress:   cfg.Compress,
	})
}

// NewOutput 新增日志输出位置
func (l *Logger) NewOutput(writer io.Writer) {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	l.writers = append(l.writers, writer)
	l.AddCore(zapcore.NewCore(zapcore.NewJSONEncoder(cfg), zapcore.AddSync(writer), l.lv))
}

// AddCore 添加Core
func (l *Logger) AddCore(core ...zapcore.Core) {
	l.addCoreOnly(core...)
	l.flushLogger()
}

// Flush 将缓冲区日志刷新至目标
func (l *Logger) Flush() {
	err := l.Logger.Sync()
	if err != nil {
		With("err", err.Error()).Warn("flushLogger log error")
	}
}

// 保存内部core(不负责刷新 logger)
func (l *Logger) addCoreOnly(core ...zapcore.Core) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.cores = append(l.cores, core...)
}

func (l *Logger) clone() *Logger {
	newL := &Logger{cores: l.cores, lv: l.lv, loggerConfig: l.loggerConfig}
	newL.flushLogger()
	return newL
}

// 刷新内部 logger
// 刷新互斥，触发刷新后，刷新完成前依旧按照旧的配置执行
func (l *Logger) flushLogger() {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.Logger = *zap.
		New(
			zapcore.NewTee(l.cores...),
			zap.AddCaller(),
		).
		Sugar().
		WithOptions(
			zap.AddCallerSkip(l.SkipCaller),
			zap.WithCaller(true),
		)
}
