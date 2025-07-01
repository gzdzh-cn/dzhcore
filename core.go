package dzhcore

import (
	"context"
	"time"

	"github.com/gzdzh-cn/dzhcore/config"
	"github.com/gzdzh-cn/dzhcore/log"
	"github.com/gzdzh-cn/dzhcore/utility/util"

	"github.com/bwmarrin/snowflake"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gbuild"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/guid"
	"gorm.io/gorm"
)

var (
	GormDBS        = make(map[string]*gorm.DB) // 定义全局gorm.DB对象集合 仅供内部使用
	CacheEPS       = gcache.New()              // 定义全局缓存对象	供EPS使用
	CacheManager   = gcache.New()              // 定义全局缓存对象	供其他业务使用
	ProcessFlag    = guid.S()                  // 定义全局进程标识
	RunMode        = "dev"                     // 定义全局运行模式
	IsRedisMode    = false                     // 定义全局是否为redis模式
	I18n           = gi18n.New()               // 定义全局国际化对象
	Versions       = g.Map{}                   // 版本列表
	NodeSnowflake  *snowflake.Node             // 雪花
	DbCacheManager = gcache.New()
	DbRedisEnable  = false // 开启db 查询结果使用 redis 缓存
	redisConfig    = &gredis.Config{}
	DbExpire       int64
	IsProd         = config.IsProd
	AppName        = config.AppName
	IsDesktop      = config.IsDesktop // 是否为桌面端
	ConfigMap      = config.ConfigMap
	RunLogger      = log.RunLogger // 日志记录器

)

func init() {
	IsDesktop = GetCfgWithDefault(ctx, "core.isDesktop", g.NewVar(false)).Bool()
	AppName = GetCfgWithDefault(ctx, "core.appName", g.NewVar("dzhgo")).String()
	gbuildData := gbuild.Data()
	if !IsDesktop {
		if _, ok := gbuildData["builtTime"]; ok {
			IsProd = true
		} else {
			IsProd = false
		}
	} else {
		IsProd = GetCfgWithDefault(ctx, "core.isProd", g.NewVar(false)).Bool()
	}

	if RunLogger == nil {
		defaultPath := GetCfgWithDefault(ctx, "core.gfLogger.path", g.NewVar("path")).String()
		logPath := util.GetLoggerPath(IsProd, AppName, IsDesktop, defaultPath)
		ConfigMap = g.Map{
			"path":     logPath,
			"level":    GetCfgWithDefault(ctx, "core.gfLogger.level", g.NewVar("debug")).String(),
			"stdout":   GetCfgWithDefault(ctx, "core.gfLogger.stdout", g.NewVar(true)).Bool(),
			"flags":    GetCfgWithDefault(ctx, "core.gfLogger.flags", g.NewVar(44)).Int(),
			"stStatus": GetCfgWithDefault(ctx, "core.gfLogger.stStatus", g.NewVar(1)).Int(),
			"stSkip":   GetCfgWithDefault(ctx, "core.gfLogger.stSkip", g.NewVar(1)).Int(),
		}
		RunLogger = log.NewRunLogger(ConfigMap) // 初始化 RunLogger 变量
	}

	SetVersions("dzhcore", Version)
	NodeSnowflake = CreateSnowflake(ctx) //雪花节点创建
	buildData := gbuild.Data()
	if _, ok := buildData["mode"]; ok {
		RunMode = buildData["mode"].(string)
	}
	if RunMode == "core-tools" {
		return
	}
	IsRedisMode = GetCfgWithDefault(ctx, "redis.enable", g.NewVar(false)).Bool()
	if IsRedisMode {

		redisVar, err := g.Cfg().Get(ctx, "redis.core")
		if err != nil {
			g.Log().Error(ctx, "初始化缓存失败,请检查配置文件")
			panic(err)
		}
		if !redisVar.IsEmpty() {
			err = redisVar.Struct(redisConfig)
			if err != nil {
				g.Log().Error(ctx, "初始化缓存失败,请检查配置文件")
				return
			}
			redis, err := gredis.New(redisConfig)
			if err != nil {
				g.Log().Error(ctx, "初始化缓存失败,请检查配置文件")
				panic(err)
			}
			CacheManager.SetAdapter(gcache.NewAdapterRedis(redis))
		}

		//db 查询使用指定缓存分组
		DbRedisEnable = GetCfgWithDefault(ctx, "redis.dbRedis.enable", g.NewVar(false)).Bool()
		DbExpire = GetCfgWithDefault(ctx, "redis.dbRedis.expire", g.NewVar(60000)).Int64() * int64(time.Millisecond)
		if DbRedisEnable {
			dbRedisVar, err := g.Cfg().Get(ctx, "redis.core")
			if err != nil {
				g.Log().Error(ctx, "初始化缓存失败,请检查配置文件")
				panic(err)
			}
			if !dbRedisVar.IsEmpty() {
				dbNum := GetCfgWithDefault(ctx, "redis.dbRedis.db", g.NewVar(9)).Int()
				redisConfig.Db = dbNum
				redis, err := gredis.New(redisConfig)
				if err != nil {
					g.Log().Error(ctx, "初始化缓存失败,请检查配置文件")
					panic(err)
				}
				redisCache := gcache.NewAdapterRedis(redis)
				DbCacheManager.SetAdapter(redisCache)
			}
		}
	}

}

