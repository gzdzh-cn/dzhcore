package dzhcore

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

func DDAO(d IDao, ctx context.Context) *gdb.Model {
	return d.Ctx(ctx)
}

// DBM 根据model获取 *gdb.Model
func DBM(m IModel) *gdb.Model {
	return g.DB(m.GroupName()).Model(m.TableName())
}
