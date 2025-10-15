package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nickygerritsen/switchboard/internal/config"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  *config.Config
		wantErr bool
	}{
		{
			name: "basic config without logging",
			config: &config.Config{
				DefaultBrowser: "chrome",
				Debug:          false,
				LogFile:        "",
			},
			wantErr: false,
		},
		{
			name: "config with debug enabled",
			config: &config.Config{
				DefaultBrowser: "chrome",
				Debug:          true,
				LogFile:        "auto",
			},
			wantErr: false,
		},
		{
			name: "config with custom log file",
			config: &config.Config{
				DefaultBrowser: "chrome",
				Debug:          false,
				LogFile:        filepath.Join(t.TempDir(), "test.log"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if logger != nil {
				defer logger.Close()

				if tt.config.Debug && logger.level != LevelDebug {
					t.Errorf("Expected debug level when debug is enabled, got %d", logger.level)
				}

				if !tt.config.Debug && logger.level != LevelInfo {
					t.Errorf("Expected info level when debug is disabled, got %d", logger.level)
				}
			}
		})
	}
}

func TestLogger_Logging(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	cfg := &config.Config{
		DefaultBrowser: "chrome",
		Debug:          true,
		LogFile:        logFile,
	}

	logger, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Write some log messages
	logger.Debug("Debug message: %s", "test")
	logger.Info("Info message: %s", "test")
	logger.Warn("Warn message: %s", "test")
	logger.Error("Error message: %s", "test")

	// Close to flush
	logger.Close()

	// Read log file
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	// Check that all messages are present
	if !strings.Contains(logContent, "[DEBUG]") {
		t.Error("Log file should contain [DEBUG]")
	}
	if !strings.Contains(logContent, "[INFO]") {
		t.Error("Log file should contain [INFO]")
	}
	if !strings.Contains(logContent, "[WARN]") {
		t.Error("Log file should contain [WARN]")
	}
	if !strings.Contains(logContent, "[ERROR]") {
		t.Error("Log file should contain [ERROR]")
	}
	if !strings.Contains(logContent, "Debug message: test") {
		t.Error("Log file should contain debug message")
	}
}

func TestLogger_LogLevels(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	cfg := &config.Config{
		DefaultBrowser: "chrome",
		Debug:          false, // Info level
		LogFile:        logFile,
	}

	logger, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Write messages at different levels
	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Error("Error message")

	logger.Close()

	// Read log file
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	// Debug should not be logged at info level
	if strings.Contains(logContent, "[DEBUG]") {
		t.Error("Log file should not contain [DEBUG] when debug is disabled")
	}
	// Info and Error should be logged
	if !strings.Contains(logContent, "[INFO]") {
		t.Error("Log file should contain [INFO]")
	}
	if !strings.Contains(logContent, "[ERROR]") {
		t.Error("Log file should contain [ERROR]")
	}
}

func TestInit(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	cfg := &config.Config{
		DefaultBrowser: "chrome",
		Debug:          true,
		LogFile:        logFile,
	}

	err := Init(cfg)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}
	defer Close()

	// Test package-level functions
	Debug("Package debug")
	Info("Package info")
	Warn("Package warn")
	Error("Package error")

	Close()

	// Read log file
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	if !strings.Contains(logContent, "Package debug") {
		t.Error("Log should contain package-level debug message")
	}
	if !strings.Contains(logContent, "Package info") {
		t.Error("Log should contain package-level info message")
	}
}

func TestLogger_NoFile(t *testing.T) {
	cfg := &config.Config{
		DefaultBrowser: "chrome",
		Debug:          false,
		LogFile:        "",
	}

	logger, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	if logger.logToFile {
		t.Error("Logger should not log to file when LogFile is empty")
	}

	// These should not error even without a file
	logger.Info("Test message")
	logger.Error("Test error")
}
