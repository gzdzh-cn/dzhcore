package cmd

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
)

var (
	// GenCode 代码生成命令
	GenCode = &gcmd.Command{
		Name:  "gen",
		Usage: "gen",
		Brief: "生成控制器、模型和服务代码",
		Func:  genCodeFunc,
		Arguments: []gcmd.Argument{
			{
				Name:  "addons",
				Short: "a",
				Brief: "生成整个插件模块目录到 /addons 目录下，例如: user",
			},
			{
				Name:  "name",
				Short: "n",
				Brief: "模块名称，例如: user",
			},
			{
				Name:  "module",
				Short: "m",
				Brief: "所属模块，例如: admin 或 app，(可不填) 默认: admin和 app 同时生成",
			},
		},
	}
)

// 代码生成器执行函数
func genCodeFunc(ctx context.Context, parser *gcmd.Parser) (err error) {

	addons := parser.GetOpt("addons").String()
	name := parser.GetOpt("name").String()
	module := parser.GetOpt("module").String()

	if addons == "" {
		return fmt.Errorf("请提供插件名称，使用 -a 或 --addons 参数")
	}
	if name == "" {
		return fmt.Errorf("请提供模块名称，使用 -n 或 --name 参数")
	}
	if module != "" {
		// 验证模块类型
		if module != "admin" && module != "app" {
			if parser.GetOpt("addons").String() == "" {
				return fmt.Errorf("模块类型必须是 admin 或 app，或请使用 --addons 生成插件目录")
			}
		}
	}

	// 只要有 --addons，统一先生成基础结构
	if addons != "" {

		// 生成基础结构
		if err := generateAddonModule(addons, module); err != nil {
			return err
		}
		// 如果 name 和 module 都有，再补充 user.go 等
		if name != "" || module != "" {
			return genAddonWithModuleOnly(addons, name, module)
		}
		return nil
	}

	fmt.Println("代码生成完成!")
	return nil
}

// 新增：生成插件模块目录结构
func generateAddonModule(name, module string) error {
	basePath := filepath.Join("addons", name)

	// 根据 module 决定生成哪些目录
	var subDirs []string
	if module == "admin" {
		subDirs = []string{
			"controller/admin",
			"model",
			"logic",
			"service",
			"middleware",
			"funcs",
			"config",
			"consts",
			"packed",
			"resource/initjson",
			"api/v1",
			"controller",
		}
	} else if module == "app" {
		subDirs = []string{
			"controller/app",
			"model",
			"logic",
			"service",
			"middleware",
			"funcs",
			"config",
			"consts",
			"packed",
			"resource/initjson",
			"api/v1",
			"controller",
		}
	} else {
		subDirs = []string{
			"controller/admin",
			"controller/app",
			"model",
			"logic",
			"service",
			"middleware",
			"funcs",
			"config",
			"consts",
			"packed",
			"resource/initjson",
			"api/v1",
			"controller",
		}
	}
	for _, dir := range subDirs {
		fullPath := filepath.Join(basePath, dir)
		if !gfile.Exists(fullPath) {
			if err := gfile.Mkdir(fullPath); err != nil {
				return fmt.Errorf("创建目录失败: %s, 错误: %v", fullPath, err)
			}
		}
	}

	// 合并 controller 下 admin、app 子目录及 go 文件的生成
	for _, parent := range []string{"controller"} {
		var subs []string
		if module == "admin" {
			subs = []string{"admin"}
		} else if module == "app" {
			subs = []string{"app"}
		} else {
			subs = []string{"admin", "app"}
		}
		parentBase := filepath.Join(basePath, parent)
		for _, sub := range subs {
			subDir := filepath.Join(parentBase, sub)
			if !gfile.Exists(subDir) {
				if err := gfile.Mkdir(subDir); err != nil {
					return fmt.Errorf("创建目录失败: %s, 错误: %v", subDir, err)
				}
			}
			filePath := filepath.Join(subDir, sub+".go")
			content := fmt.Sprintf("package %s\n\n// %s 插件的 %s/%s 代码\n", sub, name, parent, sub)
			if err := gfile.PutContents(filePath, content); err != nil {
				return fmt.Errorf("写入 %s 失败: %v", filePath, err)
			}
		}
	}

	// 合并 api 下 v1 子目录及 go 文件的生成
	apiV1Dir := filepath.Join(basePath, "api", "v1")
	if !gfile.Exists(apiV1Dir) {
		if err := gfile.Mkdir(apiV1Dir); err != nil {
			return fmt.Errorf("创建目录失败: %s, 错误: %v", apiV1Dir, err)
		}
	}
	filePath := filepath.Join(apiV1Dir, name+".go")
	content := fmt.Sprintf("package v1\n\n// %s 插件的 api/v1/%s.go 代码\n", name, name)
	if err := gfile.PutContents(filePath, content); err != nil {
		return fmt.Errorf("写入 %s 失败: %v", filePath, err)
	}

	// 生成插件根目录下的 config.go
	configPath := filepath.Join(basePath, "config.go")
	configContent := fmt.Sprintf("package %s\n\n// %s 插件的配置\n", name, name)
	if err := gfile.PutContents(configPath, configContent); err != nil {
		return fmt.Errorf("写入 %s 失败: %v", configPath, err)
	}
	// 生成插件根目录下的 插件名.go
	mainPath := filepath.Join(basePath, name+".go")
	mainContent := fmt.Sprintf("package %s\n\n// %s 插件主入口\n", name, name)
	if err := gfile.PutContents(mainPath, mainContent); err != nil {
		return fmt.Errorf("写入 %s 失败: %v", mainPath, err)
	}

	// 生成 logic/sys/{name}.go，用 dict.go 的模板
	if err := generateLogicSysAtPath(name, basePath, fmt.Sprintf("addons/%s", name), name); err != nil {
		return err
	}

	fmt.Printf("插件模块 %s 目录结构已生成于 %s\n", name, basePath)
	return nil
}

