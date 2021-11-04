package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	gormLogger "gorm.io/gorm/logger"
	"reflect"
	"testing"
	"time"
)

type LogRow struct {
	Level string `json:"level"`
	Message	string	`json:"msg"`
	Module	string	`json:"module"`
	Data	[]interface{}	`json:"data"`
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

func compareLogData(x, y []interface{}) (bool, error) {
	xBuf, err := json.Marshal(x)
	if err != nil {
		return false, fmt.Errorf("failed to marshal first value: %v", err)
	}

	yBuf, err := json.Marshal(y)
	if err != nil {
		return false, fmt.Errorf("failed to marshall second value: %v", err)
	}

	return bytes.Equal(xBuf, yBuf), nil
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

		if logData.Message != infoTest.msg {
			t.Errorf("Info() output test failed, unexpected log msg '%v', expected '%v'", logData.Message, infoTest.msg)
		}

		if logData.Module != moduleName {
			t.Errorf("Info() output test failed, unexpected module name '%v', expected '%v'", logData.Message, moduleName)
		}

		if infoTest.args != nil {
			isEqual, err := compareLogData(logData.Data, infoTest.args)

			if err != nil {
				t.Errorf("Info() output test failed, failed to compare log data, compared values %v and %v", logData.Data, infoTest.args)
			}

			if !isEqual {
				t.Errorf("Info() output test failed, compared values %v and %v don't match", logData.Data, infoTest.args)
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

		if logData.Message != infoTest.msg {
			t.Errorf("Warn() output test failed, unexpected log msg '%v', expected '%v'", logData.Message, infoTest.msg)
		}

		if logData.Module != moduleName {
			t.Errorf("Warn() output test failed, unexpected module name '%v', expected '%v'", logData.Message, moduleName)
		}

		if infoTest.args != nil {
			isEqual, err := compareLogData(logData.Data, infoTest.args)

			if err != nil {
				t.Errorf("Warn() output test failed, failed to compare log data, compared values %v and %v", logData.Data, infoTest.args)
			}

			if !isEqual {
				t.Errorf("Warn() output test failed, compared values %v and %v don't match", logData.Data, infoTest.args)
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

		if logData.Message != infoTest.msg {
			t.Errorf("Error() output test failed, unexpected log msg '%v', expected '%v'", logData.Message, infoTest.msg)
		}

		if logData.Module != moduleName {
			t.Errorf("Error() output test failed, unexpected module name '%v', expected '%v'", logData.Message, moduleName)
		}

		if infoTest.args != nil {
			isEqual, err := compareLogData(logData.Data, infoTest.args)

			if err != nil {
				t.Errorf("Error() output test failed, failed to compare log data, compared values %v and %v", logData.Data, infoTest.args)
			}

			if !isEqual {
				t.Errorf("Error() output test failed, compared values %v and %v don't match", logData.Data, infoTest.args)
			}
		}

		buf.Reset()
	}
}

func TestLoggerTrace(t *testing.T) {
	ctx, buf, _, l := setupLogger(Config{
		SourceField: "test",
	})

	fc := func() (string, int64) {
		return "SELECT 1 FROM table", 0
	}

	l.Trace(ctx, time.Time{}, fc, nil)

	t.Errorf("%v", buf.String())
}

