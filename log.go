package slog

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// configuration
const configPath = "./slog.json"

var config *configuration

// Init is performance a logger initialize.
// If loggersCount == -1 then keep default value (10).
// If maxWriteBufSize == -1 then keep default value (512).
func Init(loggersCount int, maxWriteBufSize int) (e error) {
	config, e = readConfiguration(configPath)

	if e != nil {
		return e
	}

	loggersID = make(map[string]int)

	if loggersCount == -1 {
		loggersCount = defaultLoggersCount
	} else {
		if loggersCount < 1 {
			e = errors.New("Minimum loggers count is 1")
			return e
		}
	}

	if maxWriteBufSize == -1 {
		bufferSize = defaultBufferSize
	} else {
		if maxWriteBufSize < 8 {
			e = errors.New("Minimum buffer len is 8")
			return e
		}

		bufferSize = maxWriteBufSize
	}

	e = initLoggersArray(loggersCount)

	if e != nil {
		return e
	}

	return nil
}

// CreateLogger is performance a logger initialize
func CreateLogger(name string, autoFlush bool) *LogHandler {
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

	return createInstanceLogger(name, w, getLevel(name), autoFlush)
}

// map[LogName]LogID
var loggersID map[string]int

func GetLogger(name string) (*LogHandler, error) {
	var e error = nil
	var lg *LogHandler = nil
	exists := false

	if loggersID == nil {
		e = Init(defaultLoggersCount, defaultBufferSize)
	} else {
		idx, exists := loggersID[name]

		if exists {
			lg = loggers[idx]
		}
	}

	if !exists {
		lg = CreateLogger(name, true)
		loggersID[name] = int(lg.ID)
	}

	return lg, e
}

func writeByLogLevel(id uint32, logLevel int, message string, needFlush bool, logName string) {
	currentTime := time.Now().Local()

	archiveLogFile(loggers[id], logName, currentTime)

	dateStr := currentTime.Format("02.01.2006 15:04:05.000")

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

	arrayLen := 33 + len(message)
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
	copy(buf[32:], message)

	buf[arrayLen-1] = 10 // "\n"

	Log(id, buf)

	if needFlush {
		Flush(id)
	}
}

func archiveLogFile(log *LogHandler, logName string, currentTime time.Time) {
	lw := log.lastWrite

	// last write date check
	if lw.Year() == 1 || (lw.Year() == currentTime.Year() && lw.YearDay() >= currentTime.YearDay()) {
		return
	}

	filePath := getFilePath(logName)
	_, file := path.Split(filePath)

	// close log file
	err := log.writeCloser.f.Close()
	if err != nil {
		Log(log.ID, []byte(fmt.Sprintf("\t[ERROR] Cannot close file. %s\n", err.Error())))
		Flush(log.ID)
		return
	}

	var needReopen bool
	var errorText string

	// rename log file
	var absArchivePath string
	if filepath.IsAbs(config.archives[filePath]) {
		absArchivePath = config.archives[filePath]
	} else {
		absArchivePath, err = filepath.Abs(config.archives[filePath])
	}

	if err != nil {
		needReopen = true
		errorText = "Cannot resolve absolute path. " + err.Error()
	} else {
		if _, err := os.Stat(absArchivePath); os.IsNotExist(err) {
			os.Mkdir(absArchivePath, os.ModeDir)
		}

		// index of dot for split
		di := strings.Index(file, ".")

		err = os.Rename(filePath, path.Join(absArchivePath, file[:di]+lw.Format("_02-01-2006.")+file[di+1:]))
		if err != nil {
			needReopen = true
			errorText = "Cannot rename file. " + err.Error()
		}
	}

	if needReopen {
		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE, 0644)

		if err == nil && len(errorText) != 0 {
			Log(log.ID, []byte("\t[ERROR] "+errorText))
			Flush(log.ID)
			return
		} // else Log(log.ID, []byte("\t[ERROR] Cannot re-open file. "+err.Error()))

		log.writeCloser.f = f
		return
	}

	// create new log file
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE, 0644)
	if err == nil {
		// set logFIleWriteCloser
		log.writeCloser = &logFileWriteCloser{bufio.NewWriterSize(f, bufferSize), f}
	}
}

// Verbose logging a message with VerboseLevel
func (h *LogHandler) Verbose(message string) {
	writeByLogLevel(h.ID, VerboseLevel, message, h.autoFlush, h.name)
}

// Info logging a message with InfoLevel
func (h *LogHandler) Info(message string) {
	writeByLogLevel(h.ID, InfoLevel, message, h.autoFlush, h.name)
}

// Warn logging a message with WarningLevel
func (h *LogHandler) Warn(message string) {
	writeByLogLevel(h.ID, WarningLevel, message, h.autoFlush, h.name)
}

// Error logging a message with ErrorLevel
func (h *LogHandler) Error(message string) {
	writeByLogLevel(h.ID, ErrorLevel, message, h.autoFlush, h.name)
}

// Log message with minimum log level
func (h *LogHandler) Log(message string) {
	writeByLogLevel(h.ID, h.minLogLevel, message, h.autoFlush, h.name)
}
