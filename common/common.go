package common

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
)

var (
	CacheManager = gcache.New()
	Versions     = g.Map{} // 全部版本
)
