package sqlite

import (
	"os"
	"path/filepath"

	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	"github.com/gogf/gf/v2/encoding/gurl"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"

	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gzdzh-cn/dzhcore/config"
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
	g.Log().Debug(ctx, "source1", source)

	if config.IsDesktop && config.IsProd {
		dbFileName := filepath.Base(source)
		d.source = util.NewToolUtil().GetDataBasePath(dbFileName, config.IsProd, config.AppName, config.IsDesktop, source)
		source = d.source
	} else {
		g.Log().Debug(ctx, "absolutePath 处理")

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

	g.Log().Debug(ctx, "source2", sourcePath)
	// if absolutePath, _ := gfile.Search("./database/dzhgo_go.sqlite"); absolutePath != "" {
	// 	sourcePath = absolutePath
	// 	g.Log().Debug(ctx, "absolutePath2", absolutePath)
	// }
	// Multiple PRAGMAs can be specified, e.g.:
	// path/to/some.db?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)
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
	// ./database/dzhgo_go.sqlite
	// /Users/lizheng/Library/Application Support/dzhgo/database/dzhgo_go.sqlite
	// /Users/lizheng/dzhgo/database/dzhgo_go.sqlite
	// g.Log().Debugf(ctx, "Will use %s to open DB", sourcePath)
	// sourcePath = "./database/1dzhgo_go.sqlite?_pragma=busy_timeout(5000)"
	// sourcePath = "/Users/lizheng/Library/Application Support/dzhgo/database/dzhgo_go.sqlite?_pragma=busy_timeout(5000)"
	g.Log().Debugf(ctx, "Will use %s to open DB", sourcePath)
	// fmt.Println("DB Dir Exists:", gfile.Exists(filepath.Dir(sourcePath)))
	// fmt.Println("DB File Exists:", gfile.Exists(sourcePath))
	// fmt.Println("DB Dir Writable:", gfile.IsWritable(filepath.Dir(sourcePath)))
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
