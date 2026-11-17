package mysql

import (
	"context"
	"database/sql"
	"net"
	"os"
	"strings"
	"testing"

	drivermysql "github.com/go-sql-driver/mysql"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/database/gdb"
)

type collationProbe struct {
	ID   int    `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(64)"`
}

func (collationProbe) TableName() string { return "collation_driver_probe" }

func TestMySQL8DefaultsIntegration(t *testing.T) {
	port := os.Getenv("DZHGO_MYSQL8_TEST_PORT")
	if port == "" {
		t.Skip("requires disposable test container port")
	}
	if !strings.HasPrefix(os.Getenv("DZHGO_COLLATION_TEST_CONTAINER"), "dzhgo-collation-test-") {
		t.Fatal("only disposable containers are allowed")
	}
	rootConfig := drivermysql.NewConfig()
	rootConfig.User = "root"
	rootConfig.Net = "tcp"
	rootConfig.Addr = net.JoinHostPort("127.0.0.1", port)
	root, err := sql.Open("mysql", rootConfig.FormatDSN())
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()
	if _, err = root.Exec("CREATE DATABASE collation_driver_test CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"); err != nil {
		t.Fatal(err)
	}
	defer root.Exec("DROP DATABASE collation_driver_test")
	config := &gdb.ConfigNode{Type: "mysql", Protocol: "tcp", Host: "127.0.0.1", Port: port, User: "root", Name: "collation_driver_test", Charset: "utf8mb4", Extra: "collation=utf8mb4_0900_ai_ci&parseTime=true"}
	orm, err := NewMysql().GetConn(config)
	if err != nil {
		t.Fatal(err)
	}
	pool, err := orm.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	pool.SetMaxIdleConns(0)
	for i := 0; i < 3; i++ {
		var charset, collation string
		if err := orm.Raw("SELECT @@character_set_connection,@@collation_connection").Row().Scan(&charset, &collation); err != nil {
			t.Fatal(err)
		}
		if charset != "utf8mb4" || collation != "utf8mb4_0900_ai_ci" {
			t.Fatal("GORM connection defaults differ")
		}
	}
	if err := orm.AutoMigrate(&collationProbe{}); err != nil {
		t.Fatal(err)
	}
	if err := orm.Exec("ALTER TABLE collation_driver_probe ADD COLUMN next_text varchar(32)").Error; err != nil {
		t.Fatal(err)
	}
	var count int
	if err := orm.Raw("SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='collation_driver_test' AND TABLE_NAME='collation_driver_probe' AND CHARACTER_SET_NAME IS NOT NULL AND (CHARACTER_SET_NAME<>'utf8mb4' OR COLLATION_NAME<>'utf8mb4_0900_ai_ci')").Row().Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatal("new table/column defaults differ")
	}
	gf, err := gdb.New(*config)
	if err != nil {
		t.Fatal(err)
	}
	row, err := gf.GetOne(context.Background(), "SELECT @@character_set_connection AS charset,@@collation_connection AS collation")
	if err != nil {
		t.Fatal(err)
	}
	if row["charset"].String() != "utf8mb4" || row["collation"].String() != "utf8mb4_0900_ai_ci" {
		t.Fatal("GoFrame connection defaults differ")
	}
	corePool, err := gf.GetCore().Master()
	if err != nil {
		t.Fatal(err)
	}
	defer corePool.Close()
}
