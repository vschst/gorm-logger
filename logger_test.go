package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	gormLogger "gorm.io/gorm/logger"
	"testing"
	"time"
)

type LogRow struct {
	Level string `json:"level"`
	Message	string	`json:"msg"`
}

func setupLogger() (context.Context, *bytes.Buffer, *Logger) {
	ctx := context.Background()

	buf := new(bytes.Buffer)
	log := logrus.New()
	log.SetOutput(buf)
	log.SetFormatter(&logrus.JSONFormatter{})

	l := New(log, Config{
		SlowThreshold: time.Second,
		LogLevel: gormLogger.Info,
	})

	return ctx, buf, l
}

func TestLoggerInfo(t *testing.T) {
	ctx, buf, l := setupLogger()

	l.Info(ctx, "test")

	var logData LogRow
	if err := json.Unmarshal(buf.Bytes(), &logData); err != nil {
		t.Errorf("failed to parse log row: %v", err)
	}

	levelWant := "info"
	if logData.Level != "info" {
		t.Errorf("Info() test failed, unexpected log level '%v', expected '%v'", logData.Level, levelWant)
	}

	messageWant := "test"
	if logData.Message != messageWant {
		t.Errorf("Info() test failed, unexpected log msg '%v', expected '%v'", logData.Message, messageWant)
	}
}
