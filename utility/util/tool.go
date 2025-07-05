package util

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/gzdzh-cn/dzhcore/coreconfig"
	"github.com/gzdzh-cn/dzhcore/envconfig"
	"github.com/gzdzh-cn/dzhcore/utility/defineType"
	"github.com/shopspring/decimal"
)

// ToolUtil 结构体，包含 ctx 字段
type ToolUtil struct {
	ctx context.Context
}

// NewToolUtil 构造函数
func NewToolUtil() *ToolUtil {
	return &ToolUtil{
		ctx: gctx.GetInitCtx(),
	}
}

// GetValueOrDefault
//
//	@Description: 给定一个interface，如果是nil或空，则返回false
//	@param value
//	@param defaultValue
//	@return interface{}
func (t *ToolUtil) GetValueOrDefault(value interface{}) bool {
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

// GetDBVersion
//
//	@Description: 获取mysql版本
//	@return string
func (t *ToolUtil) GetDBVersion() string {
	dbType := g.Cfg().MustGet(t.ctx, "database.default.type")
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
		g.Log().Warningf(t.ctx, "unsupported database type for version retrieval: %s", dbType.String())
		return ""
	}

	err := g.DB().Raw(query).Scan(&res)
	if err != nil {
		g.Log().Error(t.ctx, err.Error())
		return ""
	}
	return res.Version
}

// 配置信息
func (t *ToolUtil) GetConfig(v ...interface{}) string {
	if len(v) == 1 {
		data, err := g.Cfg().Get(t.ctx, gconv.String(v[0]))
		if err != nil {
			g.Log().Error(t.ctx, err.Error())
			return ""
		}
		return data.String()
	}
	return ""
}

// 获取 env 变量
func (t *ToolUtil) Getenv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// 获取根目录
func (t *ToolUtil) GetRootPath(isProd bool, appName string, isDesktop bool) string {
	if isDesktop {

		if !isProd {
			return "./"
		}
		// 获取适合当前操作系统的基础存储路径
		var basePath string
		switch runtime.GOOS {
		case "windows":
			appData := os.Getenv("APPDATA")
			if appData == "" {
				appData = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
			}
			basePath = filepath.Join(appData, appName)
		case "darwin":
			homeDir, _ := os.UserHomeDir()
			basePath = filepath.Join(homeDir, "Library", "Application Support", appName)
		default: // linux 和其他类 Unix 系统
			basePath = "/var/lib/" + appName
			// 如果不是 root 用户，使用用户目录
			if os.Getuid() != 0 {
				homeDir, _ := os.UserHomeDir()
				basePath = filepath.Join(homeDir, "."+appName)
			}
		}

		if err := os.MkdirAll(basePath, 0755); err != nil {
			g.Log().Error(t.ctx, err.Error())
			panic(err)
		}

		return basePath
	} else {
		return "./"
	}
}

// 获取数据库路径
func (t *ToolUtil) GetDataBasePath(dbFileName string, isProd bool, appName string, isDesktop bool, defaultPath string) string {
	if isProd && isDesktop {
		rootPath := t.GetRootPath(isProd, appName, isDesktop)
		path := filepath.Join(rootPath, "data", "database")

		// 创建目录，失败则 fallback
		if err := os.MkdirAll(path, 0755); err != nil {
			g.Log().Error(t.ctx, err.Error())
			panic(err)
		}

		return path + "/" + dbFileName
	}

	if defaultPath != "" {
		return defaultPath
	}
	return "./data/database/" + dbFileName

}

// 获取上传文件路径
func (t *ToolUtil) GetUploadPath(isProd bool, appName string, isDesktop bool, defaultPath string) string {

	rootPath := t.GetRootPath(isProd, appName, isDesktop)
	path := filepath.Join(rootPath, coreconfig.Config.Core.File.UploadPath)
	// 创建目录，失败则 fallback
	if err := os.MkdirAll(path, 0755); err != nil {
		g.Log().Error(t.ctx, err.Error())
		panic(err)
	}

	return path

}

// 获取日志路径
func (t *ToolUtil) GetLoggerPath(isProd bool, appName string, isDesktop bool, defaultPath string) string {

	if isProd && isDesktop {
		rootPath := t.GetRootPath(isProd, appName, isDesktop)
		path := filepath.Join(rootPath, "data", "logs")

		// 创建目录，失败则 fallback
		if err := os.MkdirAll(path, 0755); err != nil {
			g.Log().Error(t.ctx, err.Error())
			panic(err)
		}

		return path
	}

	if defaultPath != "" {
		return defaultPath
	}
	return coreconfig.Config.Core.RunLogger.Path

}

func (t *ToolUtil) GetSqlLoggerPath(isProd bool, appName string, isDesktop bool, defaultPath string) string {

	if isProd && isDesktop {
		rootPath := t.GetRootPath(isProd, appName, isDesktop)
		path := filepath.Join(rootPath, "data", "logs", "sql")
		if err := os.MkdirAll(path, 0755); err != nil {
			g.Log().Error(t.ctx, err.Error())
			panic(err)
		}
		return path
	}

	if defaultPath != "" {
		return defaultPath
	}
	return coreconfig.Config.Core.SQLLogger.Path

}

