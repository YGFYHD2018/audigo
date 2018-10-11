package util

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	instance = newLogger()
)

type Log interface {
	Debug(...interface{}) string
	Info(...interface{}) string
	Warn(...interface{}) string
	Error(...interface{}) string
	Fatal(...interface{}) string
}

type logger struct {
	zap *zap.Logger
}

func newLogger() Log {
	conf := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	console := zapcore.NewCore(
		zapcore.NewConsoleEncoder(conf),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)
	file := newFileLogger("log/audigo.log", conf)
	l := zap.New(zapcore.NewTee(
		console,
		file,
	))
	return &logger{l}
}

func newFileLogger(path string, conf zapcore.EncoderConfig) zapcore.Core {
	f, _ := os.Create(path)
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(conf),
		zapcore.AddSync(f),
		zapcore.DebugLevel,
	)
	return fileCore
}

func GetLogger() Log {
	return instance
}

func (l *logger) Debug(v ...interface{}) string {
	l.zap.Debug(fmt.Sprintln(v...))
	return ""
}

func (l *logger) Info(v ...interface{}) string {
	l.zap.Info(fmt.Sprintln(v...))
	return ""
}

func (l *logger) Warn(v ...interface{}) string {
	l.zap.Warn(fmt.Sprintln(v...))
	return ""
}

func (l *logger) Error(v ...interface{}) string {
	l.zap.Error(fmt.Sprintln(v...))
	return ""
}

func (l *logger) Fatal(v ...interface{}) string {
	l.zap.Fatal(fmt.Sprintln(v...))
	return ""
}
