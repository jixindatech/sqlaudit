package golog

import (
	"encoding/json"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"strings"
	"time"
)

const (
	infoLogName  = "info.log"
	errorLogName = "error.log"
)

type Logger struct {
	zapLogger *zap.Logger
}

var logLevel zapcore.Level

var logger Logger

type requestLogEntry struct {
	RemoteIp    string `json:"retmote_ip"`
	Method      string `json:"method"`
	Uri         string `json:"uri"`
	Status      int    `json:"status"`
	Lentancy    int64  `json:"latency"`
	HumaLatency string `json:"latency_human"`
	BytesIn     int    `json:"bytes_in"`
	BytesOut    int    `json:"bytes_out"`
}

func (l Logger) Write(p []byte) (n int, err error) {
	var logEntry requestLogEntry
	err = json.Unmarshal(p, &logEntry)
	if err != nil {
		return 0, err
	}
	Info("request",
		zap.String("remote_ip", logEntry.RemoteIp),
		zap.String("method", logEntry.Method),
		zap.String("uri", logEntry.Uri),
		zap.Int("status", logEntry.Status),
		zap.Int64("latency", logEntry.Lentancy),
		zap.String("latency_human", logEntry.HumaLatency),
		zap.Int("bytes_in", logEntry.BytesIn),
		zap.Int("bytes_out", logEntry.BytesOut),
	)

	return len(p), nil
}

func Debug(template string, fields ...zap.Field) {
	logger.zapLogger.Debug(template, fields...)
}

func Error(template string, fields ...zap.Field) {
	logger.zapLogger.Error(template, fields...)
}

func Warn(template string, fields ...zap.Field) {
	logger.zapLogger.Warn(template, fields...)
}

func Info(template string, fields ...zap.Field) {
	logger.zapLogger.Info(template, fields...)
}

func Fatal(template string, fields ...zap.Field) {
	logger.zapLogger.Fatal(template, fields...)
}

func setLogLevel(level string) {
	switch strings.ToLower(level) {
	case "debug":
		logLevel = zap.DebugLevel
	case "info":
		logLevel = zap.InfoLevel
	case "warn":
		logLevel = zap.WarnLevel
	case "error":
		logLevel = zap.ErrorLevel
	default:
		logLevel = zap.ErrorLevel
	}
}

func getWriter(filename string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    500, // megabytes
		MaxBackups: 5,
		MaxAge:     30, //days
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func setZapLog(level, path string) error {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "trace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	}

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	var infoWriter, errorWriter io.Writer
	if path == "stdout" {
		infoWriter = os.Stdout
		errorWriter = os.Stdout
	} else {
		infoWriter = getWriter(path + "/" + infoLogName)
		errorWriter = getWriter(path + "/" + errorLogName)
	}

	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(errorWriter), errorLevel),
	)

	//zapLog := zap.New(core, zap.AddCaller())
	logger.zapLogger = zap.New(core)

	return nil
}

func GetLogger() *Logger {
	return &logger
}

func Close() {
	_ = logger.zapLogger.Sync()
}

func SetDefaultZapLog() error {
	level := "info"
	path := "stdout"
	setLogLevel(level)
	return setZapLog(level, path)
}

func InitZapLog(level, path string) error {
	setLogLevel(level)
	return setZapLog(level, path)
}
