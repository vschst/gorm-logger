package logger

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"strconv"
	"time"
)

type Config struct {
	SlowThreshold	time.Duration
	LogLevel	gormLogger.LogLevel
	SkipErrRecordNotFound bool
	sourceField	string
}

type Logger struct {
	log	*logrus.Logger
	config Config
	traceStr, traceWarnStr, traceErrStr	string
}

func (l *Logger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newLogger := *l
	newLogger.config.LogLevel = level

	return &newLogger
}

func (l *Logger) Info(ctx context.Context, s string, args ...interface{}) {
	if l.config.LogLevel >= gormLogger.Info {
		l.log.WithContext(ctx).Infof(s, args...)
	}
}

func (l *Logger) Warn(ctx context.Context, s string, args ...interface{}) {
	if l.config.LogLevel >= gormLogger.Warn {
		l.log.WithContext(ctx).Warnf(s, args...)
	}
}

func (l *Logger) Error(ctx context.Context, s string, args ...interface{}) {
	if l.config.LogLevel >= gormLogger.Error {
		l.log.WithContext(ctx).Errorf(s, args...)
	}
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.config.LogLevel <= gormLogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := logrus.Fields{}
	if l.config.sourceField != "" {
		fields[l.config.sourceField] = utils.FileWithLineNum()
	}

	elapsedMs := float64(elapsed.Nanoseconds())/1e6
	rowsLog := strconv.FormatInt(rows, 10)
	if rows == -1 {
		rowsLog = "-"
	}

	if err != nil && l.config.LogLevel >= gormLogger.Error && !(errors.Is(err, gorm.ErrRecordNotFound) && l.config.SkipErrRecordNotFound) {
		fields[logrus.ErrorKey] = err
		l.log.WithContext(ctx).WithFields(fields).Errorf(l.traceErrStr, sql, elapsedMs, rowsLog)
	} else if l.config.SlowThreshold != 0 && elapsed > l.config.SlowThreshold && l.config.LogLevel >= gormLogger.Warn {
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.config.SlowThreshold)
		l.log.WithContext(ctx).WithFields(fields).Warnf(l.traceWarnStr, sql, slowLog, elapsedMs, rowsLog)
		return
	} else if l.config.LogLevel == gormLogger.Info {
		l.log.WithContext(ctx).WithFields(fields).Debugf(l.traceStr, sql, elapsedMs, rowsLog)
	}
}

func New(l *logrus.Logger, config Config) *Logger {
	var (
		traceStr     = "%s\n[%.3fms] [rows:%s]"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%s]"
		traceErrStr = "%s\n[%.3fms] [rows:%s]"
	)

	return &Logger{
		log: l,
		config: config,
		traceStr: traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr: traceErrStr,
	}
}
