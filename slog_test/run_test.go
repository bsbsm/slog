package slog_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/bsbsm/slog"
)

// Run is test run
func Run() {
	logPath := "../log.log"
	slog.Init(3, 16)
	logger := slog.CreateLogger(logPath, slog.InfoLevel, "", true)
	logger.Log("test message")

	// f, e := os.Open(logPath)
	// if e != nil {
	// 	panic(e)
	// }
	// if err := reader.New(f, os.Stdout).Inflate(); err != nil {
	// 	panic(err)
	// }
}

func BenchmarkLogToFileSequential(b *testing.B) {
	slog.Init(3, 16)
	logger := slog.CreateLogger("../bench_log.log", slog.InfoLevel, "", true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Log("test message for logging in log file.")
	}
}

func BenchmarkCompareToStdlib(b *testing.B) {
	b.Run("slog", func(b *testing.B) {
		slog.Init(3, 16)
		logger := slog.CreateLogger("../bench_log.log", slog.InfoLevel, "", true)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			logger.Log("test message for logging in log file.")
		}
	})

	b.Run("stdlog", func(b *testing.B) {
		var w, err = os.OpenFile("../bench_log.stndlog", os.O_APPEND|os.O_CREATE, 0644)

		if err != nil {
			panic("file")
		}

		level := 1

		l := log.New(w, "", 0)

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

			l.Print(dateStr + " " + loglevel + " test message for logging in log file.\n")
		}
	})
}
