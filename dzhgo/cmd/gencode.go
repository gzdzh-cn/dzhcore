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
				Brief: "模块名称，例如: user (必须与 addons 配合使用)",
			},
			{
				Name:  "module",
				Short: "m",
				Brief: "所属模块，例如: admin 或 app，(可不填，必须与 addons 配合使用) 默认: admin和 app 同时生成",
			},
			{
				Name:  "model",
				Short: "M",
				Brief: "单独生成模型，例如: user (可与 addons 配合使用)",
			},
			{
				Name:  "controller",
				Short: "C",
				Brief: "单独生成控制器，例如: user (可与 addons 配合使用)",
			},
			{
				Name:  "logic",
				Short: "L",
				Brief: "单独生成逻辑，例如: user (可与 addons 配合使用)",
			},
		},
	}
)

// 代码生成器执行函数
func genCodeFunc(ctx context.Context, parser *gcmd.Parser) (err error) {

	addons := parser.GetOpt("addons").String()
	name := parser.GetOpt("name").String()
	module := parser.GetOpt("module").String()
	model := parser.GetOpt("model").String()
	controller := parser.GetOpt("controller").String()
	logic := parser.GetOpt("logic").String()

	// 1. 所有参数不能全为空
	if addons == "" && name == "" && module == "" && model == "" && controller == "" && logic == "" {
		return fmt.Errorf("请至少提供一个参数，使用 -a/--addons, -n/--name, -m/--module, -M/--model, -c/--controller, -l/--logic")
	}

	// 2. 有addons时，name和module可以为空，如果name为空且没有指定特定的生成参数，则用addons名称
	if addons != "" {
		// 如果name为空且没有指定特定的生成参数，使用addons作为name
		if name == "" && model == "" && controller == "" && logic == "" {
			name = addons
		}
		// 如果module为空，设置为空字符串，在generateAddonCode中会同时生成admin和app
		return generateAddonCode(addons, name, module, model, controller, logic)
	}

	// 3. 没有addons时，不能使用name和module
	if name != "" || module != "" {
		return fmt.Errorf("没有 addons 参数时，不能使用 name 和 module 参数")
	}

	// 4. 没有addons时，只能生成internal下的单独文件
	if model == "" && controller == "" && logic == "" {
		return fmt.Errorf("没有 addons 参数时，必须提供 model、controller 或 logic 参数")
	}

	return generateInternalSingleFile(model, controller, logic)
}

// 只在 internal 目录下生成单独文件
func generateInternalSingleFile(modelName, controllerName, logicName string) error {
	// 获取 go.mod 里的 module 名称
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

	basePath := "internal"
	importPrefix := modName + "/internal"

	// 生成模型（如果指定了 model）
	if modelName != "" {
		if err := generateModelAtPath(modelName, basePath, importPrefix); err != nil {
			return err
		}
		fmt.Printf("模型 %s 已生成到 %s/model\n", modelName, basePath)
	}

	// 生成控制器（如果指定了 controller）
	if controllerName != "" {
		if err := generateControllerAtPath(controllerName, "admin", basePath, importPrefix); err != nil {
			return err
		}
		fmt.Printf("控制器 %s 已生成到 %s/controller/admin\n", controllerName, basePath)
	}

	// 生成逻辑（如果指定了 logic）
	if logicName != "" {
		if err := generateLogicSysAtPath(logicName, basePath, importPrefix, logicName); err != nil {
			return err
		}
		fmt.Printf("逻辑 %s 已生成到 %s/logic/sys\n", logicName, basePath)
	}

	return nil
}

