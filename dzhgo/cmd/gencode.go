package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
)

// 下划线转斜杠
func underscoreToSlash(s string) string {
	return strings.ReplaceAll(s, "_", "/")
}

// processAddonsPathWithUnderscore 处理包含下划线的 addons 路径
// 例如：addons 名称为 "cur_pro"，输入 "cur_pro_user" 时
// 只转换 addons 名称后面的下划线，保持 addons 名称内的下划线不变
func processAddonsPathWithUnderscore(str string, undersAddons string) string {
	if str == "" || undersAddons == "" {
		fmt.Println("str or undersAddons is empty")
		return underscoreToSlash(str)
	}

	// 如果字符串以 addons 名称开头
	if gstr.HasPrefix(str, undersAddons) {
		// 获取 addons 名称后面的部分
		remaining := strings.TrimPrefix(str, undersAddons)

		// 如果后面还有内容且以下划线开头
		if remaining != "" && gstr.HasPrefix(remaining, "_") {
			// 移除开头的下划线，然后转换剩余部分的下划线为斜杠
			convertedRemaining := strings.ReplaceAll(remaining[1:], "_", "/")
			return undersAddons + "/" + convertedRemaining
		}

		// 如果后面没有内容，直接返回 addons 名称
		return undersAddons
	}

	// 如果不以 addons 名称开头，使用普通的转换
	return underscoreToSlash(str)
}

/*
使用示例：

假设 addons 名称为 "cur_pro"：

1. processAddonsPathWithUnderscore("cur_pro_user", "cur_pro")
   结果: "cur_pro/user"

2. processAddonsPathWithUnderscore("cur_pro_user_admin", "cur_pro")
   结果: "cur_pro/user/admin"

3. processAddonsPathWithUnderscore("cur_pro", "cur_pro")
   结果: "cur_pro"

4. processAddonsPathWithUnderscore("other_module_user", "cur_pro")
   结果: "other/module/user" (不匹配 addons 名称，使用普通转换)

核心逻辑：
- 保持 addons 名称 "cur_pro" 中的下划线不变
- 只转换 addons 名称后面的下划线为斜杠
- 这样生成的路径结构更符合预期
*/

// 识别-a 的参数 A，输入一个字符，转换这个字符，下划线转换成斜杠的时候，先排除掉A字符不转换，把剩下的字符斜下划线转斜杠

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
			// {
			// 	Name:  "name",
			// 	Short: "n",
			// 	Brief: "模块名称，例如: user (必须与 addons 配合使用)",
			// },
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

	addons := parser.GetOpt("addons").String()         //addons 名称
	module := parser.GetOpt("module").String()         //所属模块
	model := parser.GetOpt("model").String()           //模型名称
	controller := parser.GetOpt("controller").String() //控制器名称
	logic := parser.GetOpt("logic").String()           //逻辑名称

	// 1. 所有参数不能全为空
	if addons == "" && module == "" && model == "" && controller == "" && logic == "" {
		return fmt.Errorf("请至少提供一个参数，使用 -a/--addons, -m/--module, -M/--model, -c/--controller, -l/--logic")
	}

	// 2. 有addons时，module可以为空
	if addons != "" {
		// 转换下划线为驼峰命名
		modelCamel := gstr.CaseCamelLower(model)
		controllerCamel := gstr.CaseCamelLower(controller)
		logicCamel := gstr.CaseCamelLower(logic)
		addonsCamel := gstr.CaseCamelLower(addons)

		// 如果module为空，设置为空字符串，在generateAddonCode中会同时生成admin和app
		return generateAddonCode(addonsCamel, module, modelCamel, controllerCamel, logicCamel)
	}

	// 3. 没有addons时，不能使用name和module
	if module != "" {
		return fmt.Errorf("没有 addons 参数时，不能使用 module 参数")
	}

	// 4. 没有addons时，只能生成internal下的单独文件
	if model == "" && controller == "" && logic == "" {
		return fmt.Errorf("没有 addons 参数时，必须提供 model、controller 或 logic 参数")
	}

	// 转换下划线为驼峰命名
	modelCamel := gstr.CaseCamelLower(model)
	controllerCamel := gstr.CaseCamelLower(controller)
	logicCamel := gstr.CaseCamelLower(logic)
	addonsCamel := gstr.CaseCamelLower(addons)
	return generateInternalSingleFile(addonsCamel, modelCamel, controllerCamel, logicCamel)
}

