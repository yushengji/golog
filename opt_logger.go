package fslog

type loggerConfig struct {
	SkipCaller int
}

type LoggerOption func(*loggerConfig)

func WithSkipCaller(skip int) LoggerOption {
	return func(config *loggerConfig) {
		config.SkipCaller = skip
	}
}