func NewInit() {

	g.Log().Debug(ctx, "------------ dzhcore NewInit start")
	g.Log().Debugf(ctx, "IsProd:%v, AppName:%v, IsDesktop:%v", IsProd, AppName, IsDesktop)

	if IsProd {
		g.Log().Info(ctx, "生产环境")
	} else {
		g.Log().Info(ctx, "开发环境")
	}
	g.Log().Debugf(ctx, "dzhcore version:%v", Version)
	g.Log().Debugf(ctx, "当前运行模式:%v", RunMode)
	g.Log().Debugf(ctx, "当前实例ID:%v", ProcessFlag)
	g.Log().Debugf(ctx, "是否redis缓存模式:%v", IsRedisMode)
	g.Log().Debugf(ctx, "是否DbRedisEnable缓存模式:%v", DbRedisEnable)

	// 创建全部表
	InitModels()

	// 注册路由
	RegisterControllers()
	RegisterControllerSimples()

	g.Log().Debug(ctx, "------------ dzhcore NewInit end")

}

// BaseRes core.OK 正常返回
type BaseRes struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Ok 返回正常结果
func Ok(data interface{}, message ...string) *BaseRes {

	var msg string
	if len(msg) == 0 {
		msg = "BaseResMessage" // 默认值
	} else {
		msg = message[0]
	}
	return &BaseRes{
		Code:    1000,
		Data:    data,
		Message: I18n.Translate(context.TODO(), msg),
	}
}

// Fail 失败返回结果
func Fail(message string) *BaseRes {
	return &BaseRes{
		Code:    1001,
		Message: message,
	}
}

// 分布式函数
// func DistributedFunc(ctx g.Ctx, f func(ctx g.Ctx) (interface{}, error)) (interface{}, error) {
// 	if ProcessFlag == ctx.Request.Header.Get("processFlag") {
// 		return f(ctx)
// 	}
// 	return nil, nil
// }

// 存储版本
func SetVersions(name string, v string) {
	Versions[name] = v
}

// 获取版本
func GetVersions(name string) interface{} {
	if name == "all" {
		return Versions
	} else {
		return Versions[name]
	}
}

// 雪花
func CreateSnowflake(ctx context.Context) *snowflake.Node {
	node, err := snowflake.NewNode(1) // 1 是节点的ID
	if err != nil {
		g.Log().Error(ctx, err.Error())
	}

	return node
}
