package dzhcore

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

type IService interface {
	ServiceAdd(ctx context.Context, req *AddReq) (data any, err error)              // 新增
	ServiceDelete(ctx context.Context, req *DeleteReq) (data any, err error)        // 删除
	ServiceUpdate(ctx context.Context, req *UpdateReq) (data any, err error)        // 修改
	ServiceInfo(ctx context.Context, req *InfoReq) (data any, err error)            // 详情
	ServiceList(ctx context.Context, req *ListReq) (data any, err error)            // 列表
	ServicePage(ctx context.Context, req *PageReq) (data any, err error)            // 分页
	ModifyBefore(ctx context.Context, method string, param g.MapStrAny) (err error) // 新增|删除|修改前的操作
	ModifyAfter(ctx context.Context, method string, param g.MapStrAny) (err error)  // 新增|删除|修改后的操作
	CacheDo(ctx context.Context, method string, param g.MapStrAny) (err error)      // 处理 db 缓存
	GetModel() IModel
	GetDao() IDao
}

type Service struct {
	Dao                IDao
	Model              IModel
	ListQueryOp        *QueryOp
	PageQueryOp        *QueryOp
	InsertParam        func(ctx context.Context) g.MapStrAny // Add时插入参数
	Before             func(ctx context.Context) (err error) // CRUD前的操作
	InfoIgnoreProperty string                                // Info时忽略的字段,多个字段用逗号隔开
	UniqueKey          g.MapStrStr                           // 唯一键 key:字段名 value:错误信息
	NotNullKey         g.MapStrStr                           // 非空键 key:字段名 value:错误信息
}

// List/Add接口条件配置
type QueryOp struct {
	FieldEQ      []string                                 // 字段等于
	KeyWordField []string                                 // 模糊搜索匹配的数据库字段
	AddOrderby   g.MapStrStr                              // 添加排序
	Where        func(ctx context.Context) []g.Array      // 自定义条件
	OrWhere      func(ctx context.Context) []g.Array      // 自定义条件
	Select       string                                   // 查询字段,多个字段用逗号隔开 如: id,name  或  a.id,a.name,b.name AS bname
	As           string                                   //主表别名
	Join         []*JoinOp                                // 关联查询
	Extend       func(ctx g.Ctx, m *gdb.Model) *gdb.Model // 追加其他条件
	ModifyResult func(ctx g.Ctx, data any) any            // 修改结果
}

// 关联查询
type JoinOp struct {
	Dao       IDao
	Model     IModel   // 关联的model
	Alias     string   // 别名
	Condition string   // 关联条件
	Type      JoinType // 关联类型  LeftJoin RightJoin InnerJoin
}

// 关联类型
type JoinType string

// 新增
func (s *Service) ServiceAdd(ctx context.Context, req *AddReq) (data any, err error) {

	r := g.RequestFromCtx(ctx)
	rmap := r.GetMap()
	m := DDAO(s.Dao, ctx)
	// 非空键
	if s.NotNullKey != nil {
		for k, v := range s.NotNullKey {
			if rmap[k] == nil {
				return nil, gerror.New(v)
			}
		}
	}
	// 唯一键
	if s.UniqueKey != nil {
		for k, v := range s.UniqueKey {
			if rmap[k] != nil {

				count, err := m.Where(k, rmap[k]).Count()
				if err != nil {
					return nil, err
				}
				if count > 0 {
					err = gerror.New(v)
					return nil, err
				}
			}
		}
	}
	if s.InsertParam != nil {
		insertParams := s.InsertParam(ctx)
		if len(insertParams) > 0 {
			for k, v := range insertParams {
				rmap[k] = v
			}
		}
	}

	rmap["id"] = NodeSnowflake.Generate().String()
	_, err = m.Insert(rmap)
	if err != nil {
		return
	}

	data = g.Map{"id": rmap["id"]}

	return
}

// 删除
func (s *Service) ServiceDelete(ctx context.Context, req *DeleteReq) (data any, err error) {
	ids := g.RequestFromCtx(ctx).Get("ids").Slice()
	m := DDAO(s.Dao, ctx)
	data, err = m.WhereIn("id", ids).Delete()

	return
}

