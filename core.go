package dzhcore

import (
	"context"
	"path/filepath"
	"time"

	"github.com/gzdzh-cn/dzhcore/coreconfig"
	"github.com/gzdzh-cn/dzhcore/envconfig"
	"github.com/gzdzh-cn/dzhcore/log"
	"github.com/gzdzh-cn/dzhcore/utility/env"
	"github.com/gzdzh-cn/dzhcore/utility/util"
	"gorm.io/gorm"

	"github.com/bwmarrin/snowflake"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gbuild"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/guid"

	_ "github.com/gzdzh-cn/dzhcore/contrib/drivers/mysql"
	_ "github.com/gzdzh-cn/dzhcore/contrib/drivers/pgsql"
	_ "github.com/gzdzh-cn/dzhcore/contrib/drivers/sqlite"
	"github.com/gzdzh-cn/dzhcore/contrib/files/local"
	_ "github.com/gzdzh-cn/dzhcore/contrib/files/local"
	_ "github.com/gzdzh-cn/dzhcore/contrib/files/oss"
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
	RedisConfig    = &gredis.Config{}
	DbExpire       int64
)

func init() {

}

func NewInit() {
	g.Log().Debug(ctx, "------------ dzhcore NewInit start")

	coreconfig.LoadEnv()
	getConfig()
	local.NewInit()
	setDataBase()
	setLogger()
	SetVersions("dzhcore", Version)

	NodeSnowflake = CreateSnowflake(ctx) //雪花节点创建

	if IsRedisMode {

		redisVar, err := g.Cfg().Get(ctx, "redis.core")
		if err != nil {
			g.Log().Errorf(ctx, "初始化缓存失败,请检查配置文件:%v", err)
			panic(err)
		}

		if !redisVar.IsEmpty() {

			err = redisVar.Struct(RedisConfig)
			if err != nil {
				g.Log().Errorf(ctx, "初始化缓存失败,请检查配置文件:%v", err)
				return
			}
			redis, err := gredis.New(RedisConfig)
			if err != nil {
				g.Log().Errorf(ctx, "初始化缓存失败,请检查配置文件:%v", err)
				panic(err)
			}
			g.Log().Debug(ctx, "redis开启成功")
			CacheManager.SetAdapter(gcache.NewAdapterRedis(redis))
		}

		//db 查询使用指定缓存分组
		if DbRedisEnable {

			dbNum := env.GetCfgWithDefault(ctx, "redis.dbRedis.db", g.NewVar(9)).Int()
			RedisConfig.Db = dbNum
			redis, err := gredis.New(RedisConfig)
			if err != nil {
				g.Log().Errorf(ctx, "初始化缓存失败,请检查配置文件:%v", err)
				panic(err)
			}
			redisCache := gcache.NewAdapterRedis(redis)
			DbCacheManager.SetAdapter(redisCache)

		}
	}

	// 创建全部表
	InitModels()
	// 注册路由
	RegisterControllers()
	RegisterControllerSimples()

	if envconfig.IsProd {
		g.Log().Info(ctx, "生产环境")
	} else {
		g.Log().Info(ctx, "开发环境")
	}
	g.Log().Debugf(ctx, "IsProd:%v, AppName:%v, IsDesktop:%v", envconfig.IsProd, envconfig.AppName, envconfig.IsDesktop)
	g.Log().Debugf(ctx, "dzhcore version:%v", Version)
	g.Log().Debugf(ctx, "当前运行模式:%v", RunMode)
	g.Log().Debugf(ctx, "当前实例ID:%v", ProcessFlag)
	g.Log().Debugf(ctx, "是否redis缓存模式:%v", IsRedisMode)
	g.Log().Debugf(ctx, "是否DbRedisEnable缓存模式:%v", DbRedisEnable)

	g.Log().Debug(ctx, "------------ dzhcore NewInit end")

}

