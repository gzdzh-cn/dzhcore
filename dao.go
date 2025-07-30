package dzhcore

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
)

// IDao is the interface for DAO objects
type IDao interface {
	DB() gdb.DB
	Table() string
	Group() string
	// Columns() any
	Ctx(ctx context.Context) *gdb.Model
	Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error)
}

// DaoWrapper 包装器，适配 AddonsCustomerProCluesDao 到 IDao 接口
type DaoWrapper struct {
	dao interface{}
}

func (w *DaoWrapper) DB() gdb.DB {
	return w.dao.(interface{ DB() gdb.DB }).DB()
}

func (w *DaoWrapper) Table() string {
	return w.dao.(interface{ Table() string }).Table()
}

func (w *DaoWrapper) Group() string {
	return w.dao.(interface{ Group() string }).Group()
}

// func (w *DaoWrapper) Columns() any {
// 	return w.dao.(interface{ Columns() any }).Columns()
// }

func (w *DaoWrapper) Ctx(ctx context.Context) *gdb.Model {
	return w.dao.(interface {
		Ctx(context.Context) *gdb.Model
	}).Ctx(ctx)
}

func (w *DaoWrapper) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return w.dao.(interface {
		Transaction(context.Context, func(context.Context, gdb.TX) error) error
	}).Transaction(ctx, f)
}

func NewDaoWrapper(dao interface{}) *DaoWrapper {
	return &DaoWrapper{dao: dao}
}
