package slog

import (
	"bufio"
	"errors"
	"io"
	"sync/atomic"
	"time"
)

var loggers []*LogHandle
var curLoggersID uint32

var bufferSize = 1024
var loggersCapacity uint32

const maxLoggers = 128
const anyLogs string = "*"

// LogHandle is named logger
type LogHandle struct {
	ID          uint32
	writer      *bufio.Writer
	autoFlush   bool
	minLogLevel int
	lastWrite   time.Time
	name        string
}

func initLoggersArray(capacity int) error {
	if capacity > maxLoggers {
		return errors.New("Too many loggers")
	}

	loggers = make([]*LogHandle, capacity)

	loggersCapacity = uint32(capacity)
	curLoggersID = 0

	return nil
}

// создавать логгеры по rules. One rule - one logger and copy pointer if requested new loggers with some rule
func createInstanceLogger(iw io.Writer, minLogLevel int, autoFlush bool) *LogHandle {
	if loggers == nil {
		panic("Loggers not initiated")
	}

	if curLoggersID >= loggersCapacity {
		panic("Too many loggers")
	}

	id := atomic.AddUint32(&curLoggersID, 1) - 1

	w := bufio.NewWriterSize(iw, bufferSize)

	lh := &LogHandle{
		ID:          id,
		writer:      w,
		autoFlush:   autoFlush,
		minLogLevel: minLogLevel,
	}

	loggers[id] = lh

	return lh
}

func getLogger(idx int) *LogHandle {
	return loggers[idx]
}

// Log is writing a message to bufio.Writer
func Log(id uint32, msg []byte) {
	loggers[id].writer.Write(msg)
	loggers[id].lastWrite = time.Now()
}

// Flush bufio.Writer
func Flush(id uint32) {
	loggers[id].writer.Flush()
}

func getFilePath(name string) string {
	// get destination log file name
	r, exist := config.rules[name]

	if !exist {
		r = config.rules[anyLogs]
	}

	dest := r.destination

	// get path by log file name
	p, exist := config.paths[dest]

	if !exist {
		p = config.paths[anyLogs]
	}

	return p
}

func getLevel(name string) int {
	r, exist := config.rules[name]

	if !exist {
		r = config.rules[anyLogs]
	}

	return r.level
}