// 只在 internal 目录下生成单独文件
func generateInternalSingleFile(addonsCamel, modelCamel, controllerCamel, logicCamel string) error {
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
	// undersAddons := gstr.CaseSnakeFirstUpper(addonsName)
	basePath := "internal"
	importPrefix := modName + "/internal"

	// 生成模型（如果指定了 model）
	if modelCamel != "" {
		if err := generateModelAtPath(addonsCamel, modelCamel, basePath); err != nil {
			return err
		}
		fmt.Printf("模型 %s 已生成到 %s/model\n", modelCamel, basePath)
	}

	// 生成控制器（如果指定了 controller）
	if controllerCamel != "" {
		if err := generateControllerAtPath(addonsCamel, controllerCamel, "admin", basePath, importPrefix); err != nil {
			return err
		}
		fmt.Printf("控制器 %s 已生成到 %s/controller/admin\n", controllerCamel, basePath)
	}

	// 生成逻辑（如果指定了 logic）
	if logicCamel != "" {
		if err := generateLogicSysAtPath(logicCamel, basePath, importPrefix, ""); err != nil {
			return err
		}
		fmt.Printf("逻辑 %s 已生成到 %s/logic/sys\n", logicCamel, basePath)
	}

	return nil
}

// 生成 addons 目录下的代码
func generateAddonCode(addonsCamel, module, modelCamel, controllerCamel, logicCamel string) error {
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
	//下划线命名
	undersAddons := gstr.CaseSnakeFirstUpper(addonsCamel)
	basePath := filepath.Join("addons", undersAddons)
	importPrefix := modName + "/addons/" + undersAddons

	// 生成基础结构
	if err := generateAddonModule(addonsCamel, module, basePath, importPrefix); err != nil {
		return err
	}

	// 生成控制器
	if controllerCamel != "" {
		// 如果module为空，同时生成admin和app
		if module == "" {
			if err := generateControllerAtPath(addonsCamel, controllerCamel, "admin", basePath, importPrefix); err != nil {
				return err
			}
			fmt.Printf("控制器 %s 已生成到 %s/controller/admin\n", controllerCamel, basePath)
			if err := generateControllerAtPath(addonsCamel, controllerCamel, "app", basePath, importPrefix); err != nil {
				return err
			}
			fmt.Printf("控制器 %s 已生成到 %s/controller/app\n", controllerCamel, basePath)
		} else {
			if err := generateControllerAtPath(addonsCamel, controllerCamel, module, basePath, importPrefix); err != nil {
				return err
			}
			fmt.Printf("控制器 %s 已生成到 %s/controller/%s\n", controllerCamel, basePath, module)
		}
	}

	// 生成模型
	if modelCamel != "" {
		if err := generateModelAtPath(addonsCamel, modelCamel, basePath); err != nil {
			return err
		}
		fmt.Printf("模型 %s 已生成到 %s/model\n", modelCamel, basePath)
	}

	// 生成逻辑
	if logicCamel != "" {
		if err := generateLogicSysAtPath(logicCamel, basePath, importPrefix, addonsCamel); err != nil {
			return err
		}
		fmt.Printf("逻辑 %s 已生成到 %s/logic/sys\n", logicCamel, basePath)
	}

	return nil
}

// 新增：生成插件模块基础目录结构
func generateAddonModule(addonsCamel, module, basePath, importPrefix string) error {

	//下划线命名
	undersName := gstr.CaseSnakeFirstUpper(addonsCamel)

	// 根据 module 决定生成哪些目录
	var subDirs []string
	switch module {
	case "admin":
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
	case "app":
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
	default:
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
		switch module {
		case "admin":
			subs = []string{"admin"}
		case "app":
			subs = []string{"app"}
		default:
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
				fmt.Printf("目录文件已存在: %s\n", filePath)
				continue
			}

			content := fmt.Sprintf("package %s\n\n// %s 插件的 %s/%s 代码\n", sub, undersName, parent, sub)
			if err := gfile.PutContents(filePath, content); err != nil {
				return fmt.Errorf("写入 %s 失败: %v", filePath, err)
			}

		}
	}

	// 合并 api 下 v1 子目录及 go 文件的生成
	if err := generateAddonApiV1(basePath, undersName); err != nil {
		return err
	}

	// 生成 controller 目录下的 controller.go
	if err := generateAddonController(basePath, importPrefix); err != nil {
		return err
	}

	// 生成 model 目录下的 model.go
	if err := generateAddonModel(basePath, undersName); err != nil {
		return err
	}

	// 生成插件根目录下的 config.go
	if err := generateAddonConfig(basePath, undersName); err != nil {
		return err
	}

	// 生成插件根目录下的 插件名.go
	if err := generateAddonMain(basePath, undersName, importPrefix); err != nil {
		return err
	}

	fmt.Printf("插件模块 %s 目录结构已生成于 %s\n", undersName, basePath)
	return nil
}

