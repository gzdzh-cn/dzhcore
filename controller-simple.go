package dzhcore

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
)

var (
	ControllerSimples []IControllerSimple
)

type IControllerSimple interface {
}
type ControllerSimple struct {
	Prefix string
}

// 添加ControllerSimple到ControllerSimples数组
func AddControllerSimple(c IControllerSimple) {
	ControllerSimples = append(ControllerSimples, c)
}

// 批量注册路由
func RegisterControllerSimples() {
	for _, controller := range ControllerSimples {
		RegisterControllerSimple(controller)
	}
}

// 注册不带crud的路由
func RegisterControllerSimple(c IControllerSimple) {
	var sController = &ControllerSimple{}
	gconv.Struct(c, &sController)
	g.Server().Group(
		sController.Prefix, func(group *ghttp.RouterGroup) {
			group.Middleware(MiddlewareHandlerResponse)
			group.Bind(
				c,
			)
		})
}
