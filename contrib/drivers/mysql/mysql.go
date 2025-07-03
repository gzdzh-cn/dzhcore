// Package mysql 扩展了 GoFrame 的 mysql 包,集成了 gorm相关功能.
package mysql

import (
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gzdzh-cn/dzhcore/coredb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	// _ "github.com/gogf/gf/contrib/drivers/mysql/v2"

	"fmt"
	"net/url"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gregex"
)

var (
	ctx         = gctx.GetInitCtx()
	err         error
	driverObj   = NewMysql()
	driverNames = g.SliceStr{"mysql", "mariadb", "tidb"}
)

type DriverMysql struct {
}

func NewMysql() coredb.Driver {
	return &DriverMysql{}
}

func (d *DriverMysql) GetConn(config *gdb.ConfigNode) (db *gorm.DB, err error) {
	var (
		source string
	)
	if config.Link != "" {
		// ============================================================================
		// Deprecated from v2.2.0.
		// ============================================================================
		source = config.Link
		// Custom changing the schema in runtime.
		if config.Name != "" {
			source, _ = gregex.ReplaceString(`/([\w\.\-]+)+`, "/"+config.Name, source)
		}
	} else {
		source = fmt.Sprintf(
			"%s:%s@%s(%s:%s)/%s?charset=%s",
			config.User, config.Pass, config.Protocol, config.Host, config.Port, config.Name, config.Charset,
		)
		if config.Timezone != "" {
			source = fmt.Sprintf("%s&loc=%s", source, url.QueryEscape(config.Timezone))
		}
		if config.Extra != "" {
			source = fmt.Sprintf("%s&%s", source, config.Extra)
		}
	}

	return gorm.Open(mysql.Open(source), &gorm.Config{})
}

func init() {

	g.Log().Debug(ctx, "------------ mysql init start ...")
	for _, driverName := range driverNames {
		if err = coredb.Register(driverName, driverObj); err != nil {
			panic(err)
		}
	}
	g.Log().Debug(ctx, "------------ mysql init end ...")
}
