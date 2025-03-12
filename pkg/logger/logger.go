package logger

import (
	"context"
	"fmt"

	configs "backend/pkg/config"
	app_middlewares "backend/pkg/middlewares"
	"os"
	"runtime"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Fields map[string]any

type Logger interface {
	Print(args ...any)
	Debug(args ...any)
	Debugln(args ...any)
	Debugf(format string, args ...any)

	Info(args ...any)
	Infoln(args ...any)
	Infof(format string, args ...any)

	Warn(args ...any)
	Warnln(args ...any)
	Warnf(format string, args ...any)

	Error(args ...any)
	Errorln(args ...any)
	Errorf(format string, args ...any)

	Fatal(args ...any)
	Fatalln(args ...any)
	Fatalf(format string, args ...any)

	Panic(args ...any)
	Panicln(args ...any)
	Panicf(format string, args ...any)

	With(key string, value any) Logger
	Withs(fields Fields) Logger
	WithSrc() Logger
	WithContext(ctx context.Context) Logger
	GetLevel() string
}

type appLogger struct {
	logger *zap.Logger
	level  zap.AtomicLevel
}

var loggerLevelMap = map[string]zapcore.Level{
	"Debug": zapcore.DebugLevel,
	"Info":  zapcore.InfoLevel,
	"Warn":  zapcore.WarnLevel,
	"Error": zapcore.ErrorLevel,
	"Panic": zapcore.PanicLevel,
	"Fatal": zapcore.FatalLevel,
}

func NewLogger(appConfig *configs.AppConfig) Logger {
	config := appConfig.Logger
	currentDate := time.Now().Format("2006-01-02")
	logFilePath := fmt.Sprintf(config.FileName, appConfig.Server.ServiceName, currentDate)

	writer := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}

	level := zap.NewAtomicLevelAt(loggerLevelMap[config.Level])
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(writer), zapcore.AddSync(os.Stdout)),
		level)

	return &appLogger{logger: zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel)), level: level}
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}
func (l *appLogger) Print(args ...any) {
	l.logger.Info(fmt.Sprint(args...))
}

func (l *appLogger) Debug(args ...any) {
	l.logger.Debug(fmt.Sprint(args...))
}

func (l *appLogger) Debugln(args ...any) {
	l.logger.Debug(fmt.Sprintln(args...))
}

func (l *appLogger) Debugf(format string, args ...any) {
	l.logger.Debug(fmt.Sprintf(format, args...))
}

func (l *appLogger) Info(args ...any) {
	l.logger.Info(fmt.Sprint(args...))
}

func (l *appLogger) Infoln(args ...any) {
	l.logger.Info(fmt.Sprintln(args...))
}

func (l *appLogger) Infof(format string, args ...any) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

func (l *appLogger) Warn(args ...any) {
	l.logger.Warn(fmt.Sprint(args...))
}

func (l *appLogger) Warnln(args ...any) {
	l.logger.Warn(fmt.Sprintln(args...))
}

func (l *appLogger) Warnf(format string, args ...any) {
	l.logger.Warn(fmt.Sprintf(format, args...))
}

func (l *appLogger) Error(args ...any) {
	l.logger.Error(fmt.Sprint(args...))
}

func (l *appLogger) Errorln(args ...any) {
	l.logger.Error(fmt.Sprintln(args...))
}

func (l *appLogger) Errorf(format string, args ...any) {
	l.logger.Error(fmt.Sprintf(format, args...))
}

func (l *appLogger) Fatal(args ...any) {
	l.logger.Fatal(fmt.Sprint(args...))
}

func (l *appLogger) Fatalln(args ...any) {
	l.logger.Fatal(fmt.Sprintln(args...))
}

func (l *appLogger) Fatalf(format string, args ...any) {
	l.logger.Fatal(fmt.Sprintf(format, args...))
}

func (l *appLogger) Panic(args ...any) {
	l.logger.Panic(fmt.Sprint(args...))
}

func (l *appLogger) Panicln(args ...any) {
	l.logger.Panic(fmt.Sprintln(args...))
}

func (l *appLogger) Panicf(format string, args ...any) {
	l.logger.Panic(fmt.Sprintf(format, args...))
}

func (l *appLogger) With(key string, value any) Logger {
	return &appLogger{logger: l.logger.With(zap.Any(key, value)), level: l.level}
}

func (l *appLogger) Withs(fields Fields) Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return &appLogger{logger: l.logger.With(zapFields...), level: l.level}
}

func (l *appLogger) WithSrc() Logger {
	_, file, line, _ := runtime.Caller(1)
	return &appLogger{logger: l.logger.With(zap.String("source", fmt.Sprintf("%s:%d", file, line))), level: l.level}
}

func (l *appLogger) GetLevel() string {
	return l.level.String()
}
func (l *appLogger) WithContext(ctx context.Context) Logger {
	key := app_middlewares.CtxKey(echo.HeaderXCorrelationID)
	correlationId := ctx.Value(key)
	if correlationId != nil {
		return &appLogger{logger: l.logger.With(zap.String("correlationId", correlationId.(string))), level: l.level}
	}
	return l
}
