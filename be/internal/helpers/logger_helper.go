package helpers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Logger holds logger configuration.
type Logger struct {
	logDir      string
	currentFile string
	file        *os.File
	logger      *log.Logger
	mu          sync.Mutex
}

// LogEntry represents a structured log entry.
type LogEntry struct {
	Function  string
	Message   string
	Level     string
	Timestamp time.Time
}

// NewLogger creates a new logger instance.
func NewLogger(logDir string) (*Logger, error) {
	// Create logs directory if not exists
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	l := &Logger{
		logDir: logDir,
	}

	// Initialize logger with today's file
	if err := l.rotateFile(); err != nil {
		return nil, fmt.Errorf("failed to rotate log file: %w", err)
	}

	// Start file rotation checker (runs at midnight)
	go l.startRotationChecker()

	// Start cleanup of old logs (30 days)
	go l.startCleanupChecker(30)

	return l, nil
}

// rotateFile opens a new log file for the current date.
func (l *Logger) rotateFile() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	today := time.Now().Format("2006-01-02")
	logPath := filepath.Join(l.logDir, fmt.Sprintf("%s.log", today))

	// Check if already using today's file
	if l.currentFile == logPath && l.file != nil {
		return nil
	}

	// Close previous file
	if l.file != nil {
		l.file.Close()
	}

	// Open new file (append mode)
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	l.currentFile = logPath
	l.file = file
	// Custom format: YYYY-MM-DD HH:MM:SS [no prefix - we add our own]
	l.logger = log.New(file, "", 0)

	return nil
}

// startRotationChecker checks daily at midnight if we need to rotate to a new file.
func (l *Logger) startRotationChecker() {
	for {
		now := time.Now()

		// Calculate time until next midnight
		midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		midnight = midnight.Add(24 * time.Hour) // Next midnight

		sleepDuration := time.Until(midnight)

		time.Sleep(sleepDuration)

		// Rotate at midnight
		if err := l.rotateFile(); err != nil {
			fmt.Printf("[Logger] Failed to rotate file: %v\n", err)
		}
	}
}

// startCleanupChecker periodically cleans up old log files at midnight.
func (l *Logger) startCleanupChecker(maxAgeDays int) {
	for {
		now := time.Now()

		// Calculate time until next midnight
		midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		midnight = midnight.Add(24 * time.Hour) // Next midnight

		sleepDuration := time.Until(midnight)

		time.Sleep(sleepDuration)

		// Cleanup at midnight
		if err := l.cleanupOldLogs(maxAgeDays); err != nil {
			fmt.Printf("[Logger] Failed to cleanup old logs: %v\n", err)
		}
	}
}

// cleanupOldLogs removes log files older than maxAge days.
func (l *Logger) cleanupOldLogs(maxAgeDays int) error {
	files, err := os.ReadDir(l.logDir)
	if err != nil {
		return err
	}

	cutoff := time.Now().AddDate(0, 0, -maxAgeDays)

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".log" {
			continue
		}

		// Parse date from filename (e.g., 2025-03-31.log)
		dateStr := strings.TrimSuffix(file.Name(), ".log")
		fileDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}

		if fileDate.Before(cutoff) {
			filePath := filepath.Join(l.logDir, file.Name())
			if err := os.Remove(filePath); err != nil {
				fmt.Printf("[Logger] Failed to delete old log file %s: %v\n", filePath, err)
			} else {
				fmt.Printf("[Logger] Deleted old log file: %s\n", filePath)
			}
		}
	}

	return nil
}

// formatLog formats a log message with timestamp and function context.
func (l *Logger) formatLog(function, message string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	if function != "" {
		return fmt.Sprintf("%s [%s] %s", timestamp, function, message)
	}
	return fmt.Sprintf("%s %s", timestamp, message)
}

// Printf prints formatted message to log with timestamp.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logger != nil {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		l.logger.Println(fmt.Sprintf("%s %s", timestamp, fmt.Sprintf(format, v...)))
	}
}

// Println prints message to log with timestamp.
func (l *Logger) Println(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logger != nil {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		l.logger.Println(fmt.Sprintf("%s %s", timestamp, fmt.Sprint(v...)))
	}
}

// Fatal prints message and exits with timestamp.
func (l *Logger) Fatal(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logger != nil {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		l.logger.Fatal(fmt.Sprintf("%s %s", timestamp, fmt.Sprint(v...)))
	}
}

// Fatalf prints formatted message and exits with timestamp.
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logger != nil {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		l.logger.Fatal(fmt.Sprintf("%s %s", timestamp, fmt.Sprintf(format, v...)))
	}
}

// Close closes the log file.
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// ==================== Structured Logging Helpers ====================

// LogStart logs the start of a function/operation.
// Usage: logger.LogStart("AuthLogin", "User login attempt: %s", email)
func (l *Logger) LogStart(function, format string, v ...interface{}) {
	l.Printf("[%s] ▶ [START] %s", function, fmt.Sprintf(format, v...))
}

// LogStep logs a step/operation inside a function (with indentation).
// Usage: logger.LogStep("AuthLogin", "User found: %s", email)
func (l *Logger) LogStep(function, format string, v ...interface{}) {
	l.Printf("[%s]   ├─ %s", function, fmt.Sprintf(format, v...))
}

// LogStepWithPrefix logs a step with custom prefix.
// Usage: logger.LogStepWithPrefix("AuthLogin", "[OK]", "Password validated")
func (l *Logger) LogStepWithPrefix(function, prefix, format string, v ...interface{}) {
	l.Printf("[%s]   %s %s", function, prefix, fmt.Sprintf(format, v...))
}

// LogEnd logs the successful end of a function/operation.
// Usage: logger.LogEnd("AuthLogin", "Login successful (duration: %v)", duration)
func (l *Logger) LogEnd(function, format string, v ...interface{}) {
	l.Printf("[%s] ✓ [END] %s", function, fmt.Sprintf(format, v...))
}

// LogError logs an error in a function/operation.
// Usage: logger.LogError("AuthLogin", "Failed to login: %v", err)
func (l *Logger) LogError(function, format string, v ...interface{}) {
	l.Printf("[%s] ✗ [ERROR] %s", function, fmt.Sprintf(format, v...))
}

// LogEndWithError logs the end of a function with error.
// Usage: logger.LogEndWithError("AuthLogin", "Login failed: %v", err)
func (l *Logger) LogEndWithError(function, format string, v ...interface{}) {
	l.Printf("[%s] ✗ [END] %s", function, fmt.Sprintf(format, v...))
}

// LogInfo logs an info message.
// Usage: logger.LogInfo("AuthLogin", "User logged in successfully")
func (l *Logger) LogInfo(function, format string, v ...interface{}) {
	l.Printf("[%s] [INFO] %s", function, fmt.Sprintf(format, v...))
}

// LogWarn logs a warning message.
// Usage: logger.LogWarn("AuthLogin", "Multiple failed login attempts")
func (l *Logger) LogWarn(function, format string, v ...interface{}) {
	l.Printf("[%s] [WARN] %s", function, fmt.Sprintf(format, v...))
}

// LogDebug logs a debug message.
// Usage: logger.LogDebug("AuthLogin", "Token payload: %+v", payload)
func (l *Logger) LogDebug(function, format string, v ...interface{}) {
	l.Printf("[%s] [DEBUG] %s", function, fmt.Sprintf(format, v...))
}
