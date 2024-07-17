package dzhcore

import (
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	_ "github.com/gzdzh-cn/dzhcore/contrib/drivers/mysql"
)

var (
	ctx = gctx.GetInitCtx()
)

func NewInit() {
	SetVersions("dzhcore", Version)
	glog.Debug(ctx, "------------ dzhcore NewInit ")
	glog.Debugf(ctx, "------------ dzhcore version: %v ", Version)

}
