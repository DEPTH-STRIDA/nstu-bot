package log

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// LogLevel определяет уровень логирования
type LogLevel string

const (
	INFO  LogLevel = "INFO"
	ERROR LogLevel = "ERROR"
	WARN  LogLevel = "WARN"
	DEBUG LogLevel = "DEBUG"
)

// LogEntry структура для записи логов в формате JSON
type LogEntry struct {
	Timestamp string   `json:"time"`
	Level     LogLevel `json:"level"`
	Message   string   `json:"message"`
}

// Logger интерфейс для логгеров
type Logger interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
}

// ConsoleLogger логгер для логирования в консоль
type ConsoleLogger struct {
	mu sync.Mutex
}

// NewConsoleLogger конструктор консольного логгера
func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

func (c *ConsoleLogger) log(level LogLevel, args ...interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry := LogEntry{
		Timestamp: time.Now().Format("02-01-2006 15:04:05"),
		Level:     level,
		Message:   fmt.Sprint(args...),
	}
	data, _ := json.Marshal(entry)
	fmt.Println(string(data))
}

func (c *ConsoleLogger) logf(level LogLevel, format string, args ...interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry := LogEntry{
		Timestamp: time.Now().Format("02-01-2006 15:04:05"),
		Level:     level,
		Message:   fmt.Sprintf(format, args...),
	}
	data, _ := json.Marshal(entry)
	fmt.Println(string(data))
}

func (c *ConsoleLogger) Warn(args ...interface{}) {
	c.log(WARN, args...)
}

func (c *ConsoleLogger) Warnf(format string, args ...interface{}) {
	c.logf(WARN, format, args...)
}

func (c *ConsoleLogger) Error(args ...interface{}) {
	c.log(ERROR, args...)
}

func (c *ConsoleLogger) Errorf(format string, args ...interface{}) {
	c.logf(ERROR, format, args...)
}

func (c *ConsoleLogger) Info(args ...interface{}) {
	c.log(INFO, args...)
}

func (c *ConsoleLogger) Infof(format string, args ...interface{}) {
	c.logf(INFO, format, args...)
}

func (c *ConsoleLogger) Debug(args ...interface{}) {
	c.log(DEBUG, args...)
}

func (c *ConsoleLogger) Debugf(format string, args ...interface{}) {
	c.logf(DEBUG, format, args...)
}

// FileLogger логгер для логирования в один файл
type FileLogger struct {
	mu        sync.Mutex
	file      *os.File
	path      string
	startTime time.Time
}

// NewFileLogger конструктор файлового логгера
func NewFileLogger(basePath string) (*FileLogger, error) {
	logger := &FileLogger{
		path:      basePath,
		startTime: time.Now(),
	}

	err := logger.initialize()
	if err != nil {
		return nil, err
	}

	go logger.rotateLogs()

	return logger, nil
}

// initialize создает файл для логирования
func (f *FileLogger) initialize() error {
	dir := filepath.Dir(f.path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	filename := filepath.Join(dir, "log_"+time.Now().Format("02-01-2006_15-04-05")+".log")
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	f.file = file
	return nil
}

// rotateLogs отвечает за периодическое обновление лог-файла
func (f *FileLogger) rotateLogs() {
	for {
		time.Sleep(time.Until(f.startTime.AddDate(0, 1, 0)))
		f.mu.Lock()
		f.file.Close()
		f.file = nil
		f.startTime = time.Now()
		err := f.initialize()
		if err != nil {
			fmt.Printf("Failed to rotate logs: %v\n", err)
		}
		f.mu.Unlock()
	}
}

func (f *FileLogger) log(level LogLevel, args ...interface{}) {
	f.mu.Lock()
	defer f.mu.Unlock()

	entry := LogEntry{
		Timestamp: time.Now().Format("02-01-2006 15:04:05"),
		Level:     level,
		Message:   fmt.Sprint(args...),
	}

	data, _ := json.Marshal(entry)
	fmt.Fprintln(f.file, string(data))
}

func (f *FileLogger) logf(level LogLevel, format string, args ...interface{}) {
	f.mu.Lock()
	defer f.mu.Unlock()

	entry := LogEntry{
		Timestamp: time.Now().Format("02-01-2006 15:04:05"),
		Level:     level,
		Message:   fmt.Sprintf(format, args...),
	}

	data, _ := json.Marshal(entry)
	fmt.Fprintln(f.file, string(data))
}

func (f *FileLogger) Warn(args ...interface{}) {
	f.log(WARN, args...)
}

func (f *FileLogger) Warnf(format string, args ...interface{}) {
	f.logf(WARN, format, args...)
}

func (f *FileLogger) Error(args ...interface{}) {
	f.log(ERROR, args...)
}

func (f *FileLogger) Errorf(format string, args ...interface{}) {
	f.logf(ERROR, format, args...)
}

func (f *FileLogger) Info(args ...interface{}) {
	f.log(INFO, args...)
}

func (f *FileLogger) Infof(format string, args ...interface{}) {
	f.logf(INFO, format, args...)
}

func (f *FileLogger) Debug(args ...interface{}) {
	f.log(DEBUG, args...)
}

func (f *FileLogger) Debugf(format string, args ...interface{}) {
	f.logf(DEBUG, format, args...)
}

// CombinedLogger логгер для логирования и в консоль, и в файл
type CombinedLogger struct {
	fileLogger    *FileLogger
	consoleLogger *ConsoleLogger
}

func (c *CombinedLogger) Warn(args ...interface{}) {
	c.fileLogger.Warn(args...)
	c.consoleLogger.Warn(args...)
}

func (c *CombinedLogger) Warnf(format string, args ...interface{}) {
	c.fileLogger.Warnf(format, args...)
	c.consoleLogger.Warnf(format, args...)
}

func (c *CombinedLogger) Error(args ...interface{}) {
	c.fileLogger.Error(args...)
	c.consoleLogger.Error(args...)
}

func (c *CombinedLogger) Errorf(format string, args ...interface{}) {
	c.fileLogger.Errorf(format, args...)
	c.consoleLogger.Errorf(format, args...)
}

func (c *CombinedLogger) Info(args ...interface{}) {
	c.fileLogger.Info(args...)
	c.consoleLogger.Info(args...)
}

func (c *CombinedLogger) Infof(format string, args ...interface{}) {
	c.fileLogger.Infof(format, args...)
	c.consoleLogger.Infof(format, args...)
}

func (c *CombinedLogger) Debug(args ...interface{}) {
	c.fileLogger.Debug(args...)
	c.consoleLogger.Debug(args...)
}

func (c *CombinedLogger) Debugf(format string, args ...interface{}) {
	c.fileLogger.Debugf(format, args...)
	c.consoleLogger.Debugf(format, args...)
}

func (level LogLevel) String() string {
	return string(level)
}
