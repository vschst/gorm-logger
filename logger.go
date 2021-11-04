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

const (
	traceStr     = "%s\n[%.3fms] [rows:%s]"
	traceWarnStr = "%s %s\n[%.3fms] [rows:%s]"
	traceErrStr = "%s\n[%.3fms] [rows:%s]"
)

type Config struct {
	SlowThreshold	time.Duration
	LogLevel	gormLogger.LogLevel
	SkipErrRecordNotFound bool
	SourceField	string
	ModuleName	string
}

type Logger struct {
	log	*logrus.Logger
	config Config
}


func (l *Logger) createFields(args []interface{}) logrus.Fields {
	fields := logrus.Fields{
		"module": l.config.ModuleName,
	}

	if len(args) > 0 {
		fields["data"] = args
	}

	return fields
}

func (l* Logger) createEntry(ctx context.Context, args ...interface{}) *logrus.Entry {
	return l.log.WithContext(ctx).WithFields(l.createFields(args))
}

func (l *Logger) GetLogger() *logrus.Logger {
	return l.log
}

func (l *Logger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newLogger := *l
	newLogger.config.LogLevel = level

	return &newLogger
}

func (l *Logger) Info(ctx context.Context, s string, args ...interface{}) {
	if l.config.LogLevel >= gormLogger.Info {
		l.createEntry(ctx, args...).Infof(s)
	}
}

func (l *Logger) Warn(ctx context.Context, s string, args ...interface{}) {
	if l.config.LogLevel >= gormLogger.Warn {
		l.createEntry(ctx, args...).Warnf(s)
	}
}

func (l *Logger) Error(ctx context.Context, s string, args ...interface{}) {
	if l.config.LogLevel >= gormLogger.Error {
		l.createEntry(ctx, args...).Errorf(s)
	}
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.config.LogLevel <= gormLogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := logrus.Fields{}
	if l.config.SourceField != "" {
		fields[l.config.SourceField] = utils.FileWithLineNum()
	}

	elapsedMs := float64(elapsed.Nanoseconds())/1e6
	rowsLog := strconv.FormatInt(rows, 10)
	if rows == -1 {
		rowsLog = "-"
	}

	switch {
	case err != nil && l.config.LogLevel >= gormLogger.Error && !(errors.Is(err, gorm.ErrRecordNotFound) && l.config.SkipErrRecordNotFound):
		fields[logrus.ErrorKey] = err.Error()
		l.log.WithContext(ctx).WithFields(fields).Errorf(traceErrStr, sql, elapsedMs, rowsLog)
	case l.config.SlowThreshold != 0 && elapsed > l.config.SlowThreshold && l.config.LogLevel >= gormLogger.Warn:
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.config.SlowThreshold)
		l.log.WithContext(ctx).WithFields(fields).Warnf(traceWarnStr, sql, slowLog, elapsedMs, rowsLog)
	case l.config.LogLevel == gormLogger.Info:
		l.log.WithContext(ctx).WithFields(fields).Infof(traceStr, sql, elapsedMs, rowsLog)
	default:
		l.log.WithContext(ctx).WithFields(fields).Debugf(traceStr, sql, elapsedMs, rowsLog)
	}
}

func New(l *logrus.Logger, config Config) *Logger {
	if config.ModuleName == "" {
		config.ModuleName = "gorm"
	}

	if config.LogLevel == 0 {
		config.LogLevel = gormLogger.Info
	}

	return &Logger{
		log: l,
		config: config,
	}
}
