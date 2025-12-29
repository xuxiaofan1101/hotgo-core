package apollo

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
)

type GoFrameLogger struct {
	ctx    context.Context
	logger *glog.Logger
}

func NewGoFrameLogger(logger *glog.Logger) *GoFrameLogger {
	if logger == nil {
		logger = g.Log() //默认使用全局 Logger
	}
	return &GoFrameLogger{
		ctx:    gctx.New(),
		logger: logger,
	}
}

func (l *GoFrameLogger) Debugf(format string, params ...interface{}) {
	l.logger.Debugf(nil, format, params...)
}

func (l *GoFrameLogger) Infof(format string, params ...interface{}) {
	l.logger.Infof(nil, format, params...)
}

func (l *GoFrameLogger) Warnf(format string, params ...interface{}) {
	l.logger.Warningf(nil, format, params...)
}

func (l *GoFrameLogger) Errorf(format string, params ...interface{}) {
	l.logger.Errorf(nil, format, params...)
}

func (l *GoFrameLogger) Debug(v ...interface{}) {
	l.logger.Debug(nil, v...)
}

func (l *GoFrameLogger) Info(v ...interface{}) {
	l.logger.Info(nil, v...)
}

func (l *GoFrameLogger) Warn(v ...interface{}) {
	l.logger.Warning(nil, v...)
}

func (l *GoFrameLogger) Error(v ...interface{}) {
	l.logger.Error(nil, v...)
}
