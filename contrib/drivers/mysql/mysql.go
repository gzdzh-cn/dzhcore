// Package mysql 扩展了 GoFrame 的 mysql 包,集成了 gorm相关功能.
package mysql

import (
	drivermysql "github.com/go-sql-driver/mysql"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gzdzh-cn/dzhcore/coredb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	// _ "github.com/gogf/gf/contrib/drivers/mysql/v2"

	"fmt"
	"net/url"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
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

// connectionConfig keeps DSN credentials and caller options intact. ConfigNode
// extras also apply to deprecated Link connections, including connection collation.
func connectionConfig(config *gdb.ConfigNode) (*drivermysql.Config, error) {
	var source string
	if config.Link != "" {
		source = config.Link
	} else {
		source = fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=%s",
			config.User, config.Pass, config.Protocol, config.Host, config.Port, config.Name, config.Charset)
		if config.Timezone != "" {
			source += "&loc=" + url.QueryEscape(config.Timezone)
		}
	}
	parsed, err := drivermysql.ParseDSN(source)
	if err != nil {
		// Parse errors can include the source or credential substrings.
		return nil, fmt.Errorf("invalid MySQL connection configuration")
	}
	if config.Name != "" {
		parsed.DBName = config.Name
	}
	if config.Extra != "" {
		extra, err := url.ParseQuery(config.Extra)
		if err != nil {
			return nil, fmt.Errorf("invalid MySQL extra parameters")
		}
		// Use the driver's own parameter parser for booleans, location, TLS, etc.
		base := parsed.FormatDSN()
		separator := "?"
		if strings.Contains(base, "?") {
			separator = "&"
		}
		parsed, err = drivermysql.ParseDSN(base + separator + extra.Encode())
		if err != nil {
			return nil, fmt.Errorf("invalid MySQL extra parameters")
		}
	}
	return parsed, nil
}

func mysqlTableOptions(config *gdb.ConfigNode, parsed *drivermysql.Config) string {
	// Only projects explicitly opting into MySQL 8's collation get these options.
	// MariaDB/TiDB and other projects retain their existing table defaults.
	if config.Type == "mysql" && parsed.Collation == "utf8mb4_0900_ai_ci" {
		return "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci"
	}
	return ""
}

func (d *DriverMysql) GetConn(config *gdb.ConfigNode) (*gorm.DB, error) {
	parsed, err := connectionConfig(config)
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(mysql.Open(parsed.FormatDSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if options := mysqlTableOptions(config, parsed); options != "" {
		// Keep options across sessions without reusing a mutable Statement between
		// independent queries/migrations on the shared connection handle.
		db = db.Set("gorm:table_options", options).Session(&gorm.Session{})
	}
	return db, nil
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
