package logger

import "bytes"

// LogEntry is an interface to retrieve the bytes of anything.
//
// The common case as the package is testlogs would be to give a log entry.
// It should be used in ToString to use any logger.
type LogEntry interface {
	Bytes() ([]byte, error)
}

// ToString transforms a slice of log entry into a string concatenation.
func ToString[S ~[]E, E LogEntry](entries S) string {
	var agg bytes.Buffer
	for _, entry := range entries {
		b, _ := entry.Bytes()
		agg.Write(b)
		agg.WriteString("\n")
	}
	return agg.String()
}
