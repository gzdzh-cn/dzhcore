package logger

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
)

func Debug(ctx context.Context, v ...interface{}) {
	g.Log().Debug(ctx, v...)
}

func Debugf(ctx context.Context, format string, v ...interface{}) {
	g.Log().Debugf(ctx, format, v...)
}

func Warning(ctx context.Context, v ...interface{}) {
	g.Log().Warning(ctx, v...)
}

func Warningf(ctx context.Context, format string, v ...interface{}) {
	g.Log().Warningf(ctx, format, v...)
}

func Info(ctx context.Context, v ...interface{}) {
	g.Log().Info(ctx, v...)
}

func Infof(ctx context.Context, format string, v ...interface{}) {
	g.Log().Infof(ctx, format, v...)
}

func Error(ctx context.Context, v ...interface{}) {
	g.Log().Error(ctx, v...)
}

func Errorf(ctx context.Context, format string, v ...interface{}) {
	g.Log().Errorf(ctx, format, v...)
}
