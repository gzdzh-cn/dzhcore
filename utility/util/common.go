package util

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
)

var ctx = gctx.GetInitCtx()

// GetValueOrDefault
//
//	@Description: 给定一个interface，如果是nil或空，则返回false
//	@param value
//	@param defaultValue
//	@return interface{}
func GetValueOrDefault(value interface{}) bool {

	if value == nil {
		return false
	}

	// 使用类型断言判断是否为空值
	switch v := value.(type) {
	case string:
		if v == "" {
			return false
		}
	case int:
		if v == 0 {
			return false
		}
	case []interface{}:
		if len(v) == 0 {
			return false
		}
	// 可以根据需要添加更多类型的检查
	default:
		// 其他类型不为空的情况
		return true
	}

	// 如果不是空值，返回原值
	return true
}

// GetDatabaseVersion
//
//	@Description: 获取mysql版本
//	@return string
func GetDBVersion() string {

	dbType := g.Cfg().MustGet(ctx, "database.default.type")

	type result struct {
		Version string `json:"version"`
	}
	var res *result
	query := ""
	switch strings.ToLower(dbType.String()) {
	case "mysql":
		query = "SELECT VERSION() as version"
	case "sqlite":
		query = "SELECT sqlite_version() as version"
	case "pgsql":
		query = "SELECT version() as version"
	default:
		g.Log().Warningf(ctx, "unsupported database type for version retrieval: %s", dbType.String())
		return ""
	}

	err := g.DB().Raw(query).Scan(&res)

	if err != nil {
		g.Log().Error(ctx, err.Error())
		return ""
	}

	return res.Version
}

// 配置信息
func GetConfig(v ...interface{}) string {

	if len(v) == 1 {
		data, err := g.Cfg().Get(ctx, gconv.String(v[0]))
		if err != nil {
			g.Log().Error(ctx, err.Error())
			return ""
		}
		return data.String()
	}

	return ""
}

// 获取 env 变量
func Getenv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// 获取日志路径
func GetLoggerPath(isProd bool, appName string) string {
	if !isProd {
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
	return basePath
}
