package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	gormLogger "gorm.io/gorm/logger"
	"reflect"
	"strings"
	"testing"
	"time"
)

type LogRow struct {
	Level string `json:"level"`
	Message	string	`json:"msg"`
	Module	string	`json:"module"`
	Test	string `json:"test"`
	Error	string	`json:"error"`
}

func setupLogger(config Config) (context.Context, *bytes.Buffer, *logrus.Logger, *Logger) {
	ctx := context.Background()

	buf := new(bytes.Buffer)
	log := logrus.New()
	log.SetOutput(buf)
	log.SetFormatter(&logrus.JSONFormatter{})

	l := New(log, config)

	return ctx, buf, log, l
}

func compareExtraArgs(msg string, args []interface{}) bool {
	for _, arg := range args {
		if !strings.Contains(msg, fmt.Sprintf("%T=%v", arg, arg)) {
			return false
		}
	}

	return true
}

func TestInstance(t *testing.T) {
	_, _, log, l := setupLogger(Config{})
	want := &Logger{
		log: log,
		config: Config{
			ModuleName: "gorm",
			LogLevel: gormLogger.Info,
		},
	}

	if !reflect.DeepEqual(l, want) {
		t.Errorf("New() test failed, unexpected instance %v, expected %v", l, want)
	}
}

func TestLoggerGetLogger(t *testing.T) {
	_, _, log, l := setupLogger(Config{})

	got := l.GetLogger()
	if !reflect.DeepEqual(got, log) {
		t.Errorf("GetLogger() test failed, unexpected logrus instance %v, expected %v", got, log)
	}
}

func TestLoggerLogMode(t *testing.T) {
	_, _, log, l := setupLogger(Config{
		LogLevel: gormLogger.Info,
	})
	want := &Logger{
		log: log,
		config: Config{
			ModuleName: "gorm",
			LogLevel: gormLogger.Warn,
		},
	}

	got := l.LogMode(gormLogger.Warn)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("LogMode() test failed, unexpected instance %v, expected %v", got, want)
	}
}

var loggerOutputTests = []struct{
	msg	string
	args []interface{}
} {
	{
		msg: "test msg",
	},
	{
		msg: "test args",
		args: []interface{}{123, []int{1, 2}, "str", map[string]string{"key": "value"}},
	},
}

func TestLoggerInfo(t *testing.T) {
	moduleName := "gorm-info-logger"
	ctx, buf, _, l := setupLogger(Config{
		LogLevel: gormLogger.Info,
		ModuleName: moduleName,
	})

	l.Info(ctx, "info msg")

	var logData LogRow
	if err := json.Unmarshal(buf.Bytes(), &logData); err != nil {
		t.Errorf("Info() test error, failed to parse log row: %v", err)
	}

	levelWant := "info"
	if logData.Level != levelWant {
		t.Errorf("Info() test failed, unexpected log level '%v', expected '%v'", logData.Level, levelWant)
	}

	buf.Reset()

	for _, infoTest := range loggerOutputTests {
		l.Info(ctx, infoTest.msg, infoTest.args...)

		if err := json.Unmarshal(buf.Bytes(), &logData); err != nil {
			t.Errorf("Info() output test error, failed to parse log row: %v, got string: %v", err, buf.String())
		}

		if !strings.HasPrefix(logData.Message, infoTest.msg) {
			t.Errorf("Info() output test failed, unexpected start of log msg '%v', expected '%v'", logData.Message, infoTest.msg)
		}

		if logData.Module != moduleName {
			t.Errorf("Info() output test failed, unexpected module name '%v', expected '%v'", logData.Message, moduleName)
		}

		if infoTest.args != nil {
			isEqual := compareExtraArgs(logData.Message, infoTest.args)

			if !isEqual {
				t.Errorf("Info() output test failed, log msg %v doesn't contain extra args %v", logData.Message, infoTest.args)
			}
		}

		buf.Reset()
	}
}

func TestLoggerWarn(t *testing.T) {
	moduleName := "gorm-warn-logger"
	ctx, buf, _, l := setupLogger(Config{
		LogLevel: gormLogger.Warn,
		ModuleName: moduleName,
	})

	l.Warn(ctx, "warn msg")

	var logData LogRow
	if err := json.Unmarshal(buf.Bytes(), &logData); err != nil {
		t.Errorf("Warn() test error, failed to parse log row: %v", err)
	}

	levelWant := "warning"
	if logData.Level != levelWant {
		t.Errorf("Warn() test failed, unexpected log level '%v', expected '%v'", logData.Level, levelWant)
	}

	buf.Reset()

	for _, infoTest := range loggerOutputTests {
		l.Warn(ctx, infoTest.msg, infoTest.args...)

		if err := json.Unmarshal(buf.Bytes(), &logData); err != nil {
			t.Errorf("Warn() output test error, failed to parse log row: %v, got string: %v", err, buf.String())
		}

		if !strings.HasPrefix(logData.Message, infoTest.msg) {
			t.Errorf("Warn() output test failed, unexpected start of log msg '%v', expected '%v'", logData.Message, infoTest.msg)
		}

		if logData.Module != moduleName {
			t.Errorf("Warn() output test failed, unexpected module name '%v', expected '%v'", logData.Message, moduleName)
		}

		if infoTest.args != nil {
			isEqual := compareExtraArgs(logData.Message, infoTest.args)

			if !isEqual {
				t.Errorf("Warn() output test failed, log msg %v doesn't contain extra args %v", logData.Message, infoTest.args)
			}
		}

		buf.Reset()
	}
}