// 生成目录下的api/v1/api.go
func generateAddonApiV1(basePath, undersName string) error {
	apiV1Dir := filepath.Join(basePath, "api", "v1")
	if err := gfile.Mkdir(apiV1Dir); err != nil {
		return fmt.Errorf("创建目录失败: %s, 错误: %v", apiV1Dir, err)
	}

	filePath := filepath.Join(apiV1Dir, undersName+".go")
	if gfile.Exists(filePath) {
		fmt.Printf("api/v1/%s.go 文件已存在: %s\n", undersName, filePath)
		return nil
	}

	content := fmt.Sprintf("package v1\n\n// %s 插件的 api/v1/%s.go 代码\n", undersName, undersName)
	if err := gfile.PutContents(filePath, content); err != nil {
		return fmt.Errorf("写入 %s 失败: %v", filePath, err)
	}
	return nil
}

// 生成目录下的controller.go
func generateAddonController(basePath, importPrefix string) error {
	controllerDir := filepath.Join(basePath, "controller")
	if err := gfile.Mkdir(controllerDir); err != nil {
		return fmt.Errorf("创建目录失败: %s, 错误: %v", controllerDir, err)
	}

	// 生成 controller 目录下的 go 文件
	filePath := filepath.Join(controllerDir, "controller.go")
	if gfile.Exists(filePath) {
		fmt.Printf("controller/controller.go 文件已存在: %s\n", filePath)
		return nil
	}

	templaleData := `package controller

import (
	_ "%s/controller/admin"
	_ "%s/controller/app"
)
`

	content := fmt.Sprintf(
		templaleData,
		importPrefix,
		importPrefix,
	)
	if err := gfile.PutContents(filePath, content); err != nil {
		return fmt.Errorf("写入 %s 失败: %v", filePath, err)
	}
	return nil
}

// 生成目录下的model.go
func generateAddonModel(basePath, undersName string) error {
	modelDir := filepath.Join(basePath, "model")
	if err := gfile.Mkdir(modelDir); err != nil {
		return fmt.Errorf("创建目录失败: %s, 错误: %v", modelDir, err)
	}

	// 生成 model 目录下的 go 文件
	filePath := filepath.Join(modelDir, "model.go")
	if gfile.Exists(filePath) {
		fmt.Printf("model/model.go 文件已存在: %s\n", filePath)
		return nil
	}

	content := fmt.Sprintf("package model\n\n// %s 插件的 model/model.go 代码\n", undersName)
	if err := gfile.PutContents(filePath, content); err != nil {
		return fmt.Errorf("写入 %s 失败: %v", filePath, err)
	}
	return nil
}

// 生成插件根目录下的 config.go
func generateAddonConfig(basePath, undersName string) error {
	configPath := filepath.Join(basePath, "config.go")
	if gfile.Exists(configPath) {
		fmt.Printf("配置文件已存在: %s\n", configPath)
		return nil
	}

	templaleData := `package %s

	import "github.com/gzdzh-cn/dzhcore"
	var (
		Version = "v1.0.0"
	)

	func init() {
		dzhcore.SetVersions("%s", Version)
	}
`
	content := fmt.Sprintf(templaleData, undersName, undersName)
	if err := gfile.PutContents(configPath, content); err != nil {
		return fmt.Errorf("写入 %s 失败: %v", configPath, err)
	}

	return nil
}

// 生成插件根目录下的插件名.go
func generateAddonMain(basePath, undersName, importPrefix string) error {
	mainPath := filepath.Join(basePath, undersName+".go")
	if gfile.Exists(mainPath) {
		fmt.Printf("插件主入口文件已存在: %s\n", mainPath)
		return nil
	}
	addonsCamel := gstr.CaseCamelLower(undersName)
	lcFirstAddonsCamel := gstr.LcFirst(addonsCamel)
	templaleData := `package %s

	import (
		"github.com/gogf/gf/v2/frame/g"
		"github.com/gogf/gf/v2/os/gctx"
		"github.com/gzdzh-cn/dzhcore"
		_ "%s/controller"
		_ "%s/model"
	)
	
	func init() {
		dzhcore.AddAddon(&%sAddon{Version: Version, Name: "%s"})
	}
	
	type %sAddon struct {
		Version string
		Name    string
	}
	
	func (a *%sAddon) GetName() string {
		return a.Name
	}
	
	func (a *%sAddon) GetVersion() string {
		return a.Version
	}
	
	func (a *%sAddon) NewInit() {
		var (
			ctx = gctx.GetInitCtx()
		)
		g.Log().Debug(ctx, "------------ addon %s init start ...")
		g.Log().Debugf(ctx, "%s version:%%v", Version)
	}
	
	`

	content := fmt.Sprintf(
		templaleData,
		undersName,
		importPrefix,
		importPrefix,
		lcFirstAddonsCamel,
		undersName,
		lcFirstAddonsCamel,
		lcFirstAddonsCamel,
		lcFirstAddonsCamel,
		lcFirstAddonsCamel,
		undersName,
		undersName,
	)
	if err := gfile.PutContents(mainPath, content); err != nil {
		return fmt.Errorf("写入 %s 失败: %v", mainPath, err)
	}
	return nil
}

