package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var Log *Logger

// LogLevel represents the severity level of the log.
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
)

// Logger struct holds the configuration and loggers for both console and file.
type Logger struct {
	level      LogLevel
	consoleLog bool
	fileLog    bool
	file       *os.File
	mu         sync.Mutex
}

// NewLogger initializes the logger.
// `consoleLog` enables logging to console, `fileLog` enables logging to file, and `filePath` specifies the file to log to.
func NewLogger(consoleLog, fileLog bool, filePath string, level LogLevel) (*Logger, error) {
	logger := &Logger{
		level:      level,
		consoleLog: consoleLog,
		fileLog:    fileLog,
	}

	if fileLog {
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %v", err)
		}
		logger.file = file
	}

	Log = logger
	return logger, nil
}

// Close closes the log file if file logging is enabled.
func (l *Logger) Close() {
	if l.file != nil {
		l.file.Close()
	}
}

// log writes a message to the log based on the level.
func (l *Logger) log(level LogLevel, format string, v ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	message := fmt.Sprintf("[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, v...))
	levelStr := ""

	switch level {
	case DEBUG:
		levelStr = "[DEBUG]"
	case INFO:
		levelStr = "[INFO]"
	case WARNING:
		levelStr = "[WARNING]"
	case ERROR:
		levelStr = "[ERROR]"
	}

	// Log to console
	if l.consoleLog {
		log.Println(levelStr, message)
	}

	// Log to file
	if l.fileLog && l.file != nil {
		_, err := l.file.WriteString(levelStr + " " + message)
		if err != nil {
			fmt.Printf("Error writing to log file: %v\n", err)
		}
	}
}

// Debug logs a debug message.
func (l *Logger) Debug(format string, v ...interface{}) {
	l.log(DEBUG, format, v...)
}

// Info logs an info message.
func (l *Logger) Info(format string, v ...interface{}) {
	l.log(INFO, format, v...)
}

// Warning logs a warning message.
func (l *Logger) Warning(format string, v ...interface{}) {
	l.log(WARNING, format, v...)
}

// Error logs an error message.
func (l *Logger) Error(format string, v ...interface{}) {
	l.log(ERROR, format, v...)
}

// Timer is a struct for timing execution.
type Timer struct {
	start time.Time
}

// StartTimer starts a timer.
func StartTimer() *Timer {
	return &Timer{start: time.Now()}
}

// Stop logs the elapsed time and returns the duration.
func (t *Timer) Stop(l *Logger, msg string) time.Duration {
	duration := time.Since(t.start)
	l.Info("%s took %v", msg, duration)
	return duration
}
