// Package logger предоставляет контекстный логгер с поддержкой spread ID.
package logger

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tokyobordel/traineepkg/context/trace"
)

const criticalFieldKey = "_critical"

// Formatter - кастомный форматтер для logrus
type Formatter struct {
	TimestampFormat string
}

// Format реализует интерфейс logrus.Formatter
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}

	// Форматируем время
	timestamp := entry.Time.Format(f.TimestampFormat)

	// Форматируем уровень (выравниваем по длине)
	isCritical := false
	if _, ok := entry.Data[criticalFieldKey]; ok {
		isCritical = true
		delete(entry.Data, criticalFieldKey)
	}

	level := strings.ToUpper(entry.Level.String())
	if isCritical {
		level = "CRITICAL"
	}
	if level == "WARN" {
		level = "WARNING"
	}
	level = fmt.Sprintf("%-7s", level) // например, "WARNING"

	// Получаем caller info (пакет, функция, строка)
	caller := getCaller(8) // глубина стека для вызова логгера

	// Получаем UUID из полей, если есть
	spread := ""
	if val, ok := entry.Data["spread"]; ok {
		spread = fmt.Sprintf("%v", val)
		delete(entry.Data, "spread") // чтобы не дублировать в полях
	} else {
		spread = "Spread not found in fields"
	}

	traceback := "<"
	if entry.Level == logrus.ErrorLevel && !isCritical {
		traceback = "\n\n Traceback:\n " + getTraceback() + " <"
	}

	// Формируем основную строку лога
	fmt.Fprintf(b, ">> [%s] [%s] [spread-%s] [%s] msg:> %s %s", timestamp, level, spread, caller, entry.Message, traceback)

	// Сортируем и выводим остальные поля
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(b, " | %s=%v", k, entry.Data[k])
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

// getCaller возвращает строку с пакетом, функцией и строкой кода
func getCaller(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown:0"
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return fmt.Sprintf("%s:%d", file, line)
	}

	fullFuncName := fn.Name() // полный путь пакета + функция
	parts := strings.Split(fullFuncName, "/")
	shortName := parts[len(parts)-1] // оставляем только последний сегмент

	return fmt.Sprintf("%s:%d %s", file, line, shortName)
}

// getTraceback возвращает форматированный stack trace
func getTraceback() string {
	const maxStackDepth = 10
	pcs := make([]uintptr, maxStackDepth)

	// Получаем стек вызовов, пропускаем первые 4 кадра:
	// 0: runtime.Callers
	// 1: getTraceback
	// 2: Format
	// 3: logrus internal
	n := runtime.Callers(4, pcs)

	if n == 0 {
		return "  no stack trace available"
	}

	frames := runtime.CallersFrames(pcs[:n])
	var traceback strings.Builder

	frameNum := 1
	for {
		frame, more := frames.Next()

		// Фильтруем системные пакеты
		if strings.Contains(frame.Function, "runtime.") ||
			strings.Contains(frame.Function, "github.com/sirupsen/logrus") {
			if !more {
				break
			}
			continue
		}

		// Форматируем кадр стека
		shortFile := frame.File
		if lastSlash := strings.LastIndex(shortFile, "/"); lastSlash >= 0 {
			shortFile = shortFile[lastSlash+1:]
		}

		fmt.Fprintf(&traceback, "  %d. %s:%d in %s\n",
			frameNum, shortFile, frame.Line, frame.Function)

		frameNum++
		if !more {
			break
		}
	}

	if traceback.Len() == 0 {
		return "  no relevant stack frames found"
	}

	return traceback.String()
}

func openLogFile(path string) (*os.File, error) {
	if path == "" {
		return nil, fmt.Errorf("путь к файлу логов не может быть пустым")
	}

	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("не удалось создать директорию для логов %q: %w", dir, err)
		}
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл логов %q: %w", path, err)
	}

	return file, nil
}

