package slog

import (
	"bufio"
	"errors"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

var loggers []*LogHandler
var curLoggersID uint32

var bufferSize int
var loggersCapacity uint32

const maxLoggersLimit = 128
const anyLogs string = "*"

// Log levels
const (
	VerboseLevel = iota
	InfoLevel
	WarningLevel
	ErrorLevel
)

// default values
const defaultBufferSize int = 512
const defaultLoggersCount int = 10

// LogHandler is named logger
type LogHandler struct {
	ID          uint32
	writeCloser *logFileWriteCloser
	autoFlush   bool
	minLogLevel int
	lastWrite   time.Time
	name        string
	mutex       sync.Mutex
}

type logFileWriteCloser struct {
	*bufio.Writer
	f *os.File
}

func initLoggersArray(capacity int) error {
	if capacity > maxLoggersLimit {
		return errors.New("Too many loggers")
	}

	loggers = make([]*LogHandler, capacity)

	loggersCapacity = uint32(capacity)
	curLoggersID = 0

	return nil
}

// создавать логгеры по rules. One rule - one logger and copy pointer if requested new loggers with some rule
func createInstanceLogger(name string, f *os.File, minLogLevel int, autoFlush bool) *LogHandler {
	if loggers == nil {
		panic("Loggers not initiated")
	}

	if curLoggersID >= loggersCapacity {
		panic("Too many loggers")
	}

	id := atomic.AddUint32(&curLoggersID, 1) - 1

	fw := &logFileWriteCloser{bufio.NewWriterSize(f, bufferSize), f}
	lh := &LogHandler{
		ID:          id,
		writeCloser: fw,
		autoFlush:   autoFlush,
		minLogLevel: minLogLevel,
		name:        name,
		mutex:       sync.Mutex{},
	}

	loggers[id] = lh

	return lh
}

func getLogger(idx int) *LogHandler {
	return loggers[idx]
}

// Log is writing a message to bufio.Writer
func Log(id uint32, msg []byte) {
	loggers[id].mutex.Lock()
	loggers[id].writeCloser.Write(msg)
	loggers[id].mutex.Unlock()

	loggers[id].lastWrite = time.Now()
}

// Flush bufio.Writer
func Flush(id uint32) {
	loggers[id].mutex.Lock()
	loggers[id].writeCloser.Flush()
	loggers[id].mutex.Unlock()
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
