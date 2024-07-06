package dzhCore

import "github.com/gzdzh/dzhcore/coreconfig"

var (
	Config            = coreconfig.Config            // 配置中的core节相关配置
	GetCfgWithDefault = coreconfig.GetCfgWithDefault // GetCfgWithDefault 获取配置，如果配置不存在，则使用默认值
)