// 获取配置
func getConfig() {
	envconfig.IsDesktop = coreconfig.Config.Core.IsDesktop
	envconfig.AppName = coreconfig.Config.Core.AppName
	gbuildData := gbuild.Data()
	if envconfig.IsDesktop {
		envconfig.IsProd = coreconfig.Config.Core.IsProd
	} else {
		if _, ok := gbuildData["builtTime"]; ok {
			envconfig.IsProd = true
		} else {
			envconfig.IsProd = false
		}
	}

	if _, ok := gbuildData["mode"]; ok {
		RunMode = gbuildData["mode"].(string)
	}
	if RunMode == "core-tools" {
		return
	}
	IsRedisMode = coreconfig.Config.Redis.Enable
	DbRedisEnable = coreconfig.Config.Redis.DBRedis.Enable
	DbExpire = coreconfig.Config.Redis.DBRedis.Expire * int64(time.Millisecond)

	g.Log().Debugf(ctx, "config.IsProd:%v, config.IsDesktop:%v, config.AppName:%v", envconfig.IsProd, envconfig.IsDesktop, envconfig.AppName)
}

// 数据库配置
func setDataBase() {
	setDbConfig()
	setSqlLogger()
}

// database.default 配置
func setDbConfig() {
	// 读取 database.default 配置
	dbConfVar, err := g.Cfg().Get(ctx, "database.default")
	if err != nil {
		g.Log().Error(ctx, "读取数据库配置失败", err)
		return
	}
	if dbConfVar.IsEmpty() {
		g.Log().Error(ctx, "未找到数据库配置 database.default")
		return
	}
	var dbNode *gdb.ConfigNode
	dbConfVar.Struct(&dbNode)

	// sqlite 只需要 type、name、extra、createdAt、updatedAt、deletedAt、debug
	if dbNode.Type == "sqlite" {
		dbNode.Host = ""
		dbNode.Port = ""
		dbNode.User = ""
		dbNode.Pass = ""
		dbNode.Charset = ""
		dbNode.Timezone = ""
	}

	if envconfig.IsDesktop && envconfig.IsProd {
		var (
			source string
		)
		if dbNode.Link != "" {
			source = dbNode.Link
			dbNode.Link = ""
		} else {
			source = dbNode.Name
		}
		dbFileName := filepath.Base(source)
		dbNode.Name = util.NewToolUtil().GetDataBasePath(dbFileName, envconfig.IsProd, envconfig.AppName, envconfig.IsDesktop, source)

	}
	g.Log().Debugf(ctx, "sqlite sourcePath:%v", dbNode.Name)
	gdb.SetConfig(gdb.Config{
		"default": gdb.ConfigGroup{
			*dbNode,
		},
	})
}

// 设置sql日志
func setSqlLogger() {
	defaultPath := env.GetCfgWithDefault(ctx, "core.sqlLogger.path", g.NewVar("path")).String()
	logPath := util.NewToolUtil().GetSqlLoggerPath(envconfig.IsProd, envconfig.AppName, envconfig.IsDesktop, defaultPath)
	configMap := g.Map{
		"path":     logPath,
		"file":     coreconfig.Config.Core.SQLLogger.File,
		"level":    coreconfig.Config.Core.SQLLogger.Level,
		"stdout":   coreconfig.Config.Core.SQLLogger.Stdout,
		"flags":    coreconfig.Config.Core.SQLLogger.Flags,
		"stStatus": coreconfig.Config.Core.SQLLogger.StStatus,
		"stSkip":   coreconfig.Config.Core.SQLLogger.StSkip,
	}
	dbLogger := glog.New()
	dbLogger.SetConfigWithMap(configMap)
	g.DB().SetLogger(dbLogger)
}

// 自定义日志
func setLogger() {
	if log.Logger == nil {
		defaultPath := coreconfig.Config.Core.GFLogger.Path
		logPath := util.NewToolUtil().GetLoggerPath(envconfig.IsProd, envconfig.AppName, envconfig.IsDesktop, defaultPath)
		envconfig.ConfigMap = g.Map{
			"path":     logPath,
			"file":     coreconfig.Config.Core.GFLogger.File,
			"level":    coreconfig.Config.Core.GFLogger.Level,
			"stdout":   coreconfig.Config.Core.GFLogger.Stdout,
			"flags":    coreconfig.Config.Core.GFLogger.Flags,
			"stStatus": coreconfig.Config.Core.GFLogger.StStatus,
			"stSkip":   coreconfig.Config.Core.GFLogger.StSkip,
		}
		log.Logger = log.NewLogger(envconfig.ConfigMap) // 初始化 RunLogger 变量
		log.SetLogger(envconfig.IsProd, envconfig.AppName, envconfig.IsDesktop, defaultPath, envconfig.ConfigMap)
	}
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
