package dzhcore

import (
	"context"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/gogf/gf/i18n/gi18n"
	"github.com/gogf/gf/util/guid"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gbuild"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
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
	DbExpire       int64
	IsProd         = false
)

func init() {
	gbuildData := gbuild.Data()
	if _, ok := gbuildData["builtTime"]; ok {
		g.Log().Warning(ctx, "生产环境")
		IsProd = true
	} else {
		g.Log().Warning(ctx, "开发环境")
		IsProd = false
	}
}

func NewInit() {

	g.Log().Debug(ctx, "------------ dzhcore NewInit start")

	var (
		ctx         = gctx.GetInitCtx()
		redisConfig = &gredis.Config{}
	)

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
				return
			}
			redis, err := gredis.New(redisConfig)
			if err != nil {
				panic(err)
			}
			CacheManager.SetAdapter(gcache.NewAdapterRedis(redis))
			// IsRedisMode = true
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
					panic(err)
				}
				redisCache := gcache.NewAdapterRedis(redis)
				DbCacheManager.SetAdapter(redisCache)
			}
		}
	}

	g.Log().Infof(ctx, "dzhcore version:%v", Version)
	g.Log().Info(ctx, "当前运行模式", RunMode)
	g.Log().Info(ctx, "当前实例ID:", ProcessFlag)
	g.Log().Info(ctx, "是否redis缓存模式:", IsRedisMode)
	g.Log().Info(ctx, "是否DbRedisEnable缓存模式:", DbRedisEnable)

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
