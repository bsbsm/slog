package slog

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

// Log levels
const (
	VerboseLevel = iota
	InfoLevel
	WarningLevel
	ErrorLevel
)

// default values
const defaultBufferSize int = 512
const defaultMaxLoggerCount int = 10

// var format string
// var forceFlush bool

// configuration
const configPath = "./slog.json"

var config *configuration

// Init is performance a logger initialize
func Init(maxLoggersCount int, maxWriteBufSize int) (e error) {
	config, e = readConfiguration(configPath)

	if e != nil {
		return e
	}

	loggerNames = make(map[string]int)

	if maxLoggersCount < 1 {
		e = errors.New("Minimum loggers count is 1")
		return e
	}

	if maxWriteBufSize < 8 {
		e = errors.New("Minimum buffer len is 8")
		return e
	}

	e = initLoggersArray(maxLoggersCount)

	if e != nil {
		return e
	}

	bufferSize = maxWriteBufSize

	go ArchiveLogFilesLoop(config.paths, config.archives)

	return nil
}

// CreateLogger is performance a logger initialize
func CreateLogger(name string, autoFlush bool) *LogHandle {
	logFilePath := getFilePath(name)
	w, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Create log file ERROR " + fmt.Sprint(err))

		if os.IsNotExist(err) {
			err = nil
			w, err = os.Create(logFilePath)
		} else {
			return nil
		}
	}

	return createInstanceLogger(w, getLevel(name), autoFlush)
}

var loggerNames map[string]int

func GetLogger(name string) (*LogHandle, error) {
	var e error = nil
	var lg *LogHandle = nil
	exists := false

	if loggerNames == nil {
		e = Init(defaultMaxLoggerCount, defaultBufferSize)
	} else {
		idx, exists := loggerNames[name]

		if exists {
			lg = loggers[idx]
		}
	}

	if !exists {
		lg = CreateLogger(name, true)
		loggerNames[name] = int(lg.ID)
	}

	return lg, e
}

func writeByLogLevel(id uint32, logLevel int, message string, needFlush bool, logName string) {
	currentTime := time.Now().Local()
	lw := loggers[id].lastWrite
	// TODO проверка текущей даты
	if lw.YearDay() < currentTime.YearDay() {
		filePath := getFilePath(logName)
		_, file := path.Split(filePath)
		di := strings.Index(file, ".")
		err := os.Rename(filePath, path.Join(config.archives[filePath], file[:di]+lw.Format("02-01-2006")+file[di+1:]))
		if err != nil {
			Log(id, []byte("cannot rename file "+filePath))
		} else {
			f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE, 0644)
			if err == nil {
				loggers[id].writer = bufio.NewWriterSize(f, bufferSize)
			}
		}
	}

	dateStr := currentTime.Format("01.02.2006 15:04:05")

	var levelStr string
	switch {
	case logLevel == 0:
		levelStr = "Verbose"
	case logLevel == 1:
		levelStr = "Info"
	case logLevel == 2:
		levelStr = "Warning"
	case logLevel == 3:
		levelStr = "Error"
	default:
		levelStr = "Unknown"
	}

	var buf = make([]byte, 28+len(message))

	copy(buf, dateStr)
	copy(buf, " ")
	copy(buf, levelStr)
	copy(buf, " ")
	copy(buf, message)

	Log(id, buf)

	if needFlush {
		Flush(id)
	}
}

// Verbose logging a message with VerboseLevel
func (h *LogHandle) Verbose(message string) {
	writeByLogLevel(h.ID, VerboseLevel, message, h.autoFlush, h.name)
}

// Info logging a message with InfoLevel
func (h *LogHandle) Info(message string) {
	writeByLogLevel(h.ID, InfoLevel, message, h.autoFlush, h.name)
}

// Warn logging a message with WarningLevel
func (h *LogHandle) Warn(message string) {
	writeByLogLevel(h.ID, WarningLevel, message, h.autoFlush, h.name)
}

// Error logging a message with ErrorLevel
func (h *LogHandle) Error(message string) {
	writeByLogLevel(h.ID, ErrorLevel, message, h.autoFlush, h.name)
}

// Log message with minimum log level
func (h *LogHandle) Log(message string) {
	writeByLogLevel(h.ID, h.minLogLevel, message, h.autoFlush, h.name)
}
