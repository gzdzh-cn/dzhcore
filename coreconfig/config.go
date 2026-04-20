package coreconfig

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gzdzh-cn/dzhcore/defineStruct"
	"github.com/gzdzh-cn/dzhcore/utility/env"
	"github.com/joho/godotenv"
)

var (
	ctx    = gctx.GetInitCtx()
	Config *defineStruct.Config
)

func init() {
	Config = newConfig()
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		g.Log().Debug(ctx, "未找到.env文件，使用默认环境变量")
	}
}

// NewConfig new config
func newConfig() *defineStruct.Config {
	return &defineStruct.Config{
		Server: defineStruct.ServerConfig{
			Address:           env.GetCfgWithDefault(ctx, "server.address", g.NewVar(":8200")).String(),
			OpenapiPath:       env.GetCfgWithDefault(ctx, "server.openapiPath", g.NewVar("/api.json")).String(),
			SwaggerPath:       env.GetCfgWithDefault(ctx, "server.swaggerPath", g.NewVar("/swagger")).String(),
			ServerRoot:        env.GetCfgWithDefault(ctx, "server.serverRoot", g.NewVar("resource/public")).String(),
			ClientMaxBodySize: env.GetCfgWithDefault(ctx, "server.clientMaxBodySize", g.NewVar(104857600)).Int(),
			Paths:             env.GetCfgWithDefault(ctx, "server.paths", g.NewVar([]string{"template"})).Strings(),
			DefaultFile:       env.GetCfgWithDefault(ctx, "server.defaultFile", g.NewVar("index.html")).String(),
			Delimiters:        env.GetCfgWithDefault(ctx, "server.delimiters", g.NewVar([]string{"{{", "}}"})).Strings(),
		},
		Database: defineStruct.DatabaseConfig{
			Link:      env.GetCfgWithDefault(ctx, "database.link", g.NewVar("")).String(),
			Type:      env.GetCfgWithDefault(ctx, "database.type", g.NewVar("sqlite")).String(),
			Name:      env.GetCfgWithDefault(ctx, "database.name", g.NewVar("./data/database/dzhgo_go.sqlite")).String(),
			Host:      env.GetCfgWithDefault(ctx, "database.host", g.NewVar("127.0.0.1")).String(),
			Port:      env.GetCfgWithDefault(ctx, "database.port", g.NewVar("3306")).String(),
			User:      env.GetCfgWithDefault(ctx, "database.user", g.NewVar("")).String(),
			Pass:      env.GetCfgWithDefault(ctx, "database.pass", g.NewVar("")).String(),
			Charset:   env.GetCfgWithDefault(ctx, "database.charset", g.NewVar("utf8mb4")).String(),
			Timezone:  env.GetCfgWithDefault(ctx, "database.timezone", g.NewVar("Asia/Shanghai")).String(),
			Extra:     env.GetCfgWithDefault(ctx, "database.extra", g.NewVar("")).String(),
			CreatedAt: env.GetCfgWithDefault(ctx, "database.createdAt", g.NewVar("createTime")).String(),
			UpdatedAt: env.GetCfgWithDefault(ctx, "database.updatedAt", g.NewVar("updateTime")).String(),
			DeletedAt: env.GetCfgWithDefault(ctx, "database.deletedAt", g.NewVar("deletedAt")).String(),
			Debug:     env.GetCfgWithDefault(ctx, "database.debug", g.NewVar(false)).Bool(),
		},
		Redis: defineStruct.RedisConfig{
			Enable: env.GetCfgWithDefault(ctx, "redis.enable", g.NewVar(false)).Bool(),
			CfRedis: defineStruct.RedisCore{
				Address: env.GetCfgWithDefault(ctx, "redis.cfRedis.address", g.NewVar("127.0.0.1:16379")).String(),
				DB:      env.GetCfgWithDefault(ctx, "redis.cfRedis.db", g.NewVar(0)).Int(),
				Pass:    env.GetCfgWithDefault(ctx, "redis.cfRedis.pass", g.NewVar("")).String(),
				Expire:  env.GetCfgWithDefault(ctx, "redis.cfRedis.expire", g.NewVar(12960000)).Uint(),
			},
			DBRedis: defineStruct.DBRedisConfig{
				Address: env.GetCfgWithDefault(ctx, "redis.dbRedis.address", g.NewVar("127.0.0.1:16379")).String(),
				Enable:  env.GetCfgWithDefault(ctx, "redis.dbRedis.enable", g.NewVar(false)).Bool(),
				Expire:  env.GetCfgWithDefault(ctx, "redis.dbRedis.expire", g.NewVar(60000)).Uint(),
				DB:      env.GetCfgWithDefault(ctx, "redis.dbRedis.db", g.NewVar(9)).Int(),
			},
		},
		Elasticsearch: defineStruct.ElasticsearchConfig{
			Enable:   env.GetCfgWithDefault(ctx, "elasticsearch.enable", g.NewVar(true)).Bool(),
			Host:     env.GetCfgWithDefault(ctx, "elasticsearch.host", g.NewVar("http://127.0.0.1:9200")).String(),
			Username: env.GetCfgWithDefault(ctx, "elasticsearch.username", g.NewVar("")).String(),
			Password: env.GetCfgWithDefault(ctx, "elasticsearch.password", g.NewVar("")).String(),
		},
		Core: defineStruct.CoreConfig{
			AppName:     env.GetCfgWithDefault(ctx, "core.appName", g.NewVar("dzhgo")).String(),
			IsDesktop:   env.GetCfgWithDefault(ctx, "core.isDesktop", g.NewVar(false)).Bool(),
			IsProd:      env.GetCfgWithDefault(ctx, "core.isProd", g.NewVar(false)).Bool(),
			AutoMigrate: env.GetCfgWithDefault(ctx, "core.autoMigrate", g.NewVar(true)).Bool(),
			Eps:         env.GetCfgWithDefault(ctx, "core.eps", g.NewVar(true)).Bool(),
			Notice: defineStruct.NoticeConfig{
				Enable: env.GetCfgWithDefault(ctx, "core.notice.enable", g.NewVar(true)).Bool(),
				Send:   env.GetCfgWithDefault(ctx, "core.notice.send", g.NewVar(5)).Int(),
			},
			SQLLogger: defineStruct.LoggerConfig{
				Path:   env.GetCfgWithDefault(ctx, "core.sqlLogger.path", g.NewVar("./data/logs/sql")).String(),
				File:   env.GetCfgWithDefault(ctx, "core.sqlLogger.file", g.NewVar("sql-{Y-m-d}.log")).String(),
				Level:  env.GetCfgWithDefault(ctx, "core.sqlLogger.level", g.NewVar("all")).String(),
				Stdout: env.GetCfgWithDefault(ctx, "core.sqlLogger.stdout", g.NewVar(false)).Bool(),
			},
			GFLogger: defineStruct.LoggerConfig{
				Path:   env.GetCfgWithDefault(ctx, "core.gfLogger.path", g.NewVar("./data/logs/")).String(),
				File:   env.GetCfgWithDefault(ctx, "core.gfLogger.file", g.NewVar("{Y-m-d}.log")).String(),
				Level:  env.GetCfgWithDefault(ctx, "core.gfLogger.level", g.NewVar("debug")).String(),
				Stdout: env.GetCfgWithDefault(ctx, "core.gfLogger.stdout", g.NewVar(true)).Bool(),
				Flags:  env.GetCfgWithDefault(ctx, "core.gfLogger.flags", g.NewVar(44)).Int(),
			},
			RunLogger: defineStruct.RunLogger{
				LoggerConfig: defineStruct.LoggerConfig{
					Path:   env.GetCfgWithDefault(ctx, "core.runLogger.path", g.NewVar("./data/logs/run/")).String(),
					File:   env.GetCfgWithDefault(ctx, "core.runLogger.file", g.NewVar("run-{Y-m-d}.log")).String(),
					Level:  env.GetCfgWithDefault(ctx, "core.runLogger.level", g.NewVar("debug")).String(),
					Stdout: env.GetCfgWithDefault(ctx, "core.runLogger.stdout", g.NewVar(false)).Bool(),
					Flags:  env.GetCfgWithDefault(ctx, "core.runLogger.flags", g.NewVar(44)).Int(),
				},
				Enable:     env.GetCfgWithDefault(ctx, "core.runLogger.enable", g.NewVar(true)).Bool(),
				RotateSize: env.GetCfgWithDefault(ctx, "core.runLogger.rotateSize", g.NewVar("3MB")).String(),
			},
			File: defineStruct.FileConfig{
				Mode:       env.GetCfgWithDefault(ctx, "core.file.mode", g.NewVar("local")).String(),
				FilePrefix: env.GetCfgWithDefault(ctx, "core.file.filePrefix", g.NewVar("")).String(),
				UploadPath: env.GetCfgWithDefault(ctx, "core.file.uploadPath", g.NewVar("resource/uploads")).String(),
				Oss: defineStruct.OssConfig{
					Endpoint:        env.GetCfgWithDefault(ctx, "core.file.oss.endpoint", g.NewVar("")).String(),
					AccessKeyID:     env.GetCfgWithDefault(ctx, "core.file.oss.accessKeyID", g.NewVar("")).String(),
					SecretAccessKey: env.GetCfgWithDefault(ctx, "core.file.oss.secretAccessKey", g.NewVar("")).String(),
					BucketName:      env.GetCfgWithDefault(ctx, "core.file.oss.bucketName", g.NewVar("")).String(),
					UseSSL:          env.GetCfgWithDefault(ctx, "core.file.oss.useSSL", g.NewVar(false)).Bool(),
					Location:        env.GetCfgWithDefault(ctx, "core.file.oss.location", g.NewVar("")).String(),
				},
			},
		},
		Modules: defineStruct.ModulesConfig{
			Base: defineStruct.BaseModuleConfig{
				JWT: defineStruct.JWTConfig{
					SSO:    env.GetCfgWithDefault(ctx, "modules.base.jwt.sso", g.NewVar(false)).Bool(),
					Secret: env.GetCfgWithDefault(ctx, "modules.base.jwt.secret", g.NewVar("88888888")).String(),
					Token: defineStruct.TokenConfig{
						Expire:        env.GetCfgWithDefault(ctx, "modules.base.jwt.token.expire", g.NewVar(604800)).Uint(),
						RefreshExpire: env.GetCfgWithDefault(ctx, "modules.base.jwt.token.refreshExpire", g.NewVar(1296000)).Uint(),
					},
				},
				Middleware: defineStruct.MiddlewareConfig{
					CORS: env.GetCfgWithDefault(ctx, "modules.base.middleware.cors", g.NewVar(false)).Bool(),
					Authority: defineStruct.AuthorityConf{
						Enable: env.GetCfgWithDefault(ctx, "modules.base.middleware.authority.enable", g.NewVar(true)).Bool(),
					},
					Log: defineStruct.LogConf{
						Enable:     env.GetCfgWithDefault(ctx, "modules.base.middleware.log.enable", g.NewVar(true)).Bool(),
						IgnorePath: env.GetCfgWithDefault(ctx, "modules.base.middleware.log.ignorePath", g.NewVar("/admin/base/sys/log/getKeep")).String(),
						IgnoreReg:  env.GetCfgWithDefault(ctx, "modules.base.middleware.log.ignoreReg", g.NewVar("/(page|list)$")).String(),
					},
				},
				HTTP: defineStruct.HTTPConfig{
					ProxyOpen: env.GetCfgWithDefault(ctx, "modules.base.http.proxy_open", g.NewVar(false)).Bool(),
					ProxyURL:  env.GetCfgWithDefault(ctx, "modules.base.http.proxy_url", g.NewVar("")).String(),
				},
				Img: defineStruct.ImgConfig{
					CDNUrl: env.GetCfgWithDefault(ctx, "modules.base.img.cdn_url", g.NewVar("")).String(),
				},
			},
		},
	}
}