// 新增：支持自定义 basePath 和 importPrefix 的模型生成函数
func generateModelAtPath(addonsCamel, modelCamel, basePath string) error {

	modelTemplate := `package model

import (
	"github.com/gzdzh-cn/dzhcore"
)

const TableName%s = "addons_%s"

// %s 模型，映射表 <addons_%s>
type %s struct {
	*dzhcore.Model
	Title     string  %cgorm:"column:title;comment:标题;type:varchar(255);not null" json:"title"%c // 标题
	Status   int     %cgorm:"column:status;comment:状态;type:int(11);default:1" json:"status"%c // 状态
	OrderNum int32   %cgorm:"column:order_num;comment:排序;type:int;not null;default:99" json:"orderNum"%c  // 排序
	Remark   *string %cgorm:"column:remark;comment:备注;type:varchar(255)" json:"remark"%c      // 备注
}

// TableName %s 的表名
func (*%s) TableName() string {
	return TableName%s
}

// GroupName %s 的表分组
func (*%s) GroupName() string {
	return "default"
}

// New%s 创建一个新的 %s 实例
func New%s() *%s {
	return &%s{
		Model: dzhcore.NewModel(),
	}
}

// init 注册模型
func init() {
	dzhcore.AddModel(&%s{})
}
`

	// 驼峰转下划线
	undersAddonsName := gstr.CaseSnakeFirstUpper(addonsCamel)

	modelDir := filepath.Join(basePath, "model")

	// 下划线命名
	undersModelName := gstr.CaseSnakeFirstUpper(modelCamel)
	if addonsCamel != "" {
		undersModelName = undersAddonsName + "_" + gstr.CaseSnakeFirstUpper(modelCamel)
	}
	modelPath := filepath.Join(modelDir, undersModelName+".go")
	if !gfile.Exists(modelDir) {
		if err := gfile.Mkdir(modelDir); err != nil {
			return fmt.Errorf("创建目录失败: %s, 错误: %v", modelDir, err)
		}
	}
	if gfile.Exists(modelPath) {
		fmt.Printf("模型文件已存在: %s\n", modelPath)
		return nil
	}

	// 驼峰
	caseCamelTableName := gstr.CaseCamel(undersModelName)
	modelContent := fmt.Sprintf(modelTemplate,
		caseCamelTableName, undersModelName, // 第1-2个: TableName%s, "addons_%s"
		caseCamelTableName, undersModelName, // 第3-4个: // %s 模型, <addons_%s>
		caseCamelTableName,                     // 第5个: type %s struct
		'`', '`', '`', '`', '`', '`', '`', '`', // 第6-13个: 8 个 struct tag 反引号
		caseCamelTableName, caseCamelTableName, caseCamelTableName, // 第14-16个: // TableName, func (*), return TableName
		caseCamelTableName, caseCamelTableName, // 第17-18个: // GroupName, func (*)
		caseCamelTableName, caseCamelTableName, caseCamelTableName, caseCamelTableName, caseCamelTableName, // 第19-23个: // New, 实例, func New, *%s, &%s{
		caseCamelTableName, // 第24个: dzhcore.AddModel(&%s{})
	)
	if err := gfile.PutContents(modelPath, modelContent); err != nil {
		return fmt.Errorf("写入模型文件失败: %v", err)
	}
	fmt.Printf("生成模型: %s\n", modelPath)
	return nil
}