// 支持 --addons --name --module 生成到 addons 下的插件模块
func genAddonWithModuleOnly(addons, name, module string) error {
	basePath := filepath.Join("addons", addons)

	// 获取 go.mod 里的 module 名称，供 importPrefix 用
	modName := ""
	if modData := gfile.GetContents("go.mod"); modData != "" {
		for _, line := range gstr.Split(modData, "\n") {
			if gstr.HasPrefix(line, "module ") {
				modName = gstr.Trim(gstr.TrimLeftStr(line, "module "))
				break
			}
		}
	}
	if modName == "" {
		return fmt.Errorf("无法获取go.mod中的module名称")
	}
	importPrefix := modName + "/addons/" + addons

	// 根据 module 参数生成 controller 下的 {name}.go，使用统一模板
	var controllerSubs []string
	if module == "admin" {
		controllerSubs = []string{"admin"}
	} else if module == "app" {
		controllerSubs = []string{"app"}
	} else {
		controllerSubs = []string{"admin", "app"}
	}
	for _, sub := range controllerSubs {
		if err := generateControllerAtPath(name, sub, basePath, importPrefix); err != nil {
			return err
		}
	}
	// 生成 model/user.go
	if err := generateModelAtPath(name, basePath, importPrefix); err != nil {
		return err
	}
	// 生成 logic/sys/{name}.go，用 dict.go 的模板
	if err := generateLogicSysAtPath(name, basePath, importPrefix, addons); err != nil {
		return err
	}
	return nil
}

// 新增：支持自定义 basePath 和 importPrefix 的控制器生成函数
func generateControllerAtPath(name, module, basePath, importPrefix string) error {
	modName := ""
	if modData := gfile.GetContents("go.mod"); modData != "" {
		for _, line := range gstr.Split(modData, "\n") {
			if gstr.HasPrefix(line, "module ") {
				modName = gstr.Trim(gstr.TrimLeftStr(line, "module "))
				break
			}
		}
	}
	if modName == "" {
		return fmt.Errorf("无法获取go.mod中的module名称")
	}
	upperName := gstr.UcFirst(name)
	controllerDir := filepath.Join(basePath, "controller", module)
	if err := gfile.Mkdir(controllerDir); err != nil && !gfile.Exists(controllerDir) {
		return fmt.Errorf("创建目录失败: %s, 错误: %v", controllerDir, err)
	}
	controllerPath := filepath.Join(controllerDir, name+".go")
	if gfile.Exists(controllerPath) {
		return fmt.Errorf("控制器文件已存在: %s", controllerPath)
	}
	const controllerTemplate = `package %s

import (
	"context"
 
	logic "%s/logic/sys"

	"github.com/gzdzh-cn/dzhcore"
)

type %sController struct {
	*dzhcore.Controller
}

func init() {
	var %sController = &%sController{
		&dzhcore.Controller{
			Prefix:  "/%s/%s",
			Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page"},
			Service: logic.New%sService(),
		},
	}

	// 注册路由
	dzhcore.AddController(%sController)
}
`
	content := fmt.Sprintf(
		controllerTemplate,
		module,
		importPrefix,
		upperName,
		name, upperName,
		module, name,
		upperName,
		name,
	)
	if err := gfile.PutContents(controllerPath, content); err != nil {
		return fmt.Errorf("写入控制器文件失败: %v", err)
	}
	fmt.Printf("生成控制器: %s\n", controllerPath)
	return nil
}

