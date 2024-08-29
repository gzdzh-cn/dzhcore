package util

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gzdzh-cn/dzhcore/common"
	"github.com/gzdzh-cn/dzhcore/utility/util/logger"
)

var ctx = gctx.GetInitCtx()

// 清理orm缓存
func ClearOrmCache(ctx context.Context, key string) {

	_, err := common.CacheManager.Remove(ctx, "SelectCache:"+key)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}

}

// GetValueOrDefault
//
//	@Description: 给定一个interface，如果是nil或空，则返回false
//	@param value
//	@param defaultValue
//	@return interface{}
func GetValueOrDefault(value interface{}) bool {

	if value == nil {
		return false
	}

	// 使用类型断言判断是否为空值
	switch v := value.(type) {
	case string:
		if v == "" {
			return false
		}
	case int:
		if v == 0 {
			return false
		}
	case []interface{}:
		if len(v) == 0 {
			return false
		}
	// 可以根据需要添加更多类型的检查
	default:
		// 其他类型不为空的情况
		return true
	}

	// 如果不是空值，返回原值
	return true
}

// GetMySQLVersion
//
//	@Description: 获取mysql版本
//	@return string
func GetMySQLVersion() string {

	type result struct {
		Version string `json:"version"`
	}
	var res *result
	err := g.DB().Raw("SELECT VERSION() as version").Scan(&res)

	if err != nil {
		logger.Error(ctx, err.Error())
		return ""
	}

	return res.Version
}

// 配置信息
func GetConfig(v ...interface{}) string {

	if len(v) == 1 {
		data, err := g.Cfg().Get(ctx, gconv.String(v[0]))
		if err != nil {
			logger.Error(ctx, err.Error())
			return ""
		}
		return data.String()
	}

	return ""
}

// 获取版本
func GetVersions(name string) interface{} {
	if name == "all" {
		return common.Versions
	} else {
		return common.Versions[name]
	}
}
