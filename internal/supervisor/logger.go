package supervisor

import (
	"context"
	"fmt"
	"log"
	"os"
)

// Logger implements the types.Logger interface from opamp-go.
type Logger struct {
	logger *log.Logger
}

// NewLogger creates a new Logger.
func NewLogger() *Logger {
	return &Logger{
		logger: log.New(os.Stdout, "[supervisor] ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *Logger) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Printf("DEBUG: "+format, args...)
}

func (l *Logger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Printf("ERROR: "+format, args...)
}

// Ensure Logger implements the required interface.
var _ interface {
	Debugf(ctx context.Context, format string, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
} = (*Logger)(nil)

// getLogger returns a formatted string for logging.
func getLogger() {
	fmt.Println("Logger initialized")
}