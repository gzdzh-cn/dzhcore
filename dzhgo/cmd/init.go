package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
)

var (
	// InitProject 初始化项目命令
	InitProject = &gcmd.Command{
		Name:  "init",
		Usage: "init",
		Brief: "初始化项目结构",
		Func:  initProjectFunc,
		Arguments: []gcmd.Argument{
			{
				Name:  "name",
				Short: "n",
				Brief: "项目名称，例如: myproject",
			},
			{
				Name:  "path",
				Short: "p",
				Brief: "项目路径，默认为当前目录",
			},
		},
	}
)

// 初始化项目执行函数
func initProjectFunc(ctx context.Context, parser *gcmd.Parser) (err error) {
	// 获取参数
	name := parser.GetOpt("name").String()
	if name == "" {
		name = parser.GetArg(1).String()
	}
	if name == "" {
		return fmt.Errorf("请提供项目名称，使用 -n 或 --name 参数")
	}

	// 获取项目路径
	path := parser.GetOpt("path").String()
	if path == "" {
		path = parser.GetArg(2).String()
	}
	if path == "" {
		path = "." // 默认是当前目录
	}

	projectPath := filepath.Join(path, name)

	// 检查目录是否已存在
	if gfile.Exists(projectPath) {
		return fmt.Errorf("项目目录已存在: %s", projectPath)
	}

	// 创建项目目录
	fmt.Printf("正在创建项目: %s, 路径: %s\n", name, projectPath)

	// 创建项目基本目录结构
	dirs := []string{
		filepath.Join(projectPath, "internal", "controller", "admin"),
		filepath.Join(projectPath, "internal", "controller", "app"),
		filepath.Join(projectPath, "internal", "model", "entity"),
		filepath.Join(projectPath, "internal", "service"),
		filepath.Join(projectPath, "internal", "logic", "admin"),
		filepath.Join(projectPath, "internal", "logic", "app"),
		filepath.Join(projectPath, "cmd"),
		filepath.Join(projectPath, "api"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录失败: %s, 错误: %v", dir, err)
		}
	}

	// 创建基本文件
	files := map[string]string{
		filepath.Join(projectPath, "main.go"): `package main

import (
	"{{.Name}}/cmd"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	ctx := gctx.New()
	cmd.Main.Run(ctx)
}`,

		filepath.Join(projectPath, "go.mod"): `module {{.Name}}

go 1.20

require (
	github.com/gogf/gf/v2 v2.5.7
)`,

		filepath.Join(projectPath, "cmd", "cmd.go"): `package cmd

import (
	"context"
	"{{.Name}}/internal"
	
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// 初始化
			internal.Init()
			
			// 启动服务器
			s := g.Server()
			s.Run()
			
			return nil
		},
	}
)`,

		filepath.Join(projectPath, "internal", "init.go"): `package internal

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
)

// Init 初始化应用
func Init() {
	ctx := context.Background()
	g.Log().Info(ctx, "应用初始化...")
}`,

		filepath.Join(projectPath, "README.md"): "# {{.Name}}\n\n基于 GoFrame 的项目\n\n## 目录结构\n\n- cmd: 命令行入口\n- internal: 内部代码\n  - controller: 控制器\n  - model: 数据模型\n  - service: 服务接口\n  - logic: 业务逻辑\n- api: API文档\n\n## 运行项目\n\n```bash\ngo run main.go\n```",
	}

	// 写入文件
	for filePath, content := range files {
		// 替换模板变量
		content = strings.Replace(content, "{{.Name}}", name, -1)

		if err := gfile.PutContents(filePath, content); err != nil {
			return fmt.Errorf("写入文件失败: %s, 错误: %v", filePath, err)
		}
	}

	fmt.Println("项目创建完成!")
	fmt.Printf("可以通过以下命令运行项目:\n")
	fmt.Printf("cd %s\n", projectPath)
	fmt.Printf("go mod tidy\n")
	fmt.Printf("go run main.go\n")

	return nil
}
