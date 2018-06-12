package slog

import (
	"bufio"
	"io"
	"sync/atomic"
)

var loggers []LogHandle
var curLoggersID uint32

var maxLoggers uint32 = 128
var maxBufSize = 1024

// LogHandle is named logger
type LogHandle struct {
	ID          uint32
	format      string
	writer      *bufio.Writer
	autoFlush   bool
	minLogLevel int
}

func createInstanceLogger(iw io.Writer, frmt string, minLogLevel int, autoFlush bool) *LogHandle {
	if loggers == nil || len(loggers) == 0 {
		loggers = make([]LogHandle, maxLoggers)
		curLoggersID = 0
	}

	if curLoggersID >= maxLoggers {
		panic("Too many loggers")
	}

	id := atomic.AddUint32(&curLoggersID, 1) - 1

	w := bufio.NewWriterSize(iw, maxBufSize)

	lh := &LogHandle{
		ID:          id,
		format:      frmt,
		writer:      w,
		autoFlush:   autoFlush,
		minLogLevel: minLogLevel,
	}

	loggers[id] = *lh

	return lh
}

// Log is writing a message to bufio.Writer
func Log(id uint32, msg *string) {

	loggers[id].writer.Write([]byte(*msg))
}

// Flush bufio.Writer
func Flush(id uint32) {
	loggers[id].writer.Flush()
}
