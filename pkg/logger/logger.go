package logger

import (
	"log"
)

// Logger is a simplified interface for logging purposes in craft features.
//
// It's a interface because it would allow anyone to use any logger implementation (logrus, log, slog, etc.).
//
// In case you don't need or want a specific implementation, you can use default implementations Default (returns a std log implementation).
type Logger interface {
	// Info should log with the INFO level.
	Info(...any)

	// Infof should log with the INFO level and use format subtitution to take care of input args.
	Infof(string, ...any)

	// Warn should log with the WARN level.
	Warn(...any)

	// Warnf should log with the WARN level and use format subtitution to take care of input args.
	Warnf(string, ...any)
}

// std is a simple implementation of Logger for log std library.
type std struct {
	*log.Logger
}

var _ Logger = &std{} // ensure interface is implemented

// Default returns the default std logger (log library).
//
// It uses the global private logger (not a new instantiation).
func Default() Logger {
	return &std{log.Default()}
}

// Info logs with std logger using Print function.
//
// No logging level is involved since base std library doesn't handle logging level.
func (s *std) Info(args ...any) {
	s.Print(args...)
}

// Infof logs with std logger using Printf function.
//
// No logging level is involved since base std library doesn't handle logging level.
func (s *std) Infof(msg string, args ...any) {
	s.Printf(msg, args...)
}

// Warn logs with std logger using Print function.
//
// No logging level is involved since base std library doesn't handle logging level.
func (s *std) Warn(args ...any) {
	s.Print(args...)
}

// Warnf logs with std logger using Printf function.
//
// No logging level is involved since base std library doesn't handle logging level.
func (s *std) Warnf(msg string, args ...any) {
	s.Printf(msg, args...)
}
