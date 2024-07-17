package dzhcore

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
)

type IDao interface {
	DB() gdb.DB
	Table() string
	Group() string
	Ctx(ctx context.Context) *gdb.Model
}