// 修改
func (s *Service) ServiceUpdate(ctx context.Context, req *UpdateReq) (data any, err error) {
	r := g.RequestFromCtx(ctx)
	rmap := r.GetMap()
	m := DDAO(s.Dao, ctx)
	if rmap["id"] == nil {
		err = gerror.New("id不能为空")
		g.Log().Error(ctx, err.Error())
		return
	}
	if s.UniqueKey != nil {
		for k, v := range s.UniqueKey {
			if rmap[k] != nil {
				count, err := m.Where(k, rmap[k]).WhereNot("id", rmap["id"]).Count()
				if err != nil {
					return nil, err
				}
				if count > 0 {
					err = gerror.New(v)
					g.Log().Error(ctx, err.Error())
					return nil, err
				}
			}
		}
	}

	_, err = m.Data(rmap).Where("id", gconv.String(rmap["id"])).Update()
	return
}

// 查询
func (s *Service) ServiceInfo(ctx context.Context, req *InfoReq) (data any, err error) {
	if s.Before != nil {
		err = s.Before(ctx)
		if err != nil {
			return
		}
	}

	m := DDAO(s.Dao, ctx)
	// 如果InfoIgnoreProperty不为空 则忽略相关字段
	if len(s.InfoIgnoreProperty) > 0 {
		m.FieldsEx(s.InfoIgnoreProperty)
	}
	data, err = m.Clone().Where("id", gconv.String(req.Id)).One()
	return
}

// 列表
func (s *Service) ServiceList(ctx context.Context, req *ListReq) (data any, err error) {
	if s.Before != nil {
		err = s.Before(ctx)
		if err != nil {
			return
		}
	}
	r := g.RequestFromCtx(ctx)

	m := DDAO(s.Dao, ctx)

	// 如果 req.Order 和 req.Sort 均不为空 则添加排序
	if !r.Get("order").IsEmpty() && !r.Get("sort").IsEmpty() {
		m = m.Order(r.Get("order").String() + " " + r.Get("sort").String())
	}
	// 如果 ListQueryOp 不为空 则使用 ListQueryOp 进行查询
	if s.ListQueryOp != nil {
		//主表别名
		if s.ListQueryOp.As != "" {
			m = m.As(s.ListQueryOp.As)
		}
		if Select := s.ListQueryOp.Select; Select != "" {
			m = m.Fields(Select)
		}
		// 如果Join不为空 则添加Join
		if len(s.ListQueryOp.Join) > 0 {
			for _, join := range s.ListQueryOp.Join {
				switch join.Type {
				case LeftJoin:
					m = m.LeftJoin(join.Model.TableName(), join.Condition).As(join.Alias)
				case RightJoin:
					m = m.RightJoin(join.Model.TableName(), join.Condition).As(join.Alias)
				case InnerJoin:
					m = m.InnerJoin(join.Model.TableName(), join.Condition).As(join.Alias)
				}
			}
		}

		// 如果fileldEQ不为空 则添加查询条件
		if len(s.ListQueryOp.FieldEQ) > 0 {
			for _, field := range s.ListQueryOp.FieldEQ {
				if !r.Get(field).IsEmpty() {
					m = m.Where(field, r.Get(field))
				}
			}
		}
		// 如果KeyWordField不为空 则添加查询条件
		if !r.Get("keyWord").IsEmpty() {
			if len(s.ListQueryOp.KeyWordField) > 0 {
				builder := m.Builder()
				for _, field := range s.ListQueryOp.KeyWordField {
					builder = builder.WhereOrLike(field, "%"+r.Get("keyWord").String()+"%")
				}
				m = m.Where(builder)
			}
		}
		if s.ListQueryOp.Where != nil {
			where := s.ListQueryOp.Where(ctx)
			if len(where) > 0 {
				for _, v := range where {
					if len(v) == 3 {
						if gconv.Bool(v[2]) {
							m = m.Where(v[0], v[1])
						}
					}
					if len(v) == 2 {
						m = m.Where(v[0], v[1])
					}
					if len(v) == 1 {
						m = m.Where(v[0])
					}
				}
			}
		}
		// 如果ListQueryOp的Extend不为空 则执行Extend
		if s.ListQueryOp.Extend != nil {
			m = s.ListQueryOp.Extend(ctx, m)
		}
		// 如果 addOrderby 不为空 则添加排序
		if len(s.ListQueryOp.AddOrderby) > 0 && r.Get("order").IsEmpty() && r.Get("sort").IsEmpty() {
			for field, order := range s.ListQueryOp.AddOrderby {
				m = m.Order(field, order)
			}
		}
	}

	// 增加默认数据限制，防止查询所有数据
	m = m.Limit(10000)

	result, err := m.All()
	if err != nil {
		g.Log().Errorf(ctx, "ServiceList error:%v", err.Error())
	}
	if result == nil {
		data = garray.New()
	} else {
		data = result
	}
	if s.ListQueryOp != nil {
		if s.ListQueryOp.ModifyResult != nil {
			data = s.ListQueryOp.ModifyResult(ctx, data)
		}
	}
	return
}

