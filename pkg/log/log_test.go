package log

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var mockLogger *zaplogger

func TestMain(m *testing.M) {
	opts := &Options{
		Level:             "debug",
		DisableCaller:     false,
		DisableStacktrace: false,
		Format:            "json",
		OutputPaths:       []string{"stdout"},
	}

	Init(opts)
	mockLogger = std
	m.Run()
}

func TestLoggerMethods(t *testing.T) {
	assert.NotPanics(t, func() {
		Debugw("debug message", "key1", "value1")
		Infow("nfo message", "key2", "value2")
		Warnw("warn message", "key3", "value3")
		Errorw("error messgae", "key4", "value4")
	}, "Log methods should not cause a panic in this test")

	assert.Panics(t, func() {
		Panicw("fatal message", "key5", "value5")
	}, "Panicw should cause a panic and exit the program")
}

func TestLoggerInitialization(t *testing.T) {
	opts := NewOptions()
	logger := New(opts)

	assert.NotNilf(t, logger, "Logger should not be nil after initialization")
	assert.IsType(t, &zaplogger{}, logger, "Logger should be of type *zapLogger")
}

func TestSync(t *testing.T) {
	assert.NotPanics(t, func() {
		Sync()
	}, "Sync should not panic")
}