// 新增：支持自定义 basePath 和 importPrefix 的模型生成函数
func generateModelAtPath(name, basePath, importPrefix string) error {
	upperName := gstr.UcFirst(name)
	modelDir := filepath.Join(basePath, "model")
	modelPath := filepath.Join(modelDir, name+".go")
	if !gfile.Exists(modelDir) {
		if err := gfile.Mkdir(modelDir); err != nil {
			return fmt.Errorf("创建目录失败: %s, 错误: %v", modelDir, err)
		}
	}
	if gfile.Exists(modelPath) {
		return fmt.Errorf("模型文件已存在: %s", modelPath)
	}
	modelTemplate := "package model\n\nimport (\n\t\"github.com/gzdzh-cn/dzhcore\"\n)\n\nconst TableName%s = \"addons_%s\"\n\n// %s 模型，映射表 <addons_%s>\ntype %s struct {\n\t*dzhcore.Model\n\tName     string  `gorm:\"column:name;type:varchar(255);not null\" json:\"name\"` // 名称\n\tValue    string  `gorm:\"column:value;type:varchar(255)\" json:\"value\"`        // 值\n\tOrderNum int32   `gorm:\"column:orderNum;type:int;not null\" json:\"orderNum\"`  // 排序\n\tRemark   *string `gorm:\"column:remark;type:varchar(255)\" json:\"remark\"`      // 备注\n}\n\n// TableName %s 的表名\nfunc (*%s) TableName() string {\n\treturn TableName%s\n}\n\n// GroupName %s 的表分组\nfunc (*%s) GroupName() string {\n\treturn \"default\"\n}\n\n// New%s 创建一个新的 %s 实例\nfunc New%s() *%s {\n\treturn &%s{\n\t\tModel: dzhcore.NewModel(),\n\t}\n}\n\n// init 注册模型\nfunc init() {\n\tdzhcore.AddModel(&%s{})\n}\n"
	modelContent := fmt.Sprintf(modelTemplate,
		upperName, name, // const TableName%s = "addons_%s"
		upperName, name, // 注释
		upperName,            // type %s struct
		upperName,            // TableName 注释
		upperName, upperName, // TableName 方法
		upperName, upperName, // GroupName 方法
		upperName, upperName, // New%s 创建
		upperName, upperName, upperName, // New%s() *%s { return &%s{...}
		upperName, // init 注册
	)
	if err := gfile.PutContents(modelPath, modelContent); err != nil {
		return fmt.Errorf("写入模型文件失败: %v", err)
	}
	fmt.Printf("生成模型: %s\n", modelPath)
	return nil
}

