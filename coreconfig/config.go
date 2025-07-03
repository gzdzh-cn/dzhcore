package coreconfig

import (
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gbuild"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gzdzh-cn/dzhcore/config"
	"github.com/gzdzh-cn/dzhcore/utility/env"
	"github.com/gzdzh-cn/dzhcore/utility/util"
	"github.com/joho/godotenv"
)

var (
	ctx    = gctx.GetInitCtx()
	Config *sConfig
)

// core config
type sConfig struct {
	AutoMigrate bool  `json:"auto_migrate,omitempty"` // 是否自动创建表
	Eps         bool  `json:"eps,omitempty"`          // 是否开启eps
	File        *file `json:"file,omitempty"`         // 文件上传配置
}

// OSS相关配置
type oss struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"accessKeyID"`
	SecretAccessKey string `json:"secretAccessKey"`
	UseSSL          bool   `json:"useSSL"`
	BucketName      string `json:"bucketName"`
	Location        string `json:"location"`
}

// 文件上传配置
type file struct {
	Mode   string `json:"mode"`   // 模式 local oss
	Domain string `json:"domain"` // 域名 http://
	Oss    *oss   `json:"oss,omitempty"`
}

func init() {
	gbuildData := gbuild.Data()
	if _, ok := gbuildData["builtTime"]; ok {
		config.IsProd = true
	}
	err := godotenv.Load()
	if err != nil {
		g.Log().Debug(ctx, "未找到.env文件，使用默认环境变量")
	}
	Config = newConfig()

}

// NewConfig new config
func newConfig() *sConfig {

	newCfg := &sConfig{
		AutoMigrate: env.GetCfgWithDefault(ctx, "core.autoMigrate", g.NewVar(false)).Bool(),
		Eps:         env.GetCfgWithDefault(ctx, "core.eps", g.NewVar(false)).Bool(),
		File: &file{
			Mode:   env.GetCfgWithDefault(ctx, "core.file.mode", g.NewVar("none")).String(),
			Domain: env.GetCfgWithDefault(ctx, "core.file.domain", g.NewVar("http://127.0.0.1:8200")).String(),
			Oss: &oss{
				Endpoint: func() string {
					if !config.IsProd {
						return util.Getenv("OSS_ENDPOINT", env.GetCfgWithDefault(ctx, "core.file.oss.endpoint", g.NewVar("")).String())
					}
					return env.GetCfgWithDefault(ctx, "core.file.oss.endpoint", g.NewVar("")).String()
				}(),
				AccessKeyID: func() string {
					if !config.IsProd {
						return util.Getenv("OSS_ACCESS_KEY_ID", env.GetCfgWithDefault(ctx, "core.file.oss.accessKeyID", g.NewVar("")).String())
					}
					return env.GetCfgWithDefault(ctx, "core.file.oss.accessKeyID", g.NewVar("")).String()
				}(),
				SecretAccessKey: func() string {
					if !config.IsProd {
						return util.Getenv("OSS_SECRET_ACCESS_KEY", env.GetCfgWithDefault(ctx, "core.file.oss.secretAccessKey", g.NewVar("")).String())
					}
					return env.GetCfgWithDefault(ctx, "core.file.oss.secretAccessKey", g.NewVar("")).String()
				}(),
				BucketName: func() string {
					if !config.IsProd {
						return util.Getenv("OSS_BUCKET_NAME", env.GetCfgWithDefault(ctx, "core.file.oss.bucketName", g.NewVar("")).String())
					}
					return env.GetCfgWithDefault(ctx, "core.file.oss.bucketName", g.NewVar("")).String()
				}(),
				UseSSL: func() bool {
					if !config.IsProd {
						// 先从环境变量获取，若未设置则用配置文件
						envVal := util.Getenv("OSS_USESSL", "")
						if envVal != "" {
							// 支持常见的布尔字符串
							lower := strings.ToLower(envVal)
							return lower == "1" || lower == "true" || lower == "yes" || lower == "on"
						}
					}
					return env.GetCfgWithDefault(ctx, "core.file.oss.useSSL", g.NewVar(false)).Bool()
				}(),
				Location: util.Getenv("OSS_LOCATION", env.GetCfgWithDefault(ctx, "core.file.oss.location", g.NewVar("")).String()),
			},
		},
	}
	return newCfg
}
