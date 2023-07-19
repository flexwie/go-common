package logger

import (
	"context"
	"time"

	"github.com/charmbracelet/log"
	glogger "gorm.io/gorm/logger"
)

type GormLogger struct {
	Logger *log.Logger
}

func NewGormLogger(logger *log.Logger) *GormLogger {
	return &GormLogger{
		Logger: logger.WithPrefix("gorm"),
	}
}

func (g *GormLogger) LogMode(ll glogger.LogLevel) glogger.Interface {
	switch ll {
	case glogger.Silent:
		g.Logger.SetLevel(log.DebugLevel)
	case glogger.Info:
		g.Logger.SetLevel(log.InfoLevel)
	case glogger.Warn:
		g.Logger.SetLevel(log.WarnLevel)
	case glogger.Error:
		g.Logger.SetLevel(log.ErrorLevel)
	}

	return g
}

func (g *GormLogger) Info(ctx context.Context, s string, opts ...interface{}) {
	g.Logger.Infof(s, opts...)
}

func (g *GormLogger) Warn(ctx context.Context, s string, opts ...interface{}) {
	g.Logger.Warnf(s, opts...)
}

func (g *GormLogger) Error(ctx context.Context, s string, opts ...interface{}) {
	g.Logger.Errorf(s, opts...)
}

func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rows := fc()
	g.Logger.Info("sql", "query", sql, "affected", rows)
}
