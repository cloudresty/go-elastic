package elastic

// Logger defines the interface for pluggable logging within go-elastic.
// This allows users to integrate their preferred logging solution.
type Logger interface {
	// Info logs an informational message with optional structured fields
	Info(msg string, fields ...any)
	// Warn logs a warning message with optional structured fields
	Warn(msg string, fields ...any)
	// Error logs an error message with optional structured fields
	Error(msg string, fields ...any)
	// Debug logs a debug message with optional structured fields
	Debug(msg string, fields ...any)
}

// NopLogger is a no-operation logger that produces no output.
// This is used as the default logger when WithLogger is not provided.
type NopLogger struct{}

// Info implements Logger.Info with no operation
func (n *NopLogger) Info(msg string, fields ...any) {}

// Warn implements Logger.Warn with no operation
func (n *NopLogger) Warn(msg string, fields ...any) {}

// Error implements Logger.Error with no operation
func (n *NopLogger) Error(msg string, fields ...any) {}

// Debug implements Logger.Debug with no operation
func (n *NopLogger) Debug(msg string, fields ...any) {}
