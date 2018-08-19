package slog

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestGetLoggerToReturnNotNilLogger(t *testing.T) {
	// arrange

	// act
	logger, e := GetLogger("AnyLog")

	// assert
	if e != nil {
		t.Fatalf("Expected error == Nil but got '%s'", e.Error())
	}

	if logger == nil {
		t.Fatalf("Expected LogHandler != Nil but got Nil")
	}
}

func TestWriteByLogLevelToReturn(t *testing.T) {
	// arrange
	logger, _ := GetLogger("AnyLog")

	// act
	logger.Log("test")

	// assert
	if _, err := os.Stat("./log.log"); os.IsNotExist(err) {
		t.Fatalf("Expected 'LogFile is exist' but got '%s'", err.Error())
	}

	if logger.lastWrite.YearDay() != time.Now().YearDay() {
		t.Fatalf("Expected 'lastWrite is today' but got not today")
	}
}

func TestArchiveLogFileToReturn(t *testing.T) {
	// arrange
	logger, _ := GetLogger("AnyLog")
	logger.Log("test")
	currentDate := time.Now()
	logger.lastWrite = currentDate.Add(-24 * time.Hour)

	// act
	archiveLogFile(logger, logger.name, currentDate)

	// assert
	if _, err := os.Stat(fmt.Sprintf("./archive/log%slog",
		logger.lastWrite.Format("_02-01-2006."))); os.IsNotExist(err) {
		t.Fatalf("Expected ArchiveLogFile exist but got '%s'", err.Error())
	}
}
