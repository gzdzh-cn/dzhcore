package log

import (
	"context"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gbuild"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"

	"github.com/gzdzh-cn/dzhcore/utility/env"
)

var (
	ctx       = gctx.New()
	gLog      = glog.New()
	IsProd    = false
	AppName   = "dzhgo"
	RunLogger *SRunLogger
	IsDesktop = false // 是否为桌面端
	ConfigMap = g.Map{}
)

type SRunLogger struct {
	gLog *glog.Logger
	ctx  context.Context
}

func init() {
	IsDesktop = env.GetCfgWithDefault(ctx, "core.isDesktop", g.NewVar(false)).Bool()
	AppName = env.GetCfgWithDefault(ctx, "core.appName", g.NewVar("dzhgo")).String()
	gbuildData := gbuild.Data()
	if !IsDesktop {
		if _, ok := gbuildData["builtTime"]; ok {
			IsProd = true
		} else {
			IsProd = false
		}
	} else {
		IsProd = env.GetCfgWithDefault(ctx, "core.isProd", g.NewVar(false)).Bool()
	}

	defaultPath := env.GetCfgWithDefault(ctx, "core.gfLogger.path", g.NewVar("path")).String()
	logPath := GetLoggerPath(IsProd, AppName, IsDesktop, defaultPath)
	ConfigMap = g.Map{
		"path":     logPath,
		"level":    env.GetCfgWithDefault(ctx, "core.gfLogger.level", g.NewVar("debug")).String(),
		"stdout":   env.GetCfgWithDefault(ctx, "core.gfLogger.stdout", g.NewVar(true)).Bool(),
		"flags":    env.GetCfgWithDefault(ctx, "core.gfLogger.flags", g.NewVar(44)).Int(),
		"stStatus": env.GetCfgWithDefault(ctx, "core.gfLogger.stStatus", g.NewVar(1)).Int(),
		"stSkip":   env.GetCfgWithDefault(ctx, "core.gfLogger.stSkip", g.NewVar(1)).Int(),
	}

}

func SetLogger() {

	RunLogger = NewRunLogger(ConfigMap) // 初始化 RunLogger 变量
	CfgM := g.Map{
		"path":     ConfigMap["path"],
		"level":    env.GetCfgWithDefault(ctx, "core.gfLogger.level", g.NewVar("debug")).String(),
		"stdout":   env.GetCfgWithDefault(ctx, "core.gfLogger.stdout", g.NewVar(true)).Bool(),
		"flags":    env.GetCfgWithDefault(ctx, "core.gfLogger.flags", g.NewVar(44)).Int(),
		"stStatus": env.GetCfgWithDefault(ctx, "core.gfLogger.stStatus", g.NewVar(1)).Int(),
	}
	g.Log().SetConfigWithMap(CfgM) //
}

func NewRunLogger(configMap g.Map) *SRunLogger {
	gLog.SetConfigWithMap(configMap)
	logger := &SRunLogger{
		gLog: gLog,
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

// 获取日志路径
func GetLoggerPath(isProd bool, appName string, isDesktop bool, defaultPath string) string {

	if !isProd {
		if defaultPath != "" {
			return defaultPath
		}
		return "./data/logs/"
	}

	if isProd && !isDesktop {
		if defaultPath != "" {
			return defaultPath
		}
		return "./data/logs/"
	}

	// 获取适合当前操作系统的基础存储路径
	var basePath string
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			appData = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
		}
		basePath = filepath.Join(appData, appName+"/logs/")
	case "darwin":
		homeDir, _ := os.UserHomeDir()
		basePath = filepath.Join(homeDir, "Library", "Application Support", appName+"/logs/")
	default: // linux 和其他类 Unix 系统
		basePath = "/var/lib/" + appName + "/logs/"
		// 如果不是 root 用户，使用用户目录
		if os.Getuid() != 0 {
			homeDir, _ := os.UserHomeDir()
			basePath = filepath.Join(homeDir, "."+appName+"/logs/")
		}
	}
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		os.MkdirAll(basePath, 0755)
	}
	return basePath

}
