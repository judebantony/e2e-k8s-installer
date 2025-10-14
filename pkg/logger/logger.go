package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// Logger wraps zerolog with enhanced functionality for the installer
type Logger struct {
	logger zerolog.Logger
}

// LogLevel represents the logging level
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// LogFormat represents the logging format
type LogFormat string

const (
	LogFormatJSON LogFormat = "json"
	LogFormatText LogFormat = "text"
)

// Config holds logger configuration
type Config struct {
	Level  LogLevel
	Format LogFormat
	Output io.Writer
}

// NewLogger creates a new enhanced logger
func NewLogger(config Config) *Logger {
	// Set global log level
	level, err := zerolog.ParseLevel(string(config.Level))
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Configure output writer
	var writer io.Writer = os.Stdout
	if config.Output != nil {
		writer = config.Output
	}

	// Configure format
	var logger zerolog.Logger
	if config.Format == LogFormatText {
		// Human-readable console output
		consoleWriter := zerolog.ConsoleWriter{
			Out:        writer,
			TimeFormat: time.RFC3339,
			FormatLevel: func(i interface{}) string {
				level := strings.ToUpper(fmt.Sprintf("%s", i))
				switch level {
				case "DEBUG":
					return colorize("🔍 DEBUG", "37") // White
				case "INFO":
					return colorize("ℹ️  INFO ", "36") // Cyan
				case "WARN":
					return colorize("⚠️  WARN ", "33") // Yellow
				case "ERROR":
					return colorize("❌ ERROR", "31") // Red
				default:
					return level
				}
			},
			FormatMessage: func(i interface{}) string {
				return colorize(fmt.Sprintf("%s", i), "0") // Default color
			},
			FormatFieldName: func(i interface{}) string {
				return colorize(fmt.Sprintf("%s=", i), "36") // Cyan
			},
			FormatFieldValue: func(i interface{}) string {
				return fmt.Sprintf("%s", i)
			},
		}
		logger = zerolog.New(consoleWriter).With().Timestamp().Logger()
	} else {
		// JSON output
		logger = zerolog.New(writer).With().Timestamp().Logger()
	}

	return &Logger{
		logger: logger,
	}
}

// colorize adds ANSI color codes to text
func colorize(text, colorCode string) string {
	return fmt.Sprintf("\033[%sm%s\033[0m", colorCode, text)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string) *LogEvent {
	return &LogEvent{event: l.logger.Debug(), msg: msg}
}

// Info logs an info message
func (l *Logger) Info(msg string) *LogEvent {
	return &LogEvent{event: l.logger.Info(), msg: msg}
}

// Warn logs a warning message
func (l *Logger) Warn(msg string) *LogEvent {
	return &LogEvent{event: l.logger.Warn(), msg: msg}
}

// Error logs an error message
func (l *Logger) Error(msg string) *LogEvent {
	return &LogEvent{event: l.logger.Error(), msg: msg}
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string) *LogEvent {
	return &LogEvent{event: l.logger.Fatal(), msg: msg}
}

// LogEvent provides a fluent interface for adding fields to log entries
type LogEvent struct {
	event *zerolog.Event
	msg   string
}

// Str adds a string field
func (e *LogEvent) Str(key, value string) *LogEvent {
	e.event.Str(key, value)
	return e
}

// Int adds an integer field
func (e *LogEvent) Int(key string, value int) *LogEvent {
	e.event.Int(key, value)
	return e
}

// Int64 adds an int64 field
func (e *LogEvent) Int64(key string, value int64) *LogEvent {
	e.event.Int64(key, value)
	return e
}

// Float64 adds a float64 field
func (e *LogEvent) Float64(key string, value float64) *LogEvent {
	e.event.Float64(key, value)
	return e
}

// Bool adds a boolean field
func (e *LogEvent) Bool(key string, value bool) *LogEvent {
	e.event.Bool(key, value)
	return e
}

// Err adds an error field
func (e *LogEvent) Err(err error) *LogEvent {
	e.event.Err(err)
	return e
}

// Dur adds a duration field
func (e *LogEvent) Dur(key string, value time.Duration) *LogEvent {
	e.event.Dur(key, value)
	return e
}

// Time adds a time field
func (e *LogEvent) Time(key string, value time.Time) *LogEvent {
	e.event.Time(key, value)
	return e
}

// Interface adds an interface field
func (e *LogEvent) Interface(key string, value interface{}) *LogEvent {
	e.event.Interface(key, value)
	return e
}

// Step adds step-specific fields for installation tracking
func (e *LogEvent) Step(name string) *LogEvent {
	e.event.Str("step", name)
	return e
}

// Component adds component-specific fields
func (e *LogEvent) Component(name string) *LogEvent {
	e.event.Str("component", name)
	return e
}

// Progress adds progress information
func (e *LogEvent) Progress(current, total int) *LogEvent {
	e.event.Int("progress_current", current)
	e.event.Int("progress_total", total)
	e.event.Float64("progress_percent", float64(current)/float64(total)*100)
	return e
}

// Send completes the log entry
func (e *LogEvent) Send() {
	e.event.Msg(e.msg)
}

// Global logger instance
var globalLogger *Logger

// InitGlobalLogger initializes the global logger
func InitGlobalLogger(config Config) {
	globalLogger = NewLogger(config)
}

// GetLogger returns the global logger instance
func GetLogger() *Logger {
	if globalLogger == nil {
		globalLogger = NewLogger(Config{
			Level:  LogLevelInfo,
			Format: LogFormatText,
			Output: os.Stdout,
		})
	}
	return globalLogger
}

// Convenience functions for global logger
func Debug(msg string) *LogEvent {
	return GetLogger().Debug(msg)
}

func Info(msg string) *LogEvent {
	return GetLogger().Info(msg)
}

func Warn(msg string) *LogEvent {
	return GetLogger().Warn(msg)
}

func Error(msg string) *LogEvent {
	return GetLogger().Error(msg)
}

func Fatal(msg string) *LogEvent {
	return GetLogger().Fatal(msg)
}

// Step-specific logging helpers
func StepStart(step string) {
	Info("Step started").Step(step).Send()
}

func StepComplete(step string, duration time.Duration) {
	Info("Step completed").Step(step).Dur("duration", duration).Send()
}

func StepFailed(step string, err error) {
	Error("Step failed").Step(step).Err(err).Send()
}

func StepSkipped(step string, reason string) {
	Warn("Step skipped").Step(step).Str("reason", reason).Send()
}

// Progress logging
func ProgressUpdate(step string, current, total int, msg string) {
	Info(msg).Step(step).Progress(current, total).Send()
}

// Component logging
func ComponentStart(component, action string) {
	Info("Component action started").Component(component).Str("action", action).Send()
}

func ComponentComplete(component, action string, duration time.Duration) {
	Info("Component action completed").Component(component).Str("action", action).Dur("duration", duration).Send()
}

func ComponentFailed(component, action string, err error) {
	Error("Component action failed").Component(component).Str("action", action).Err(err).Send()
}