func TestLoggerError(t *testing.T) {
	moduleName := "gorm-error-logger"
	ctx, buf, _, l := setupLogger(Config{
		LogLevel:   gormLogger.Error,
		ModuleName: moduleName,
	})

	l.Error(ctx, "error msg")

	var logData LogRow
	if err := json.Unmarshal(buf.Bytes(), &logData); err != nil {
		t.Errorf("Error() test error, failed to parse log row: %v", err)
	}

	levelWant := "error"
	if logData.Level != levelWant {
		t.Errorf("Error() test failed, unexpected log level '%v', expected '%v'", logData.Level, levelWant)
	}

	buf.Reset()

	for _, infoTest := range loggerOutputTests {
		l.Error(ctx, infoTest.msg, infoTest.args...)

		if err := json.Unmarshal(buf.Bytes(), &logData); err != nil {
			t.Errorf("Error() output test error, failed to parse log row: %v, got string: %v", err, buf.String())
		}

		if !strings.HasPrefix(logData.Message, infoTest.msg) {
			t.Errorf("Error() output test failed, unexpected start of log msg '%v', expected '%v'", logData.Message, infoTest.msg)
		}

		if logData.Module != moduleName {
			t.Errorf("Error() output test failed, unexpected module name '%v', expected '%v'", logData.Message, moduleName)
		}

		if infoTest.args != nil {
			isEqual := compareExtraArgs(logData.Message, infoTest.args)

			if !isEqual {
				t.Errorf("Error() output test failed, log msg %v doesn't contain extra args %v", logData.Message, infoTest.args)
			}
		}

		buf.Reset()
	}
}

func TestLoggerTraceError(t *testing.T) {
	sourceField := "test"
	ctx, buf, _, l := setupLogger(Config{
		SourceField: sourceField,
	})

	fc := func() (string, int64) {
		return "", -1
	}

	errorText := "custom_error"
	l.Trace(ctx, time.Now(), fc, errors.New(errorText))

	var logData LogRow
	if err := json.Unmarshal(buf.Bytes(), &logData); err != nil {
		t.Errorf("Trace() error checking error, failed to parse log row: %v", err)
	}

	levelWant := "error"
	if logData.Level != levelWant {
		t.Errorf("Trace() error checking error, unexpected log level '%v', expected '%v'", logData.Level, levelWant)
	}

	if logData.Test == "" {
		t.Errorf("Trace() error checking error, source key '%v' doesn't exists on log row: %v", sourceField, buf.String())
	}

	if logData.Error != errorText {
		t.Errorf("Trace() error checking error, invalid error field value on log, got '%v', expected: %v", logData.Error, errorText)
	}
}

func TestLoggerTraceSlowThreshold(t *testing.T) {
	sourceField := "test"
	ctx, buf, _, l := setupLogger(Config{
		SourceField: sourceField,
		SlowThreshold: time.Second,
	})

	sqlQuery := "SELECT 1 FROM table"
	rows := int64(1e6)
	fc := func() (string, int64) {
		return sqlQuery, rows
	}

	now := time.Now()
	begin := now.Add((-1) * time.Minute)
	l.Trace(ctx, begin, fc, nil)

	var logData LogRow
	if err := json.Unmarshal(buf.Bytes(), &logData); err != nil {
		t.Errorf("Trace() slow threshold test error, failed to parse log row: %v", err)
	}

	levelWant := "warning"
	if logData.Level != levelWant {
		t.Errorf("Trace() slow threshold test error, unexpected log level '%v', expected '%v'", logData.Level, levelWant)
	}

	if !strings.Contains(logData.Message, sqlQuery) {
		t.Errorf("Trace() slow threshold test error, unexpected SQL query in message: %v, expected: %v", logData.Message, sqlQuery)
	}

	if !strings.Contains(logData.Message, "SLOW SQL") {
		t.Errorf("Trace() slow threshold test error, there is not info about SLOW SQL, got message: %v", logData.Message)
	}

	if !strings.Contains(logData.Message, fmt.Sprintf("rows:%d", rows)) {
		t.Errorf("Trace() slow threshold test error, invalid rows info, expected 'rows:%d', got message: %v", rows, logData.Message)
	}

	if logData.Test == "" {
		t.Errorf("Trace() slow threshold test error, source key '%v' doesn't exists on log row: %v", sourceField, buf.String())
	}
}

func TestLoggerTraceInfo(t *testing.T) {
	ctx, buf, _, l := setupLogger(Config{})

	fc := func() (string, int64) {
		return "", -1
	}

	l.Trace(ctx, time.Now(), fc, nil)

	var logData LogRow
	if err := json.Unmarshal(buf.Bytes(), &logData); err != nil {
		t.Errorf("Trace() info test error, failed to parse log row: %v", err)
	}

	levelWant := "info"
	if logData.Level != levelWant {
		t.Errorf("Trace() info test error, unexpected log level '%v', expected '%v'", logData.Level, levelWant)
	}

	buf.Reset()

	nl := l.LogMode(gormLogger.Silent)

	nl.Trace(ctx, time.Now(), fc, nil)

	silentLog := buf.String()
	if silentLog != "" {
		t.Errorf("Trace() info test error, with a silent mode, the following log is output: %v", silentLog)
	}
}

