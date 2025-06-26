package dzhcore

import (
	"github.com/gzdzh-cn/dzhcore/coreconfig"
)

var (
	Version           = "v1.2.1"
	Config            = coreconfig.Config            // 配置中的core节相关配置
	GetCfgWithDefault = coreconfig.GetCfgWithDefault // GetCfgWithDefault 获取配置，如果配置不存在，则使用默认值

)
