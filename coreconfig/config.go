package coreconfig

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
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
	err := godotenv.Load()
	if err != nil {
		g.Log().Debug(ctx, "未找到.env文件，使用默认环境变量")
	}
	Config = newConfig()

}

// NewConfig new config
func newConfig() *sConfig {

	config := &sConfig{
		AutoMigrate: GetCfgWithDefault(ctx, "core.autoMigrate", g.NewVar(false)).Bool(),
		Eps:         GetCfgWithDefault(ctx, "core.eps", g.NewVar(false)).Bool(),
		File: &file{
			Mode:   GetCfgWithDefault(ctx, "core.file.mode", g.NewVar("none")).String(),
			Domain: GetCfgWithDefault(ctx, "core.file.domain", g.NewVar("http://127.0.0.1:8200")).String(),
			Oss: &oss{
				Endpoint:        util.Getenv("OSS_ENDPOINT", ""),
				AccessKeyID:     util.Getenv("OSS_ACCESS_KEY_ID", ""),
				SecretAccessKey: util.Getenv("OSS_SECRET_ACCESS_KEY", ""),
				BucketName:      util.Getenv("OSS_BUCKET_NAME", ""),
				UseSSL:          util.Getenv("OSS_USESSL", "false") == "true",
				Location:        util.Getenv("OSS_LOCATION", ""),
			},
		},
	}
	return config
}

// GetCfgWithDefault get config with default value
func GetCfgWithDefault(ctx g.Ctx, key string, defaultValue *g.Var) *g.Var {
	value, err := g.Cfg().GetWithEnv(ctx, key)
	if err != nil {
		return defaultValue
	}
	if value.IsEmpty() || value.IsNil() {
		return defaultValue
	}
	return value
}