// 生成服务代码
func generateService(name, module, basePath string) error {
	// 首字母大写的模块名
	upperName := gstr.UcFirst(name)

	// 服务接口路径
	serviceDir := filepath.Join(basePath, "internal", "service")
	servicePath := filepath.Join(serviceDir, name+".go")

	// 检查目录是否存在
	if !gfile.Exists(serviceDir) {
		if err := gfile.Mkdir(serviceDir); err != nil {
			return fmt.Errorf("创建目录失败: %s, 错误: %v", serviceDir, err)
		}
	}

	// 检查文件是否已存在
	if gfile.Exists(servicePath) {
		return fmt.Errorf("服务接口文件已存在: %s", servicePath)
	}

	// 逻辑实现路径
	logicDir := filepath.Join(basePath, "internal", "logic", module)
	logicPath := filepath.Join(logicDir, name+".go")

	// 检查目录是否存在
	if !gfile.Exists(logicDir) {
		if err := gfile.Mkdir(logicDir); err != nil {
			return fmt.Errorf("创建目录失败: %s, 错误: %v", logicDir, err)
		}
	}

	// 检查文件是否已存在
	if gfile.Exists(logicPath) {
		return fmt.Errorf("服务逻辑文件已存在: %s", logicPath)
	}

	// 服务接口模板
	serviceTemplate := `package service

import (
	"context"
	
	"internal/model"
	"internal/service/internal/%s"
)

// %sService 接口
type %sService interface {
	GetById(ctx context.Context, id uint) (*model.%s, error)
	GetList(ctx context.Context, page, size int) ([]*model.%s, int, error)
	Add(ctx context.Context, data *model.%s) error
	Update(ctx context.Context, data *model.%s) error
	Delete(ctx context.Context, id uint) error
}

// %s 获取服务实例
func %s() %sService {
	return %s.New()
}
`

	// 服务内部实现目录
	internalServiceDir := filepath.Join(basePath, "internal", "service", "internal", name)
	if !gfile.Exists(internalServiceDir) {
		if err := gfile.Mkdir(internalServiceDir); err != nil {
			return fmt.Errorf("创建目录失败: %s, 错误: %v", internalServiceDir, err)
		}
	}

	// 服务内部实现文件
	internalServicePath := filepath.Join(internalServiceDir, name+".go")

	// 检查文件是否已存在
	if gfile.Exists(internalServicePath) {
		return fmt.Errorf("服务内部实现文件已存在: %s", internalServicePath)
	}

	// 服务内部实现模板
	internalServiceTemplate := `package %s

import (
	"context"
	
	"internal/dao"
	"internal/model"
	"internal/service"
	
	"github.com/gogf/gf/v2/frame/g"
)

type s%sService struct{}

func New() service.%sService {
	return &s%sService{}
}

// GetById 根据ID获取
func (s *s%sService) GetById(ctx context.Context, id uint) (*model.%s, error) {
	var result *model.%s
	err := dao.Db().Model(result).Where("id = ?", id).Scan(&result)
	return result, err
}

// GetList 获取列表
func (s *s%sService) GetList(ctx context.Context, page, size int) ([]*model.%s, int, error) {
	var list []*model.%s
	m := dao.Db().Model(&model.%s{})
	
	// 获取总数
	var total int
	err := m.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	// 获取列表
	err = m.Limit(size).Offset((page - 1) * size).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	
	return list, total, nil
}

// Add 添加
func (s *s%sService) Add(ctx context.Context, data *model.%s) error {
	return dao.Db().Create(data).Error
}

// Update 更新
func (s *s%sService) Update(ctx context.Context, data *model.%s) error {
	return dao.Db().Save(data).Error
}

// Delete 删除
func (s *s%sService) Delete(ctx context.Context, id uint) error {
	return dao.Db().Delete(&model.%s{}, id).Error
}
`

	// 逻辑实现模板
	logicTemplate := `package %s

import (
	"context"
	
	"internal/model"
	"internal/service"
	
	"github.com/gogf/gf/v2/frame/g"
)

// %s模块逻辑处理

func init() {
	// 在这里可以添加初始化逻辑
	g.Log().Debug(context.Background(), "%s模块初始化")
}
`

	// 格式化服务模板
	serviceContent := fmt.Sprintf(serviceTemplate,
		name,
		upperName, upperName,
		upperName, upperName, upperName, upperName,
		upperName, upperName, upperName, name,
	)

	// 格式化服务内部实现模板
	internalServiceContent := fmt.Sprintf(internalServiceTemplate,
		name,
		upperName, upperName, upperName, upperName, upperName, upperName,
		upperName, upperName, upperName, upperName,
		upperName, upperName,
		upperName, upperName,
		upperName, upperName,
	)

	// 格式化逻辑实现模板
	logicContent := fmt.Sprintf(logicTemplate, module, upperName, upperName)

	// 写入服务接口文件
	if err := gfile.PutContents(servicePath, serviceContent); err != nil {
		return fmt.Errorf("写入服务接口文件失败: %v", err)
	}

	// 写入服务内部实现文件
	if err := gfile.PutContents(internalServicePath, internalServiceContent); err != nil {
		return fmt.Errorf("写入服务内部实现文件失败: %v", err)
	}

	// 写入逻辑实现文件
	if err := gfile.PutContents(logicPath, logicContent); err != nil {
		return fmt.Errorf("写入逻辑实现文件失败: %v", err)
	}

	fmt.Printf("生成服务接口: %s\n", servicePath)
	fmt.Printf("生成服务实现: %s\n", internalServicePath)
	fmt.Printf("生成逻辑实现: %s\n", logicPath)
	return nil
}

