package engine

import (
	"context"
	"fmt"
	"io"
)

// Logger is a simplified interface for logging purposes.
type Logger interface {
	// Debugf logs with the DEBUG level.
	Debugf(format string, args ...any)

	// Errorf logs with the ERROR level.
	Errorf(format string, args ...any)

	// Infof logs with the INFO level.
	Infof(format string, args ...any)

	// Warnf logs with the WARN level.
	Warnf(format string, args ...any)
}

type noopLogger struct{}

var _noopLogger Logger = &noopLogger{} // ensure interface is implemented

// Debugf does nothing.
func (*noopLogger) Debugf(string, ...any) {}

// Errorf does nothing.
func (*noopLogger) Errorf(string, ...any) {}

// Infof does nothing.
func (*noopLogger) Infof(string, ...any) {}

// Warnf does nothing.
func (*noopLogger) Warnf(string, ...any) {}

type loggerKeyType string

// loggerKey is the context key for the logger.
const loggerKey loggerKeyType = "logger"

// GetLogger returns the context logger.
//
// By default it will a noop logger, but it can be set with WithLogger run option.
func GetLogger(ctx context.Context) Logger {
	log, ok := ctx.Value(loggerKey).(Logger)
	if !ok {
		return _noopLogger
	}
	return log
}

type testLogger struct{ writer io.Writer }

var _ Logger = (*testLogger)(nil) // ensure interface is implemented

// NewTestLogger creates a new logger with the input writer.
//
// This logger is expected to be used in tests.
// In no way it should be used in production since it's unoptimized.
func NewTestLogger(writer io.Writer) Logger {
	return &testLogger{writer: writer}
}

// Debugf implements Logger.
func (b *testLogger) Debugf(format string, args ...any) {
	b.print(format, args...)
}

// Errorf implements Logger.
func (b *testLogger) Errorf(format string, args ...any) {
	b.print(format, args...)
}

// Infof implements Logger.
func (b *testLogger) Infof(format string, args ...any) {
	b.print(format, args...)
}

// Warnf implements Logger.
func (b *testLogger) Warnf(format string, args ...any) {
	b.print(format, args...)
}

func (b *testLogger) print(format string, args ...any) {
	b.writer.Write([]byte(fmt.Sprintf(format, args...)))
}
