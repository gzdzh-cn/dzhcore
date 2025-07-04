package sqlite

import (
	"os"
	"path/filepath"

	// _ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	"github.com/gogf/gf/v2/encoding/gurl"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"

	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gzdzh-cn/dzhcore/coreconfig"
	"github.com/gzdzh-cn/dzhcore/coredb"
	"github.com/gzdzh-cn/dzhcore/utility/util"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	ctx = gctx.GetInitCtx()
)

type DriverSqlite struct {
	source string
}

func NewSqlite() coredb.Driver {
	return &DriverSqlite{}
}

func (d *DriverSqlite) GetRootPath(configNode *gdb.ConfigNode) string {
	var (
		source string
	)
	if configNode.Link != "" {
		source = configNode.Link
	} else {
		source = configNode.Name
	}

	if coreconfig.Config.Core.IsDesktop && coreconfig.Config.Core.IsProd {
		dbFileName := filepath.Base(source)
		d.source = util.NewToolUtil().GetDataBasePath(dbFileName, coreconfig.Config.Core.IsProd, coreconfig.Config.Core.AppName, coreconfig.Config.Core.IsDesktop, source)
		source = d.source
	} else {

		// ./database/dzhgo_go.sqlite 截取最后一个斜杠前的
		dir := filepath.Dir(source)   // 目录部分
		base := filepath.Base(source) // 文件名部分
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "./database/" + base
		}
		// It searches the source file to locate its absolute path..
		if absolutePath, _ := gfile.Search(source); absolutePath != "" {
			source = absolutePath
			g.Log().Debug(ctx, "absolutePath", absolutePath)
		}
	}

	return source

}

func (d *DriverSqlite) GetConn(configNode *gdb.ConfigNode) (db *gorm.DB, err error) {

	sourcePath := d.GetRootPath(configNode)
	if configNode.Extra != "" {
		var (
			options  string
			extraMap map[string]interface{}
		)
		if extraMap, err = gstr.Parse(configNode.Extra); err != nil {
			return nil, err
		}
		for k, v := range extraMap {
			if options != "" {
				options += "&"
			}
			options += fmt.Sprintf(`_pragma=%s(%s)`, k, gurl.Encode(gconv.String(v)))
		}
		if len(options) > 1 {
			sourcePath += "?" + options
		}
	}

	g.Log().Debugf(ctx, "Will use %s to open DB", sourcePath)
	return gorm.Open(sqlite.Open(sourcePath), &gorm.Config{})

}

func init() {
	g.Log().Debug(ctx, "------------ sqlite init start")
	var (
		err         error
		driverObj   = NewSqlite()
		driverNames = g.SliceStr{"sqlite"}
	)
	for _, driverName := range driverNames {
		if err = coredb.Register(driverName, driverObj); err != nil {
			panic(err)
		}
	}

	g.Log().Debug(ctx, "------------ sqlite init end")
}