// 新增：支持自定义 basePath 和 importPrefix 的 logic/sys 生成函数，内容为 dict.go 模板
func generateLogicSysAtPath(name, basePath, importPrefix, addon string) error {
	const logicSysTemplate = `package sys
	
import (
	"%s/dao"
	"%s/model"
	"%s/service"
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gzdzh-cn/dzhcore"
)

func init() {
	service.Register%sService(&s%sService{})
}

type s%sService struct {
	*dzhcore.Service
}

func New%sService() *s%sService {
	return &s%sService{
		&dzhcore.Service{
			Dao:   &dao.Addons%s,
			Model: model.New%s(),
			ListQueryOp: &dzhcore.QueryOp{
				FieldEQ:      []string{""},                     // 字段等于
				KeyWordField: []string{""},                     // 模糊搜索匹配的数据库字段
				AddOrderby:   g.MapStrStr{"createTime": "ASC"}, // 添加排序
				Where: func(ctx context.Context) []g.Array { // 自定义条件
					return []g.Array{}
				},
				OrWhere: func(ctx context.Context) []g.Array { // or 自定义条件
					return []g.Array{}
				},
				Select: "",                  // 查询字段,多个字段用逗号隔开 如: id,name  或  a.id,a.name,b.name AS bname
				As:     "",                  //主表别名
				Join:   []*dzhcore.JoinOp{}, // 关联查询
				Extend: func(ctx g.Ctx, m *gdb.Model) *gdb.Model { // 追加其他条件
					return m
				},
				ModifyResult: func(ctx g.Ctx, data interface{}) interface{} { // 修改结果
					return data
				},
			},
			PageQueryOp: &dzhcore.QueryOp{
				FieldEQ:      []string{""},                     // 字段等于
				KeyWordField: []string{""},                     // 模糊搜索匹配的数据库字段
				AddOrderby:   g.MapStrStr{"createTime": "ASC"}, // 添加排序
				Where: func(ctx context.Context) []g.Array { // 自定义条件
					return []g.Array{}
				},
				OrWhere: func(ctx context.Context) []g.Array { // or 自定义条件
					return []g.Array{}
				},
				Select: "",                  // 查询字段,多个字段用逗号隔开 如: id,name  或  a.id,a.name,b.name AS bname
				As:     "",                  //主表别名
				Join:   []*dzhcore.JoinOp{}, // 关联查询
				Extend: func(ctx g.Ctx, m *gdb.Model) *gdb.Model { // 追加其他条件
					return m
				},
				ModifyResult: func(ctx g.Ctx, data interface{}) interface{} { // 修改结果
					return data
				},
			},
			InsertParam: func(ctx context.Context) g.MapStrAny { // Add时插入参数
				return g.MapStrAny{}
			},
			Before: func(ctx context.Context) (err error) { // CRUD前的操作
				return nil
			},
			InfoIgnoreProperty: "",            // Info时忽略的字段,多个字段用逗号隔开
			UniqueKey:          g.MapStrStr{}, // 唯一键 key:字段名 value:错误信息
			NotNullKey:         g.MapStrStr{}, // 非空键 key:字段名 value:错误信息
		},
	}
}
`
	upperName := gstr.UcFirst(name)
	logicSysDir := filepath.Join(basePath, "logic", "sys")
	if err := gfile.Mkdir(logicSysDir); err != nil && !gfile.Exists(logicSysDir) {
		return fmt.Errorf("创建 %s 目录失败: %v", logicSysDir, err)
	}
	logicFile := filepath.Join(logicSysDir, name+".go")
	content := fmt.Sprintf(
		logicSysTemplate,
		importPrefix,
		importPrefix,
		importPrefix,
		upperName, upperName,
		upperName,
		upperName, upperName, upperName,
		upperName, upperName,
	)
	if err := gfile.PutContents(logicFile, content); err != nil {
		return fmt.Errorf("写入 %s 失败: %v", logicFile, err)
	}
	fmt.Printf("生成逻辑实现: %s\n", logicFile)
	return nil
}