// 生成 addons 目录下的代码
func generateAddonCode(addons, name, module, model, controller, logic string) error {
	// 获取 go.mod 里的 module 名称
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

	basePath := filepath.Join("addons", addons)
	importPrefix := modName + "/addons/" + addons

	// 生成基础结构
	if err := generateAddonModule(addons, module); err != nil {
		return err
	}

	// 生成模型
	if model != "" {
		if err := generateModelAtPath(model, basePath, importPrefix); err != nil {
			return err
		}
		fmt.Printf("模型 %s 已生成到 %s/model\n", model, basePath)
	} else if name != "" && model == "" && controller == "" && logic == "" {
		// 只有在没有指定任何特定生成参数时，才使用 name 作为模型名
		if err := generateModelAtPath(name, basePath, importPrefix); err != nil {
			return err
		}
		fmt.Printf("模型 %s 已生成到 %s/model\n", name, basePath)
	}

	// 生成控制器
	if controller != "" {
		// 如果module为空，同时生成admin和app
		if module == "" {
			if err := generateControllerAtPath(controller, "admin", basePath, importPrefix); err != nil {
				return err
			}
			fmt.Printf("控制器 %s 已生成到 %s/controller/admin\n", controller, basePath)
			if err := generateControllerAtPath(controller, "app", basePath, importPrefix); err != nil {
				return err
			}
			fmt.Printf("控制器 %s 已生成到 %s/controller/app\n", controller, basePath)
		} else {
			if err := generateControllerAtPath(controller, module, basePath, importPrefix); err != nil {
				return err
			}
			fmt.Printf("控制器 %s 已生成到 %s/controller/%s\n", controller, basePath, module)
		}
	} else if name != "" && model == "" && controller == "" && logic == "" {
		// 只有在没有指定任何特定生成参数时，才使用 name 作为控制器名
		// 如果module为空，同时生成admin和app
		if module == "" {
			if err := generateControllerAtPath(name, "admin", basePath, importPrefix); err != nil {
				return err
			}
			fmt.Printf("控制器 %s 已生成到 %s/controller/admin\n", name, basePath)
			if err := generateControllerAtPath(name, "app", basePath, importPrefix); err != nil {
				return err
			}
			fmt.Printf("控制器 %s 已生成到 %s/controller/app\n", name, basePath)
		} else {
			if err := generateControllerAtPath(name, module, basePath, importPrefix); err != nil {
				return err
			}
			fmt.Printf("控制器 %s 已生成到 %s/controller/%s\n", name, basePath, module)
		}
	}

	// 生成逻辑
	if logic != "" {
		if err := generateLogicSysAtPath(logic, basePath, importPrefix, addons); err != nil {
			return err
		}
		fmt.Printf("逻辑 %s 已生成到 %s/logic/sys\n", logic, basePath)
	} else if name != "" && model == "" && controller == "" && logic == "" {
		// 只有在没有指定任何特定生成参数时，才使用 name 作为逻辑名
		if err := generateLogicSysAtPath(name, basePath, importPrefix, addons); err != nil {
			return err
		}
		fmt.Printf("逻辑 %s 已生成到 %s/logic/sys\n", name, basePath)
	}

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

		if err := gfile.Mkdir(fullPath); err != nil {
			return fmt.Errorf("创建目录失败: %s, 错误: %v", fullPath, err)
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

			if err := gfile.Mkdir(subDir); err != nil {
				return fmt.Errorf("创建目录失败: %s, 错误: %v", subDir, err)
			}

			filePath := filepath.Join(subDir, sub+".go")
			if gfile.Exists(filePath) {
				fmt.Printf("模型文件已存在: %s\n", filePath)
				continue
			}

			content := fmt.Sprintf("package %s\n\n// %s 插件的 %s/%s 代码\n", sub, name, parent, sub)
			if err := gfile.PutContents(filePath, content); err != nil {
				return fmt.Errorf("写入 %s 失败: %v", filePath, err)
			}

		}
	}

	// 合并 api 下 v1 子目录及 go 文件的生成
	apiV1Dir := filepath.Join(basePath, "api", "v1")
	if err := gfile.Mkdir(apiV1Dir); err != nil {
		return fmt.Errorf("创建目录失败: %s, 错误: %v", apiV1Dir, err)
	}

	filePath := filepath.Join(apiV1Dir, name+".go")
	if gfile.Exists(filePath) {
		fmt.Printf("api/v1/%s.go 文件已存在: %s\n", name, filePath)
		return nil
	}

	content := fmt.Sprintf("package v1\n\n// %s 插件的 api/v1/%s.go 代码\n", name, name)
	if err := gfile.PutContents(filePath, content); err != nil {
		return fmt.Errorf("写入 %s 失败: %v", filePath, err)
	}

	// 生成插件根目录下的 config.go
	configPath := filepath.Join(basePath, "config.go")
	if gfile.Exists(configPath) {
		fmt.Printf("配置文件已存在: %s\n", configPath)
		return nil
	}

	configContent := fmt.Sprintf("package %s\n\n// %s 插件的配置\n", name, name)
	if err := gfile.PutContents(configPath, configContent); err != nil {
		return fmt.Errorf("写入 %s 失败: %v", configPath, err)
	}

	// 生成插件根目录下的 插件名.go
	mainPath := filepath.Join(basePath, name+".go")
	if gfile.Exists(mainPath) {
		return fmt.Errorf("插件主入口文件已存在: %s", mainPath)
	}

	mainContent := fmt.Sprintf("package %s\n\n// %s 插件主入口\n", name, name)
	if err := gfile.PutContents(mainPath, mainContent); err != nil {
		return fmt.Errorf("写入 %s 失败: %v", mainPath, err)
	}

	// 生成 logic/sys/{name}.go，使用简单模板
	logicSysDir := filepath.Join(basePath, "logic", "sys")
	if gfile.Exists(logicSysDir) {
		fmt.Printf("逻辑文件已存在: %s\n", logicSysDir)
		return nil
	}

	if err := gfile.Mkdir(logicSysDir); err != nil {
		return fmt.Errorf("创建目录失败: %s, 错误: %v", logicSysDir, err)
	}

	logicPath := filepath.Join(logicSysDir, name+".go")
	if gfile.Exists(logicPath) {
		fmt.Printf("逻辑文件已存在: %s\n", logicPath)
		return nil
	}

	logicContent := fmt.Sprintf("package sys\n\n// %s 插件的 logic/sys/%s.go 代码\n", name, name)
	if err := gfile.PutContents(logicPath, logicContent); err != nil {
		return fmt.Errorf("写入 %s 失败: %v", logicPath, err)
	}

	fmt.Printf("插件模块 %s 目录结构已生成于 %s\n", name, basePath)
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
	if gfile.Exists(logicFile) {
		return fmt.Errorf("逻辑文件已存在: %s", logicFile)
	}

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
