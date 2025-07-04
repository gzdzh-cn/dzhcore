package log

import (
	"context"
	"os"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gzdzh-cn/dzhcore/utility/env"
)

var (
	ctx    = gctx.GetInitCtx()
	Logger *SRunLogger
)

type SRunLogger struct {
	gLog *glog.Logger
	ctx  context.Context
}

func init() {

}

func SetLogger(isProd bool, appName string, isDesktop bool, defaultPath string, configMap g.Map) {
	Logger = NewLogger(configMap) // 初始化 RunLogger 变量
	CfgM := g.Map{
		"path":     configMap["path"],
		"level":    env.GetCfgWithDefault(ctx, "core.gfLogger.level", g.NewVar("debug")).String(),
		"stdout":   env.GetCfgWithDefault(ctx, "core.gfLogger.stdout", g.NewVar(true)).Bool(),
		"flags":    env.GetCfgWithDefault(ctx, "core.gfLogger.flags", g.NewVar(44)).Int(),
		"stStatus": env.GetCfgWithDefault(ctx, "core.gfLogger.stStatus", g.NewVar(1)).Int(),
	}
	g.Log().SetConfigWithMap(CfgM)
}

func NewLogger(configMap g.Map) *SRunLogger {
	logNew := glog.New()
	logNew.SetConfigWithMap(configMap)

	logger := &SRunLogger{
		gLog: logNew,
		ctx:  context.TODO(),
	}
	return logger
}

// Print works like Sprintf.
func (l *SRunLogger) Print(ctx context.Context, message string) {
	l.gLog.Print(ctx, message)
}
func (l *SRunLogger) Printf(ctx context.Context, message string, args ...interface{}) {
	l.gLog.Printf(ctx, message, args...)
}

// Trace level logging. Works like Sprintf.
func (l *SRunLogger) Trace(ctx context.Context, message string) {
	l.gLog.Error(ctx, message)
}

func (l *SRunLogger) Tracef(ctx context.Context, message string, args ...interface{}) {
	l.gLog.Errorf(ctx, message, args...)
}

// Debug level logging. Works like Sprintf.
func (l *SRunLogger) Debug(ctx context.Context, message string) {
	l.gLog.Debug(ctx, message)
}
func (l *SRunLogger) Debugf(ctx context.Context, message string, args ...interface{}) {
	l.gLog.Debugf(ctx, message, args...)
}

// Info level logging. Works like Sprintf.
func (l *SRunLogger) Info(ctx context.Context, message string) {
	l.gLog.Info(ctx, message)
}
func (l *SRunLogger) Infof(ctx context.Context, message string, args ...interface{}) {
	l.gLog.Infof(ctx, message, args...)
}

// Warning level logging. Works like Sprintf.
func (l *SRunLogger) Warning(ctx context.Context, message string) {
	l.gLog.Warning(ctx, message)
}
func (l *SRunLogger) Warningf(ctx context.Context, message string, args ...interface{}) {
	l.gLog.Warningf(ctx, message, args...)
}

// Error level logging. Works like Sprintf.
func (l *SRunLogger) Error(ctx context.Context, message string) {
	l.gLog.Error(ctx, message)
}
func (l *SRunLogger) Errorf(ctx context.Context, message string, args ...interface{}) {
	l.gLog.Errorf(ctx, message, args...)
}

// Fatal level logging. Works like Sprintf.
func (l *SRunLogger) Fatal(ctx context.Context, message string) {
	l.gLog.Error(ctx, message)
	os.Exit(1)
}
func (l *SRunLogger) Fatalf(ctx context.Context, message string, args ...interface{}) {
	l.gLog.Errorf(ctx, message, args...)
	os.Exit(1)
}
