package dzhcore

import (
	"context"
	"strings"
	"time"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gzdzh-cn/dzhcore/coreconfig"
)

var (
	Controllers []IController
)

type IController interface {
	Add(ctx context.Context, req *AddReq) (res *BaseRes, err error)
	Delete(ctx context.Context, req *DeleteReq) (res *BaseRes, err error)
	Update(ctx context.Context, req *UpdateReq) (res *BaseRes, err error)
	Info(ctx context.Context, req *InfoReq) (res *BaseRes, err error)
	List(ctx context.Context, req *ListReq) (res *BaseRes, err error)
	Page(ctx context.Context, req *PageReq) (res *BaseRes, err error)
}
type Controller struct {
	Prefix  string     `json:"prefix"`
	Api     g.ArrayStr `json:"api"`
	Service IService   `json:"service"`
}

type AddReq struct {
	g.Meta `path:"/add" method:"POST"`
}

type DeleteReq struct {
	g.Meta `path:"/delete" method:"POST"`
	Ids    []int `json:"ids" v:"required#请选择要删除的数据"`
}

type UpdateReq struct {
	g.Meta `path:"/update" method:"POST"`
}

func (u UpdateReq) Deadline() (deadline time.Time, ok bool) {

	return

}

func (u UpdateReq) Done() <-chan struct{} {
	//TODO implement me
	panic("implement me")
}

func (u UpdateReq) Err() error {
	//TODO implement me
	panic("implement me")
}

func (u UpdateReq) Value(key any) any {
	//TODO implement me
	panic("implement me")
}

type InfoReq struct {
	g.Meta `path:"/info" method:"GET"`
	Id     int `json:"id" v:"integer|required#请选择要查询的数据"`
}

// type InfoRes struct {
// 	*BaseRes
// 	Data interface{} `json:"data"`
// }

type ListReq struct {
	g.Meta `path:"/list" method:"POST"`
	Order  string `json:"order"`
	Sort   string `json:"sort"`
}

type PageReq struct {
	g.Meta         `path:"/page" method:"POST"`
	Page           int    `d:"1" json:"page"`     // 页码
	Size           int    `d:"15" json:"size"`    //每页条数
	Order          string `json:"order"`          // 排序字段
	Sort           string `json:"sort"`           // 排序方式 asc desc
	IsExport       bool   `json:"isExport"`       // 是否导出
	MaxExportLimit int    `json:"maxExportLimit"` // 最大导出条数,不传或者小于等于0则不限制
}

func (c *Controller) Add(ctx context.Context, req *AddReq) (res *BaseRes, err error) {
	if garray.NewStrArrayFrom(c.Api).Contains("Add") {
		var data interface{}
		err = g.DB().Ctx(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			err = c.Service.ModifyBefore(ctx, "Add", g.RequestFromCtx(ctx).GetMap())
			if err != nil {
				return err
			}
			if data, err = c.Service.ServiceAdd(ctx, req); err != nil {
				return err
			}
			rmap := g.RequestFromCtx(ctx).GetMap()
			rmap["id"] = gconv.Map(data)["id"]
			if err = c.Service.ModifyAfter(ctx, "Add", rmap); err != nil {
				return err
			}
			if err = c.Service.CacheDo(ctx, "Update", g.RequestFromCtx(ctx).GetMap()); err != nil {
				return err
			}
			return err
		})
		if err != nil {
			return Fail(err.Error()), err
		}
		return Ok(data), err
	}
	g.RequestFromCtx(ctx).Response.Status = 404
	return nil, nil
}
func (c *Controller) Delete(ctx context.Context, req *DeleteReq) (res *BaseRes, err error) {
	if garray.NewStrArrayFrom(c.Api).Contains("Delete") {
		var data interface{}
		err = g.DB().Ctx(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			if err = c.Service.ModifyBefore(ctx, "Delete", g.RequestFromCtx(ctx).GetMap()); err != nil {
				return err
			}
			if data, err = c.Service.ServiceDelete(ctx, req); err != nil {
				return err
			}
			if err = c.Service.ModifyAfter(ctx, "Delete", g.RequestFromCtx(ctx).GetMap()); err != nil {
				return err
			}
			if err = c.Service.CacheDo(ctx, "Delete", g.RequestFromCtx(ctx).GetMap()); err != nil {
				return err
			}
			return err
		})
		if err != nil {
			return Fail(err.Error()), err
		}
		return Ok(data), err
	}
	g.RequestFromCtx(ctx).Response.Status = 404
	return nil, nil
}
func (c *Controller) Update(ctx context.Context, req *UpdateReq) (res *BaseRes, err error) {
	if garray.NewStrArrayFrom(c.Api).Contains("Update") {
		var data interface{}
		err = g.DB().Ctx(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			if err = c.Service.ModifyBefore(ctx, "Update", g.RequestFromCtx(ctx).GetMap()); err != nil {
				return err
			}
			if data, err = c.Service.ServiceUpdate(ctx, req); err != nil {
				return err
			}
			if err = c.Service.ModifyAfter(ctx, "Update", g.RequestFromCtx(ctx).GetMap()); err != nil {
				return err
			}
			if err = c.Service.CacheDo(ctx, "Update", g.RequestFromCtx(ctx).GetMap()); err != nil {
				return err
			}
			return err
		})

		if err != nil {
			return Fail(err.Error()), err
		}
		return Ok(data), err

	}
	g.RequestFromCtx(ctx).Response.Status = 404
	return nil, nil
}
func (c *Controller) Info(ctx context.Context, req *InfoReq) (res *BaseRes, err error) {
	if garray.NewStrArrayFrom(c.Api).Contains("Info") {
		data, err := c.Service.ServiceInfo(ctx, req)
		return Ok(data), err
	}
	g.RequestFromCtx(ctx).Response.Status = 404
	return nil, nil
}
func (c *Controller) List(ctx context.Context, req *ListReq) (res *BaseRes, err error) {
	if garray.NewStrArrayFrom(c.Api).Contains("List") {
		data, err := c.Service.ServiceList(ctx, req)
		return Ok(data), err
	}
	g.RequestFromCtx(ctx).Response.Status = 404
	return nil, nil
}
func (c *Controller) Page(ctx context.Context, req *PageReq) (res *BaseRes, err error) {
	if garray.NewStrArrayFrom(c.Api).Contains("Page") {
		data, err := c.Service.ServicePage(ctx, req)
		return Ok(data), err
	}
	g.RequestFromCtx(ctx).Response.Status = 404
	return nil, nil
}