func newLogger(logFilePath string, useLogFiles bool) (*logrus.Logger, error) {
	logger := logrus.New()

	logger.SetFormatter(&Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logger.SetReportCaller(false)
	logger.SetLevel(logrus.InfoLevel)

	if useLogFiles {
		logFile, err := openLogFile(logFilePath)
		if err != nil {
			return nil, err
		}

		logger.SetOutput(io.MultiWriter(os.Stdout, logFile))
		return logger, nil
	}

	logger.SetOutput(os.Stdout)
	return logger, nil
}

// ContextLogger - обертка над logrus.Logger с обязательным контекстом
type ContextLogger struct {
	logger          *logrus.Logger
	useLogFiles     bool
	criticalLogFile *os.File
}

// NewContextLogger создает новый ContextLogger.
// logFilePath — путь для полных логов приложения.
// criticalLogFilePath — путь для критических логов системы.
// useLogFiles — если true, логи сохраняются в файлы; если false — только в stdout.
func NewContextLogger(logFilePath, criticalLogFilePath string, useLogFiles bool) (*ContextLogger, error) {
	log, err := newLogger(logFilePath, useLogFiles)
	if err != nil {
		return nil, err
	}

	var criticalLogFile *os.File
	if useLogFiles {
		criticalLogFile, err = openLogFile(criticalLogFilePath)
		if err != nil {
			return nil, err
		}
	}

	return &ContextLogger{
		logger:          log,
		useLogFiles:     useLogFiles,
		criticalLogFile: criticalLogFile,
	}, nil
}

// NewContextLoggerFromLogger создает ContextLogger из существующего logrus.Logger
func NewContextLoggerFromLogger(logger *logrus.Logger) *ContextLogger {
	return &ContextLogger{
		logger: logger,
	}
}

// getEntryWithSpread создает logrus.Entry с spread из контекста
func (cl *ContextLogger) getEntryWithSpread(ctx context.Context) *logrus.Entry {
	spread, _ := trace.SpreadFromContext(ctx)
	if spread == "" {
		spread = "Spread not found in context"
	}

	return cl.logger.WithField("spread", spread)
}

func (cl *ContextLogger) writeFormattedLog(formatted []byte) {
	if _, err := cl.logger.Out.Write(formatted); err != nil {
		fmt.Fprintf(os.Stderr, "не удалось записать лог: %v\n", err)
	}
}

func (cl *ContextLogger) writeCriticalLog(formatted []byte) {
	if !cl.useLogFiles || cl.criticalLogFile == nil {
		return
	}

	if _, err := cl.criticalLogFile.Write(formatted); err != nil {
		fmt.Fprintf(os.Stderr, "не удалось записать критический лог в файл: %v\n", err)
	}
}

func (cl *ContextLogger) formatCriticalEntry(ctx context.Context, message string) ([]byte, error) {
	spread, _ := trace.SpreadFromContext(ctx)
	if spread == "" {
		spread = "Spread not found in context"
	}

	entry := &logrus.Entry{
		Logger:  cl.logger,
		Level:   logrus.ErrorLevel,
		Message: message,
		Time:    time.Now(),
		Data: logrus.Fields{
			"spread":         spread,
			criticalFieldKey: true,
		},
	}

	return cl.logger.Formatter.Format(entry)
}

// Info логирует сообщение уровня Info
func (cl *ContextLogger) Info(ctx context.Context, args ...interface{}) {
	cl.getEntryWithSpread(ctx).Info(args...)
}

// Infof логирует форматированное сообщение уровня Info
func (cl *ContextLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	cl.getEntryWithSpread(ctx).Infof(format, args...)
}

// Debug логирует сообщение уровня Debug
func (cl *ContextLogger) Debug(ctx context.Context, args ...interface{}) {
	cl.getEntryWithSpread(ctx).Debug(args...)
}

// Debugf логирует форматированное сообщение уровня Debug
func (cl *ContextLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	cl.getEntryWithSpread(ctx).Debugf(format, args...)
}

// Warn логирует сообщение уровня Warning
func (cl *ContextLogger) Warn(ctx context.Context, args ...interface{}) {
	cl.getEntryWithSpread(ctx).Warn(args...)
}

// Warnf логирует форматированное сообщение уровня Warning
func (cl *ContextLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	cl.getEntryWithSpread(ctx).Warnf(format, args...)
}

// Error логирует сообщение уровня Error
func (cl *ContextLogger) Error(ctx context.Context, args ...interface{}) {
	cl.getEntryWithSpread(ctx).Error(args...)
}

// Errorf логирует форматированное сообщение уровня Error
func (cl *ContextLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	cl.getEntryWithSpread(ctx).Errorf(format, args...)
}

// Criticalf логирует критическое сообщение в stdout, обычный файл и файл критических логов
func (cl *ContextLogger) Criticalf(ctx context.Context, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)

	formatted, err := cl.formatCriticalEntry(ctx, message)
	if err != nil {
		fmt.Fprintf(os.Stderr, "не удалось отформатировать критический лог: %v\n", err)
		return
	}

	cl.writeFormattedLog(formatted)
	cl.writeCriticalLog(formatted)
}

// Fatal логирует сообщение уровня Fatal
func (cl *ContextLogger) Fatal(ctx context.Context, args ...interface{}) {
	cl.getEntryWithSpread(ctx).Fatal(args...)
}

// Fatalf логирует форматированное сообщение уровня Fatal
func (cl *ContextLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	cl.getEntryWithSpread(ctx).Fatalf(format, args...)
}

// Panic логирует сообщение уровня Panic
func (cl *ContextLogger) Panic(ctx context.Context, args ...interface{}) {
	cl.getEntryWithSpread(ctx).Panic(args...)
}

// Panicf логирует форматированное сообщение уровня Panic
func (cl *ContextLogger) Panicf(ctx context.Context, format string, args ...interface{}) {
	cl.getEntryWithSpread(ctx).Panicf(format, args...)
}

// WithField добавляет поле к логгеру и возвращает новый entry с контекстом
func (cl *ContextLogger) WithField(ctx context.Context, key string, value interface{}) *logrus.Entry {
	return cl.getEntryWithSpread(ctx).WithField(key, value)
}

// WithFields добавляет несколько полей к логгеру и возвращает новый entry с контекстом
func (cl *ContextLogger) WithFields(ctx context.Context, fields logrus.Fields) *logrus.Entry {
	return cl.getEntryWithSpread(ctx).WithFields(fields)
}

// SetLevel устанавливает уровень логирования
func (cl *ContextLogger) SetLevel(level logrus.Level) {
	cl.logger.SetLevel(level)
}

// GetLevel возвращает текущий уровень логирования
func (cl *ContextLogger) GetLevel() logrus.Level {
	return cl.logger.GetLevel()
}
