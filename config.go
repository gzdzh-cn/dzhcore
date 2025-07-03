package dzhcore

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gzdzh-cn/dzhcore/coreconfig"
	"github.com/gzdzh-cn/dzhcore/utility/env"
)

var (
	Version = "v1.2.7"
	Config  = coreconfig.Config // 配置中的core节相关配置
	Cfg     = NewConfig()
)

type sConfig struct {
	Core *Core
}
type Core struct {
	RunLogger *RunLoggers
}
type RunLoggers struct {
	Enable bool
}

func NewConfig() *sConfig {
	return &sConfig{
		Core: &Core{
			RunLogger: &RunLoggers{
				Enable: env.GetCfgWithDefault(ctx, "core.runLogger.enable", g.NewVar(false)).Bool(),
			},
		},
	}

}
