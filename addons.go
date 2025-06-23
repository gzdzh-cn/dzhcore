package dzhcore

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
)

var (
	Addons []IAddon
)

type IAddon interface {
	NewInit()
	GetName() string
	GetVersion() string
}

type Addon struct {
	Name    string
	Version string
}

// GetName 返回插件名称
func (a *Addon) GetName() string {
	return a.Name
}

// GetVersion 返回插件版本
func (a *Addon) GetVersion() string {
	return a.Version
}

func AddAddon(addon IAddon) {
	Addons = append(Addons, addon)
}

func InitAddons() {
	var names []string
	for _, addon := range Addons {
		names = append(names, addon.GetName()+"-"+addon.GetVersion())
	}
	g.Log().Debugf(ctx, "InitAddons,数量： %v，分别是：%v", len(Addons), gstr.Join(names, ","))

	for _, addon := range Addons {
		g.Log().Debugf(ctx, "addon: %v", addon.GetName())
		addon.NewInit()
		SetVersions(addon.GetName(), addon.GetVersion())
	}
}
