package slog_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/bsbsm/slog"
)

//
// bench
//
func BenchmarkLogToFileSequential(b *testing.B) {
	slog.Init(1, 38)
	logger := slog.CreateLogger("logname", true)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Log("test message for logging in log")
	}
}

func BenchmarkCompareToStdlib(b *testing.B) {
	b.Run("slog", func(b *testing.B) {
		slog.Init(1, 38)
		logger := slog.CreateLogger("logname", true)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			logger.Log("test message for logging in log")
		}
	})

	b.Run("stdlog", func(b *testing.B) {
		var w, err = os.OpenFile("../bench_stndlog.log", os.O_APPEND|os.O_CREATE, 0644)

		if err != nil {
			panic("file")
		}

		level := 1

		l := log.New(w, "", 0)

		var m = "test message for logging in log"
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			// same in slog
			currentTime := time.Now().Local()
			dateStr := currentTime.Format("01.02.2006 15:04:05.000")

			var levelStr string
			switch {
			case level == 0:
				levelStr = "Verbose"
			case level == 1:
				levelStr = "Info"
			case level == 2:
				levelStr = "Warning"
			case level == 3:
				levelStr = "Error"
			default:
				levelStr = "Unknown"
			}

			arrayLen := 33 + len(m)
			var buf = make([]byte, arrayLen)
			_ = buf[32] // bound check trick

			copy(buf, dateStr)
			var whitespace byte = 32
			buf[23] = whitespace
			buf[27] = whitespace
			buf[28] = whitespace
			buf[29] = whitespace
			buf[30] = whitespace
			buf[31] = whitespace

			copy(buf[24:], levelStr)
			copy(buf[32:], m)

			buf[arrayLen-1] = 10 // "\n"

			l.Print(string(buf))
		}
	})
}
