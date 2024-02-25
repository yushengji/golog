package fslog

type outFileConfig struct {
	filename   string
	maxSize    int
	maxAge     int
	maxBackups int
	localTime  bool
	Compress   bool
}

// FileOutputOpt 文件输出选项
type FileOutputOpt func(*outFileConfig)

func WithFilename(filename string) FileOutputOpt {
	return func(c *outFileConfig) {
		c.filename = filename
	}
}

func WithMaxSize(maxSize int) FileOutputOpt {
	return func(c *outFileConfig) {
		c.maxSize = maxSize
	}
}

func WithMaxAge(maxAge int) FileOutputOpt {
	return func(c *outFileConfig) {
		c.maxAge = maxAge
	}
}

func WithMaxBackups(maxBackups int) FileOutputOpt {
	return func(c *outFileConfig) {
		c.maxBackups = maxBackups
	}
}

func WithLocalTime(localTime bool) FileOutputOpt {
	return func(c *outFileConfig) {
		c.localTime = localTime
	}
}

func WithCompress(compress bool) FileOutputOpt {
	return func(c *outFileConfig) {
		c.Compress = compress
	}
}
