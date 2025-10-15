package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/nickygerritsen/switchboard/internal/config"
)

// Logger levels
const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Logger represents a logger instance
type Logger struct {
	level     int
	debugLog  *log.Logger
	infoLog   *log.Logger
	warnLog   *log.Logger
	errorLog  *log.Logger
	logFile   *os.File
	logToFile bool
}

var defaultLogger *Logger

// Init initializes the default logger with the given configuration
func Init(cfg *config.Config) error {
	var err error
	defaultLogger, err = New(cfg)
	return err
}

// New creates a new logger instance
func New(cfg *config.Config) (*Logger, error) {
	logger := &Logger{
		level:     LevelInfo,
		logToFile: false,
	}

	if cfg.Debug {
		logger.level = LevelDebug
	}

	// Determine log file path
	var logPath string
	if cfg.LogFile != "" && cfg.LogFile != "auto" {
		logPath = cfg.LogFile
	} else if cfg.LogFile == "auto" || cfg.Debug {
		var err error
		logPath, err = config.GetLogPath()
		if err != nil {
			return nil, fmt.Errorf("failed to get log path: %w", err)
		}
	}

	// Open log file if path is set
	writer := io.Writer(io.Discard)
	if logPath != "" {
		// Ensure parent directory exists
		logDir := filepath.Dir(logPath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		logger.logFile = f
		logger.logToFile = true
		writer = f
	}

	// Create log instances
	logger.debugLog = log.New(writer, "[DEBUG] ", log.LstdFlags)
	logger.infoLog = log.New(writer, "[INFO]  ", log.LstdFlags)
	logger.warnLog = log.New(writer, "[WARN]  ", log.LstdFlags)
	logger.errorLog = log.New(writer, "[ERROR] ", log.LstdFlags)

	return logger, nil
}

// Close closes the logger and its associated file handles
func (l *Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

// Debug logs a debug message
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level <= LevelDebug {
		l.debugLog.Printf(format, v...)
	}
}

// Info logs an info message
func (l *Logger) Info(format string, v ...interface{}) {
	if l.level <= LevelInfo {
		l.infoLog.Printf(format, v...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level <= LevelWarn {
		l.warnLog.Printf(format, v...)
	}
}

// Error logs an error message
func (l *Logger) Error(format string, v ...interface{}) {
	if l.level <= LevelError {
		l.errorLog.Printf(format, v...)
	}
}

// Package-level logging functions using the default logger

// Debug logs a debug message using the default logger
func Debug(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Debug(format, v...)
	}
}

// Info logs an info message using the default logger
func Info(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Info(format, v...)
	}
}

// Warn logs a warning message using the default logger
func Warn(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Warn(format, v...)
	}
}

// Error logs an error message using the default logger
func Error(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Error(format, v...)
	}
}

// Close closes the default logger
func Close() error {
	if defaultLogger != nil {
		return defaultLogger.Close()
	}
	return nil
}
