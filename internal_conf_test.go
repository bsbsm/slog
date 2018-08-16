package slog

import (
	"strconv"
	"testing"
)

func TestReadConfigurationToReturnParsedConf(t *testing.T) {
	// arrange

	// act
	actualResult, e := readConfiguration("./slog_test/slog.json")

	// assert
	if e != nil {
		t.Fatalf("Expected nil but got '%s'", e.Error())
	}

	if actualResult == nil {
		t.Fatalf("Expected not nil but got nil")
	}

	if actualResult.archives == nil ||
		actualResult.paths == nil ||
		actualResult.rules == nil ||
		len(actualResult.archives) == 0 ||
		len(actualResult.paths) == 0 ||
		len(actualResult.rules) == 0 {
		t.Fatalf("Expected not empty but got empty")
	}

	if len(actualResult.archives) != 1 ||
		len(actualResult.paths) != 1 ||
		len(actualResult.rules) != 1 {
		t.Fatalf("Expected 1 records in each map but got %s, %s, %s", strconv.Itoa(len(actualResult.archives)), strconv.Itoa(len(actualResult.paths)), strconv.Itoa(len(actualResult.rules)))
	}
}
