package dzhCore

import (
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	_ "github.com/gzdzh/dzhcore/contrib/drivers/mysql"
)

var (
	ctx = gctx.GetInitCtx()
)

func NewInit() {
	glog.Debug(ctx, "------------ dzhCore NewInit ")
}
