package dzhcore

import (
	"runtime"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gzdzh-cn/dzhcore/coreconfig"
	"github.com/gzdzh-cn/dzhcore/utility/util"
)

func init() {
	var s = g.Server()
	//请求日志运行明细开启
	if coreconfig.Config.Core.RunLogger.Enable {
		s.BindMiddleware("/admin/*", RunLog) //请求日志明细
		s.BindMiddleware("/app/*", RunLog)   //请求日志明细
	}

}

// 请求日志运行明细开启
func RunLog(r *ghttp.Request) {
	var (
		startTime     = time.Now() //请求进入时间
		ctx           = r.Context()
		memStatsStart runtime.MemStats
	)
	runtime.ReadMemStats(&memStatsStart)

	r.Middleware.Next()

	//日志打印运行时间
	util.NewToolUtil().StdOutLog(ctx, startTime, memStatsStart)
}
