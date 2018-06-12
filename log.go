package slog

import (
	"errors"
	"fmt"
	"os"
	"time"
)

// Log levels
const (
	VerboseLevel = iota
	InfoLevel
	WarningLevel
	ErrorLevel
)

// Log format: [time] [log_level] [message]
const defaultFormat string = "%s %s %s"

var format string

// LogLevel is a minimum level for logging message
var LogLevel = VerboseLevel
var forceFlush bool

// Init is performance a logger initialize
func Init(maxLoggersCount int, maxWriteBufSize int) *error {

	if maxLoggersCount < 1 {
		e := errors.New("Minimum loggers count is 1")
		return &e
	}

	if maxBufSize < 8 {
		e := errors.New("Minimum buffer len is 8")
		return &e
	}

	return nil
}

// CreateLogger is performance a logger initialize
func CreateLogger(logFilePath string, minLogLevel int, logFormat string, autoFlush bool) *LogHandle {
	var w, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Create log file ERROR " + fmt.Sprint(err))

		if os.IsNotExist(err) {
			err = nil
			w, err = os.Create(logFilePath)
		} else {
			return nil
		}
	}

	if minLogLevel > 0 {
		LogLevel = minLogLevel
	} else {
		LogLevel = VerboseLevel
	}

	if len(logFormat) > 0 {
		format = logFormat
	} else {
		format = defaultFormat
	}

	forceFlush = autoFlush

	return createInstanceLogger(w, format, minLogLevel, autoFlush)
}

func getLevelName(logLevel int) string {
	switch {
	case logLevel == 0:
		return "Verbose"
	case logLevel == 1:
		return "Info"
	case logLevel == 2:
		return "Warning"
	case logLevel == 3:
		return "Error"
	}

	return "Unknown"
}

func writeByLogLevel(id uint32, logLevel int, message string) {
	if logLevel >= LogLevel {
		currentTime := time.Now().Local()
		dateStr := currentTime.Format("01.02.2006 15:04:05")

		toLog := fmt.Sprintf(format, dateStr, getLevelName(logLevel), message) + "\n"

		Log(id, &toLog)

		if forceFlush {
			Flush(id)
		}
	}
}

// Verbose logging a message with VerboseLevel
func (h *LogHandle) Verbose(message string) {
	writeByLogLevel(h.ID, VerboseLevel, message)
}

// Info logging a message with InfoLevel
func (h *LogHandle) Info(message string) {
	writeByLogLevel(h.ID, InfoLevel, message)
}

// Warn logging a message with WarningLevel
func (h *LogHandle) Warn(message string) {
	writeByLogLevel(h.ID, WarningLevel, message)
}

// Error logging a message with ErrorLevel
func (h *LogHandle) Error(message string) {
	writeByLogLevel(h.ID, ErrorLevel, message)
}

// Log message with minimum log level
func (h *LogHandle) Log(message string) {
	writeByLogLevel(h.ID, LogLevel, message)
}