// 分页列表
func (s *Service) ServicePage(ctx context.Context, req *PageReq) (data any, err error) {

	var (
		r            = g.RequestFromCtx(ctx)
		total        = 0
		dbRedisSlice g.SliceAny
	)

	type pagination struct {
		Page  int `json:"page"`
		Size  int `json:"size"`
		Total int `json:"total"`
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	dbRedisSlice = append(dbRedisSlice, []any{r.Router.Uri, req.Page, req.Size}...)

	m := DDAO(s.Dao, ctx)
	builder := m.Builder()
	// 如果pageQueryOp不为空 则使用pageQueryOp进行查询
	if s.PageQueryOp != nil {

		//主表别名
		if s.PageQueryOp.As != "" {
			m = m.As(s.PageQueryOp.As)
		}

		// 如果Join不为空 则添加Join
		if len(s.PageQueryOp.Join) > 0 {
			for _, join := range s.PageQueryOp.Join {
				switch join.Type {
				case LeftJoin:
					m = m.LeftJoin(join.Model.TableName(), join.Condition).As(join.Alias)
				case RightJoin:
					m = m.RightJoin(join.Model.TableName(), join.Condition).As(join.Alias)
				case InnerJoin:
					m = m.InnerJoin(join.Model.TableName(), join.Condition).As(join.Alias)
				}
			}
		}
		// 如果fileldEQ不为空 则添加查询条件
		if len(s.PageQueryOp.FieldEQ) > 0 {
			for _, field := range s.PageQueryOp.FieldEQ {
				if !r.Get(field).IsEmpty() {
					m = m.Where(field, r.Get(field))
				}
			}
		}

		// 加入where条件
		if s.PageQueryOp.Where != nil {
			where := s.PageQueryOp.Where(ctx)
			var whereSlice []string
			if len(where) > 0 {
				for _, v := range where {
					if len(v) == 3 {
						if gconv.Bool(v[2]) {
							m = m.Where(v[0], v[1])
							whereSlice = append(whereSlice, fmt.Sprintf("%s-%s", v[0], v[1]))
						}
					}
					if len(v) == 2 {
						m = m.Where(v[0], v[1])
						whereSlice = append(whereSlice, fmt.Sprintf("%s-%s", v[0], v[1]))
					}
					if len(v) == 1 {
						m = m.Where(v[0])
						whereSlice = append(whereSlice, fmt.Sprintf("%s", v[0]))
					}
				}
			}

			whereStr := gstr.Replace(gstr.JoinAny(whereSlice, "#"), " ", "&&")
			dbRedisSlice = append(dbRedisSlice, whereStr)
		}

		if s.PageQueryOp.OrWhere != nil {
			where := s.PageQueryOp.OrWhere(ctx)
			if len(where) > 0 {
				var (
					whereSlice []string
					whereStr   string
				)
				for _, v := range where {
					if len(v) == 3 {
						if gconv.Bool(v[2]) {
							whereSlice = append(whereSlice, fmt.Sprintf("%s-%s", v[0], v[1]))
						}
					}
					if len(v) == 2 {
						whereSlice = append(whereSlice, fmt.Sprintf("%s-%s", v[0], v[1]))
					}
					if len(v) == 1 {
						whereSlice = append(whereSlice, fmt.Sprintf("%s", v[0]))
					}
				}
				whereStr = gstr.Replace(gstr.JoinAny(whereSlice, "#"), " ", "&&")
				dbRedisSlice = append(dbRedisSlice, whereStr)
			}
		}

		// 如果KeyWordField不为空 则添加查询条件
		if !r.Get("keyWord").IsEmpty() {
			if len(s.PageQueryOp.KeyWordField) > 0 {
				for _, field := range s.PageQueryOp.KeyWordField {
					builder = builder.WhereOrLike(field, "%"+r.Get("keyWord").String()+"%")
				}
			}
			dbRedisSlice = append(dbRedisSlice, gstr.Trim(r.Get("keyWord").String()))
		}

		// 如果 addOrderby 不为空 则添加排序
		if len(s.PageQueryOp.AddOrderby) > 0 && r.Get("order").IsEmpty() && r.Get("sort").IsEmpty() {
			addOrderby := ""
			for field, order := range s.PageQueryOp.AddOrderby {
				m = m.Order(field, order)
				addOrderby += fmt.Sprintf("%s-%s", field, order)
			}
			dbRedisSlice = append(dbRedisSlice, addOrderby)
		}

		if Select := s.PageQueryOp.Select; Select != "" {
			m = m.Fields(Select)
		}
		// 如果PageQueryOp的Extend不为空 则执行Extend
		if s.PageQueryOp.Extend != nil {
			m = s.PageQueryOp.Extend(ctx, m)
		}
	}

	// 如果 req.Order 和 req.Sort 均不为空 则添加排序
	if !r.Get("order").IsEmpty() && !r.Get("sort").IsEmpty() {
		order := r.Get("order").String() + " " + r.Get("sort").String()
		m = m.Order(order)
		dbRedisSlice = append(dbRedisSlice, gstr.Replace(order, " ", "-"))
	}

	if DbRedisEnable {
		m = m.Cache(gdb.CacheOption{
			Duration: time.Duration(DbExpire),
			Name:     gstr.JoinAny(dbRedisSlice, "/"),
			Force:    false,
		})
	}

	m = m.Where(builder)

	var result []gdb.Record
	result, total, err = m.Offset((req.Page - 1) * req.Size).Limit(req.Size).AllAndCount(false)
	if err != nil {
		return nil, err
	}

	// 如果req.IsExport为true 则导出数据
	if req.IsExport {
		// 如果req.MaxExportSize大于0 则限制导出数据的最大条数
		if req.MaxExportLimit > 0 {
			m = m.Limit(req.MaxExportLimit)
		}
		data = g.Map{
			"list":  result,
			"total": total,
		}
		return data, nil
	} else {
		data = g.Map{
			"list": result,
			"pagination": pagination{
				Page:  req.Page,
				Size:  req.Size,
				Total: total,
			},
		}
	}

	if s.PageQueryOp != nil {
		if s.PageQueryOp.ModifyResult != nil {
			data = s.PageQueryOp.ModifyResult(ctx, data)
		}
	}

	return data, nil
}

// 新增|删除|修改前的操作
func (s *Service) ModifyBefore(ctx context.Context, method string, param g.MapStrAny) (err error) {
	return
}

// 新增|删除|修改后的操作
func (s *Service) ModifyAfter(ctx context.Context, method string, param g.MapStrAny) (err error) {
	return
}

// db 缓存处理
func (s *Service) CacheDo(ctx context.Context, method string, param g.MapStrAny) (err error) {

	var (
		r             = g.RequestFromCtx(ctx)
		containsSlice = g.Slice{}
	)
	keys, err := DbCacheManager.KeyStrings(ctx)
	if err != nil {
		return err
	}
	for _, key := range keys {
		keyStr := gconv.String(key)
		contains := gstr.Contains(keyStr, gstr.StrTillEx(keyStr, r.Router.Uri))
		if contains {
			containsSlice = append(containsSlice, keyStr)
		}
	}
	err = DbCacheManager.Removes(ctx, containsSlice)
	if err != nil {
		return err
	}
	return
}

// 获取model
func (s *Service) GetModel() IModel {
	return s.Model
}

// 获取dao
func (s *Service) GetDao() IDao {
	return s.Dao
}

func NewModelService(model IModel) *Service {
	return &Service{
		Model: model,
	}
}

func NewDaoService(dao IDao) *Service {
	return &Service{
		Dao: dao,
	}
}
