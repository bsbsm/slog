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
		var w, err = os.OpenFile("../bench_log.stndlog", os.O_APPEND|os.O_CREATE, 0644)

		if err != nil {
			panic("file")
		}

		level := 1

		l := log.New(w, "", 0)

		var m = "test message for logging in log\n"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// same in slog
			currentTime := time.Now().Local()

			dateStr := currentTime.Format("01.02.2006 15:04:05")

			var loglevel string
			switch {
			case level == 0:
				loglevel = "Verbose"
			case level == 1:
				loglevel = "Info"
			case level == 2:
				loglevel = "Warning"
			case level == 3:
				loglevel = "Error"
			}

			var buf = make([]byte, 28+len(m))
			copy(buf, dateStr)
			copy(buf, " ")
			copy(buf, loglevel)
			copy(buf, " ")
			copy(buf, m)

			l.Print(string(buf))
		}
	})
}
