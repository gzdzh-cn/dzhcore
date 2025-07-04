package dzhcore

import (
	"gorm.io/gorm"

	"github.com/gogf/gf/v2/database/gdb"

	"strings"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gzdzh-cn/dzhcore/coreconfig"
	"github.com/gzdzh-cn/dzhcore/coredb"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gres"
)

var (
	Models []IModel
)

// IModel数组
func AddModel(model IModel) {
	Models = append(Models, model)
}

func InitModels() {
	if coreconfig.Config.Core.AutoMigrate {
		g.Log().Debugf(ctx, "InitModels,数量： %v", len(Models))
		for _, model := range Models {
			g.Log().Debugf(ctx, "model: %v", model.TableName())
			CreateTable(model)
		}
	}
}

// InitDB 初始化数据库连接供gorm使用
func InitDB(group string) (*gorm.DB, error) {
	// var ctx context.Context
	var db *gorm.DB
	// 如果group为空，则使用默认的group，否则使用group参数
	if group == "" {
		group = "default"
	}
	defer func() {
		if err := recover(); err != nil {
			panic("failed to connect database")
		}
	}()
	config := g.DB(group).GetConfig()
	db, err := coredb.GetConn(config)
	if err != nil {
		panic(err.Error())
	}

	GormDBS[group] = db
	return db, nil
}

// 根据entity结构体获取 *gorm.DB
func getDBbyModel(model IModel) *gorm.DB {

	group := model.GroupName()
	// 判断是否存在 GormDBS[group] 字段，如果存在，则使用该字段的值作为DB，否则初始化DB
	if _, ok := GormDBS[group]; ok {
		return GormDBS[group]
	} else {

		db, err := InitDB(group)
		if err != nil {
			panic("failed to connect database")
		}
		// 把重新初始化的GormDBS存入全局变量中
		GormDBS[group] = db
		return db
	}
}

// CreateTable 根据entity结构体创建表
func CreateTable(model IModel) error {

	db := getDBbyModel(model)
	return db.AutoMigrate(model)

}

// FillInitData 数据库填充初始数据
func FillInitData(ctx g.Ctx, moduleName string, model IModel) {

	mInit := g.DB("default").Model("base_sys_init")
	value, err := mInit.Clone().Where("group", model.GroupName()).Where("module", moduleName).Value("tables")
	if err != nil {
		g.Log().Errorf(ctx, "读取表 base_sys_init 失败: %s", err.Error())
	}

	// 模块第一次写入数据
	if value.IsEmpty() {
		id := NodeSnowflake.Generate().String()
		if err = updateData(ctx, mInit, moduleName, model); err == nil {
			g.Log().Debugf(ctx, "分组 %v,模块 %v 中的表 %v，第一次写入", model.GroupName(), moduleName, model.TableName())
			_, err = mInit.Insert(g.Map{"id": id, "group": model.GroupName(), "module": moduleName, "tables": model.TableName()})
			if err != nil {
				g.Log().Error(ctx, err.Error())
			}
		}
		return
	}

	tableArr := strings.Split(value.String(), ",")
	tableGarry := garray.NewStrArrayFrom(tableArr)
	//写入过了，跳过
	if tableGarry.Contains(model.TableName()) {
		g.Log().Debugf(ctx, "分组 %v, 模块 %v 中的表 %v, 已经初始化过,跳过本次初始化.", model.GroupName(), moduleName, model.TableName())
		return
	}

	//更新写入
	if err = updateData(ctx, mInit, moduleName, model); err == nil {
		g.Log().Debugf(ctx, "分组 %v, 模块 %v 中的表 %v, 写入 ", model.GroupName(), moduleName, model.TableName())
		tableGarry.Append(model.TableName())
		str := strings.Join(tableGarry.Slice(), ",")
		_, err := mInit.Where("group", model.GroupName()).Where("module", moduleName).Data(g.Map{"tables": str}).Update()
		if err != nil {
			return
		}
	}

	g.Log().Debugf(ctx, "分组 %v, 模块 %v 中的表 %v, 初始化完成 ", model.GroupName(), moduleName, model.TableName())
}

// 写入文件
func updateData(ctx g.Ctx, mInit *gdb.Model, moduleName string, model IModel) error {

	m := g.DB(model.GroupName()).Model(model.TableName())
	pathName := "addons/" + moduleName
	if moduleName == "base" {
		pathName = "internal"
	}
	path := pathName + "/resource/initjson/" + model.TableName() + ".json"
	jsonData, _ := gjson.LoadContent(gres.GetContent(path))

	g.Log().Debugf(ctx, "model.TableName(): %v,path:%v", model.TableName(), path)

	if jsonData.Var().Clone().IsEmpty() {
		g.Log().Debugf(ctx, "分组%s中的表%s无可用的初始化数据,跳过本次初始化. jsonData:%v", model.GroupName(), model.TableName(), jsonData)
		return gerror.New("无可用的初始化数据,跳过本次初始化")
	}
	_, err := m.Data(jsonData).Insert()
	if err != nil {
		g.Log().Error(ctx, err.Error())
		return err
	}

	return nil
}
