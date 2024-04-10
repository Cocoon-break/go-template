package db

import (
	"context"
	"errors"
	"time"

	"go-template/pkg/zlog"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormLogger struct {
	LogLevel                  logger.LogLevel
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
}

func NewGormLogger(slowThreshold time.Duration) logger.Interface {
	return &GormLogger{
		LogLevel:      logger.Warn,
		SlowThreshold: slowThreshold,
	}
}

// Error implements logger.Interface.
func (g GormLogger) Error(_ context.Context, str string, args ...interface{}) {
	if g.LogLevel < logger.Error {
		return
	}
	zlog.Errorf(str, args...)
}

// Info implements logger.Interface.
func (g GormLogger) Info(_ context.Context, str string, args ...interface{}) {
	if g.LogLevel < logger.Info {
		return
	}
	zlog.Errorf(str, args...)
}

// Warn implements logger.Interface.
func (g GormLogger) Warn(_ context.Context, str string, args ...interface{}) {
	if g.LogLevel < logger.Info {
		return
	}
	zlog.Errorf(str, args...)
}

// LogMode implements logger.Interface.
func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := g
	newLogger.LogLevel = level
	return newLogger
}

// Trace implements logger.Interface.
func (g GormLogger) Trace(_ context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if g.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && g.LogLevel >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !g.IgnoreRecordNotFoundError):
		sql, rows := fc()
		zlog.Error("exec_sql", zap.String("sql_elapsed", elapsed.String()), zap.Int64("rows", rows), zap.String("sql", sql), zap.Error(err))
	case elapsed > g.SlowThreshold && g.SlowThreshold != 0 && g.LogLevel >= logger.Warn:
		sql, rows := fc()
		zlog.Warn("slow_sql", zap.String("sql_elapsed", elapsed.String()), zap.Int64("rows", rows), zap.String("sql", sql))
	case g.LogLevel >= logger.Info:
		sql, rows := fc()
		zlog.Debug("gorm", zap.String("sql_elapsed", elapsed.String()), zap.Int64("rows", rows), zap.String("sql", sql))
	}
}