// 添加Controller到Controllers数组
func AddController(c IController) {
	Controllers = append(Controllers, c)
}

// 批量注册路由
func RegisterControllers() {
	for _, controller := range Controllers {
		RegisterController(controller)
	}
}

// RegisterController 注册控制器到路由
func RegisterController(c IController) {
	var ctx = context.Background()
	var sController = &Controller{}
	err := gconv.Struct(c, &sController)
	if err != nil {
		return
	}
	if coreconfig.Config.Core.Eps {
		//dao := sController.Service.GetDao()
		//columns := getModelInfo(ctx, sController.Prefix, dao)
		model := sController.Service.GetModel()
		columns := getModelInfo(ctx, sController.Prefix, model)
		ModelInfo[sController.Prefix] = columns
	}
	g.Server().Group(
		sController.Prefix, func(group *ghttp.RouterGroup) {
			group.Middleware(MiddlewareHandlerResponse)
			group.Bind(
				c,
			)
		})
}

// ColumnInfo 表字段信息
type ColumnInfo struct {
	Comment      string `json:"comment"`
	Length       string `json:"length"`
	Nullable     bool   `json:"nullable"`
	PropertyName string `json:"propertyName"`
	Type         string `json:"type"`
}

// ModelInfo 路由perfix	对应的model信息
var ModelInfo = make(map[string][]*ColumnInfo)

// getModelInfo 获取模型信息
func getModelInfo(ctx g.Ctx, prefix string, model IModel) (columns []*ColumnInfo) {

	fields, err := g.DB(model.GroupName()).TableFields(ctx, model.TableName())
	if err != nil {
		panic(err)
	}
	// RunLogger.Info(ctx, "fields", fields)
	sortedFields := garray.NewArraySize(len(fields), len(fields))
	for _, field := range fields {
		sortedFields.Set(field.Index, field)
	}
	for _, field := range sortedFields.Slice() {
		if field.(*gdb.TableField).Name == "deleted_at" {
			continue
		}
		var comment string
		if field.(*gdb.TableField).Comment != "" {
			comment = field.(*gdb.TableField).Comment
		} else {
			comment = field.(*gdb.TableField).Name
		}
		// 去除 type中的长度
		var length string
		if strings.Contains(field.(*gdb.TableField).Type, "(") {
			length = field.(*gdb.TableField).Type[strings.Index(field.(*gdb.TableField).Type, "(")+1 : strings.Index(field.(*gdb.TableField).Type, ")")]
		}
		columnType := gstr.Replace(field.(*gdb.TableField).Type, "("+length+")", "")
		column := &ColumnInfo{
			Comment:      comment,
			Length:       "",
			Nullable:     field.(*gdb.TableField).Null,
			PropertyName: field.(*gdb.TableField).Name,
			Type:         columnType,
		}
		columns = append(columns, column)
	}
	return
}
