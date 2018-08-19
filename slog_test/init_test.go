package slog_test

import (
	"errors"
	"testing"

	"github.com/bsbsm/slog"
)

func TestInitToReturnNil(t *testing.T) {
	actualResult := slog.Init(1, 8)

	if actualResult != nil {
		t.Fatalf("Expected Nil but got %s", actualResult.Error())
	}
}

func TestInitToReturnIncorrectBufferSize(t *testing.T) {
	actualResult := slog.Init(1, 1)
	var expectedResult = errors.New("Minimum buffer len is 8")

	if actualResult.Error() != expectedResult.Error() {
		t.Fatalf("Expected %s but got %s", expectedResult, actualResult)
	}
}

func TestInitToReturnIncorrectLoggersCount(t *testing.T) {
	actualResult := slog.Init(0, 1)
	var expectedResult = errors.New("Minimum loggers count is 1")

	if actualResult.Error() != expectedResult.Error() {
		t.Fatalf("Expected %s but got %s", expectedResult, actualResult)
	}
}
