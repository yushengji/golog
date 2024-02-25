package fslog

import (
	"go.uber.org/zap"
	"testing"
)

// 性能对比（在进行仅控制台打印前提下）：
// fslog:
// 32.50 ns/op           16 B/op          1 allocs/op
// zap:
// 5.073 ns/op            0 B/op          0 allocs/op

func BenchmarkFSLog(b *testing.B) {
	logger := New()
	logger.SetLv(ErrorLv)
	b.Run("fslog", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			logger.Debug("fslog")
		}
		b.StopTimer()
	})
}

func BenchmarkZap(b *testing.B) {
	production, err := zap.NewProduction()
	if err != nil {
		b.Fatal(err)
	}
	zapLogger := production.Sugar()
	b.Run("zap", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			zapLogger.Debug("zap")
		}
		b.StopTimer()
	})
}