// 新增：支持自定义 basePath 和 importPrefix 的 logic/sys 生成函数，内容为 dict.go 模板
func generateLogicSysAtPath(logicCamel, basePath, importPrefix, addonsCamel string) error {

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

func News%sService() *s%sService {
	return &s%sService{
		&dzhcore.Service{
			Dao:   &dao.%s,
			Model: model.New%s(),
			ListQueryOp: &dzhcore.QueryOp{
				FieldEQ:      []string{""},                     // 字段等于
				KeyWordField: []string{""},                     // 模糊搜索匹配的数据库字段
				AddOrderby:   g.MapStrStr{"createTime": "DESC"}, // 添加排序
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
				ModifyResult: func(ctx g.Ctx, data any) any { // 修改结果
					return data
				},
			},
			PageQueryOp: &dzhcore.QueryOp{
				FieldEQ:      []string{""},                     // 字段等于
				KeyWordField: []string{""},                     // 模糊搜索匹配的数据库字段
				AddOrderby:   g.MapStrStr{"createTime": "DESC"}, // 添加排序
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
				ModifyResult: func(ctx g.Ctx, data any) any { // 修改结果
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

func (s *s%sService) Test(ctx context.Context) (err error) {
	return nil
}
`

	// 驼峰转下划线
	undersAddonsName := gstr.CaseSnakeFirstUpper(addonsCamel)

	// 首字母大写驼峰命名
	ucFirstLogicCamel := gstr.UcFirst(logicCamel)
	logicSysDir := filepath.Join(basePath, "logic", "sys")
	if err := gfile.Mkdir(logicSysDir); err != nil && !gfile.Exists(logicSysDir) {
		return fmt.Errorf("创建 %s 目录失败: %v", logicSysDir, err)
	}

	//下划线命名
	undersName := gstr.CaseSnakeFirstUpper(logicCamel)
	if addonsCamel != "" {
		undersName = undersAddonsName + "_" + gstr.CaseSnakeFirstUpper(logicCamel)
	}
	logicFile := filepath.Join(logicSysDir, undersName+".go")
	if gfile.Exists(logicFile) {
		fmt.Printf("逻辑文件已存在1: %s\n", logicFile)
		return nil
	}

	daoName := ucFirstLogicCamel
	if addonsCamel != "" {
		ucFirstLogicCamel = gstr.UcFirst(addonsCamel + ucFirstLogicCamel)
		daoName = "Addons" + ucFirstLogicCamel
	}

	content := fmt.Sprintf(
		logicSysTemplate,
		importPrefix,
		importPrefix,
		importPrefix,
		ucFirstLogicCamel, ucFirstLogicCamel,
		ucFirstLogicCamel,
		ucFirstLogicCamel, ucFirstLogicCamel, ucFirstLogicCamel,
		daoName, ucFirstLogicCamel, ucFirstLogicCamel,
	)
	if err := gfile.PutContents(logicFile, content); err != nil {
		return fmt.Errorf("写入 %s 失败: %v", logicFile, err)
	}
	fmt.Printf("生成逻辑实现: %s\n", logicFile)

	return nil
}

// 新增：支持自定义 basePath 和 importPrefix 的控制器生成函数
func generateControllerAtPath(addonsCamel, controllerCamel, module, basePath, importPrefix string) error {

	const controllerTemplate = `package %s

import (

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
			Service: logic.News%sService(),
		},
	}

	// 注册路由
	dzhcore.AddController(%sController)
}
`

	// 驼峰转下划线
	undersAddonsName := gstr.CaseSnakeFirstUpper(addonsCamel)

	upperController := gstr.UcFirst(controllerCamel)
	controllerDir := filepath.Join(basePath, "controller", module)
	if err := gfile.Mkdir(controllerDir); err != nil && !gfile.Exists(controllerDir) {
		return fmt.Errorf("创建目录失败: %s, 错误: %v", controllerDir, err)
	}
	//下划线命名
	undersName := gstr.CaseSnakeFirstUpper(controllerCamel)
	if addonsCamel != "" {
		undersName = undersAddonsName + "_" + gstr.CaseSnakeFirstUpper(controllerCamel)
	}
	controllerPath := filepath.Join(controllerDir, undersName+".go")
	if gfile.Exists(controllerPath) {
		fmt.Printf("控制器文件已存在: %s\n", controllerPath)
		return nil
	}

	if addonsCamel != "" {
		upperController = gstr.UcFirst(addonsCamel) + upperController
		controllerCamel = addonsCamel + gstr.UcFirst(controllerCamel)
	}

	content := fmt.Sprintf(
		controllerTemplate,
		module,
		importPrefix,
		upperController,
		controllerCamel, upperController,
		module, processAddonsPathWithUnderscore(undersName, undersAddonsName),
		upperController,
		controllerCamel,
	)
	if err := gfile.PutContents(controllerPath, content); err != nil {
		return fmt.Errorf("写入控制器文件失败: %v", err)
	}
	fmt.Printf("生成控制器: %s\n", controllerPath)

	return nil
}
