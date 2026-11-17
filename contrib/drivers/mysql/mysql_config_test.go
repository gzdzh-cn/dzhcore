package mysql

import (
	"strings"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
)

func TestConnectionConfigPreservesParameters(t *testing.T) {
	for _, link := range []bool{false, true} {
		config := &gdb.ConfigNode{Type: "mysql", User: "example", Pass: "secret@with/slash", Protocol: "tcp", Host: "127.0.0.1", Port: "3306", Name: "example", Charset: "utf8mb4", Timezone: "Asia/Shanghai", Extra: "collation=utf8mb4_0900_ai_ci&parseTime=true&timeout=7s"}
		if link {
			config.Link = "example:secret@with/slash@tcp(127.0.0.1:3306)/old?charset=utf8mb4&loc=Asia%2FShanghai&readTimeout=3s"
		}
		parsed, err := connectionConfig(config)
		if err != nil {
			t.Fatal(err)
		}
		if parsed.Passwd != config.Pass || parsed.DBName != "example" || parsed.Collation != "utf8mb4_0900_ai_ci" || !parsed.ParseTime || parsed.Timeout.String() != "7s" || parsed.Loc.String() != "Asia/Shanghai" {
			t.Fatal("connection settings were not preserved")
		}
		if link && parsed.ReadTimeout.String() != "3s" {
			t.Fatal("link parameters lost")
		}
		if !strings.Contains(mysqlTableOptions(config, parsed), "COLLATE=utf8mb4_0900_ai_ci") {
			t.Fatal("new table default is missing")
		}
		for _, kind := range []string{"mariadb", "tidb"} {
			config.Type = kind
			if mysqlTableOptions(config, parsed) != "" {
				t.Fatal("changed other driver defaults")
			}
		}
	}
}

func TestConnectionConfigDoesNotExposeInvalidDSN(t *testing.T) {
	_, err := connectionConfig(&gdb.ConfigNode{Link: "private-password/invalid?collation=%zz"})
	if err == nil || strings.Contains(err.Error(), "private-password") {
		t.Fatal("invalid DSN error must be sanitized")
	}
}

func TestLegacyTableDefaultsRemainUnchanged(t *testing.T) {
	config := &gdb.ConfigNode{Type: "mysql", Protocol: "tcp", Host: "localhost", Port: "3306", Charset: "utf8mb4"}
	parsed, err := connectionConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	if mysqlTableOptions(config, parsed) != "" {
		t.Fatal("non-opted-in project defaults changed")
	}
}