// 带吞吐量，响应时间参数的运行日志
func (t *ToolUtil) GetRunLoggerPath(isProd bool, appName string, isDesktop bool, defaultPath string) string {

	if isProd && isDesktop {
		rootPath := t.GetRootPath(isProd, appName, isDesktop)
		path := filepath.Join(rootPath, "data", "logs", "run")
		if err := os.MkdirAll(path, 0755); err != nil {
			g.Log().Error(t.ctx, err.Error())
			panic(err)
		}
		return path
	}

	if defaultPath != "" {
		return defaultPath
	}
	return coreconfig.Config.Core.RunLogger.Path

}

// 日式打印运行时间
func (t *ToolUtil) StdOutLog(ctx context.Context, startTime time.Time, memStatsStart runtime.MemStats) {
	var (
		r           = g.RequestFromCtx(ctx)
		ctxId       = gctx.CtxId(r.GetCtx()) //获取当前请求的ctxid
		elapsedTime = time.Since(startTime)  // 请求处理时间
		outLogger_  *defineType.OutputsForLogger
		memStatsEnd runtime.MemStats // 记录结束内存状态
		logPath     string
		prefix      = ""
		inside      = "{Y-m-d}"
		suffix      = ""
	)

	defaultPath := coreconfig.Config.Core.RunLogger.Path
	logPath = t.GetRunLoggerPath(envconfig.IsProd, envconfig.AppName, envconfig.IsDesktop, defaultPath)
	runLogger := &defineType.RunLogger{
		Path:       logPath,
		File:       coreconfig.Config.Core.RunLogger.File,
		RotateSize: coreconfig.Config.Core.RunLogger.RotateSize,
		Stdout:     coreconfig.Config.Core.RunLogger.Stdout,
	}

	matches, err := gregex.MatchString(`^(.*)\{(.+)\}(.*)\.log$`, runLogger.File)
	if err != nil {
		return
	}
	if len(matches) == 4 {
		prefix = matches[1]
		inside = matches[2]
		suffix = matches[3]
	}

	// 根据处理时间计算吞吐率
	throughput := 1.0 / elapsedTime.Seconds() //（秒）
	runtime.ReadMemStats(&memStatsEnd)
	// 计算内存消耗
	memUsed := memStatsEnd.Alloc - memStatsStart.Alloc

	outLogger_ = &defineType.OutputsForLogger{
		Time:       time.Now(),
		Host:       r.Host,
		RequestURI: r.RequestURI,
		Params:     gjson.MustEncodeString(r.GetMap()),
		RunTime:    float64(elapsedTime.Nanoseconds()) / 1e9,
		Throughput: throughput,
		MemUsed:    memUsed,
		Prefix:     prefix,
		Suffix:     suffix,
		FileRule:   inside,
		RotateSize: runLogger.RotateSize,
		Stdout:     runLogger.Stdout,
		Path:       runLogger.Path,
	}

	fname := outLogger_.Prefix + gtime.Now().Format(outLogger_.FileRule) + outLogger_.Suffix
	fileName := fmt.Sprintf("%s.log", fname)
	tempFile := fmt.Sprintf("%v/%v", outLogger_.Path, fileName)

	throughputStringFixed := decimal.NewFromFloat(outLogger_.Throughput).StringFixed(2)

	logSlice := g.SliceStr{
		fmt.Sprintf("[ %s ] %s OPTIONS %s\n", outLogger_.Time, outLogger_.Host, outLogger_.RequestURI),
		fmt.Sprintf("[ 运行时间：%vs ] [TraceId：%v ] [ 吞吐率：%vreq/s ] [ 内存消耗：%v ]\n", outLogger_.RunTime, ctxId, throughputStringFixed, humanize.Bytes(outLogger_.MemUsed)),
		fmt.Sprintf("[ info ] [ PARAM ] [ %v ]\n", t.StrTranLine(outLogger_.Params)),
	}

	//超过容量就切割
	byteSize, _ := humanize.ParseBytes(outLogger_.RotateSize)
	if gfile.Size(tempFile) > int64(byteSize) {
		endTime := gtime.Now().Format("H-i-s")
		dstPath := tempFile + "." + endTime
		gfile.Rename(tempFile, dstPath)
	}

	//写入到日志
	for _, log := range logSlice {
		gfile.PutContentsAppend(tempFile, log)
	}
	gfile.PutContentsAppend(tempFile, "----------------------------------------\n")

	//打印到控制台
	if outLogger_.Stdout {
		for _, log := range logSlice {
			g.Log().Info(ctx, t.StrTranLine(log))
		}
		g.Log().Info(ctx, "----------------------------------------")
	}

}

// 多行文本转一行
func (t *ToolUtil) StrTranLine(jsonData string) string {
	// 移除换行符和制表符
	oneLine := strings.ReplaceAll(jsonData, "\n", "")
	oneLine = strings.ReplaceAll(oneLine, "\t", "")
	oneLine = strings.ReplaceAll(oneLine, "  ", "") // 可根据需要移除多余空格
	return oneLine
}
