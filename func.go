package dzhcore

import (
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
)

type coreFunc interface {
	// Func handler
	Func(ctx g.Ctx, param string) (err error)
	// IsSingleton 是否单例,当为true时，只能有一个任务在执行,在注意函数为计划任务时使用
	IsSingleton() bool
	// IsAllWorker 是否所有worker都执行
	IsAllWorker() bool
}

// FuncMap 函数列表
var FuncMap = make(map[string]coreFunc)

// RegisterFunc 注册函数
func RegisterFunc(name string, f coreFunc) {
	FuncMap[name] = f
}

// GetFunc 获取函数
func GetFunc(name string) coreFunc {
	return FuncMap[name]
}

// RunFunc 运行函数
func RunFunc(ctx g.Ctx, funcstring string) (err error) {
	funcName := gstr.SubStr(funcstring, 0, gstr.Pos(funcstring, "("))
	funcParam := gstr.SubStr(funcstring, gstr.Pos(funcstring, "(")+1, gstr.Pos(funcstring, ")")-gstr.Pos(funcstring, "(")-1)
	if _, ok := FuncMap[funcName]; !ok {
		err = gerror.New("函数不存在:" + funcName)
		g.Log().Error(ctx, err.Error())
		return
	}
	if !FuncMap[funcName].IsAllWorker() {
		// 检查当前是否为主进程, 如果不是主进程, 则不执行
		if ProcessFlag != CacheManager.MustGetOrSet(ctx, "core:masterflag", ProcessFlag, 60*time.Second).String() {
			g.Log().Debug(ctx, "当前进程不是主进程, 不执行单例函数", funcName)
			return
		}
	}
	err = FuncMap[funcName].Func(ctx, funcParam)
	return
}

// ClusterRunFunc 集群运行函数,如果是单机模式, 则直接运行函数
func ClusterRunFunc(ctx g.Ctx, funcstring string) (err error) {
	if IsRedisMode {
		conn, err := g.Redis("core").Conn(ctx)
		if err != nil {
			return err
		}
		defer conn.Close(ctx)
		_, err = conn.Do(ctx, "publish", "core:func", funcstring)
		return err
	} else {
		return RunFunc(ctx, funcstring)
	}
}

// ListenFunc 监听函数
func ListenFunc(ctx g.Ctx) {
	if IsRedisMode {
		conn, err := g.Redis("core").Conn(ctx)
		if err != nil {
			panic(err)
		}
		defer conn.Close(ctx)
		_, err = conn.Do(ctx, "subscribe", "core:func")
		if err != nil {
			panic(err)
		}

		for {
			data, err := conn.Receive(ctx)
			//g.Dump("ListenFunc data", data)
			if err != nil {
				g.Log().Error(ctx, err.Error())
				time.Sleep(10 * time.Second)
				continue
			}
			if data != nil {

				dataMap := data.MapStrStr()
				if dataMap["Kind"] == "subscribe" {
					continue
				}
				if dataMap["Channel"] == "core:func" {
					g.Log().Debugf(ctx, "执行函数:%v", dataMap["Payload"])

					err := RunFunc(ctx, dataMap["Payload"])
					if err != nil {
						g.Log().Errorf(ctx, "执行函数失败:%v", err.Error())
					}
				}
			}
		}
	} else {
		panic(gerror.New("集群模式下, 请使用Redis作为缓存"))
	}
}
