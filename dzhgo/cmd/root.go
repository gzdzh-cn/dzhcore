package cmd

import (
	"github.com/gogf/gf/v2/os/gcmd"
)

// Root 根命令
var (
	Root = gcmd.Command{
		Name:  "dzhgo",
		Usage: "dzhgo [COMMAND] [OPTION]",
		Brief: "DzhGO 代码生成工具",
		Description: `
DzhGO 代码生成工具，用于快速生成控制器、模型和服务代码。

## 命令说明

### init 命令 - 初始化项目结构
用于快速创建一个基于 GoFrame 的项目基础结构。

**使用示例：**
1. 在当前目录创建项目: dzhgo init -n myproject
2. 在指定路径创建项目: dzhgo init -n myproject -p /path/to/project
3. 使用位置参数: dzhgo init myproject

**参数说明：**
- -n, --name: 项目名称（必填）
- -p, --path: 项目路径（可选，默认为当前目录）

**生成内容：**
- 完整的项目目录结构
- main.go 入口文件
- go.mod 模块文件
- cmd/cmd.go 命令行入口
- internal/init.go 初始化文件
- README.md 项目说明文档

### gen 命令 - 生成代码文件
用于生成控制器、模型、逻辑等代码文件。

**使用示例：**
1. 单独生成 internal 模型: dzhgo gen -M user
2. 单独生成 internal 控制器: dzhgo gen -C user
3. 单独生成 internal 逻辑: dzhgo gen -L user
4. 在 internal 中组合生成多个文件: dzhgo gen -M user -C user -L user
5. 只指定 addons（name 和 module 都为空）: dzhgo gen -a user
6. 指定 addons 和 name（module 为空）: dzhgo gen -a user -n user
7. 指定 addons 和 module（name 为空）: dzhgo gen -a user -m admin
8. 指定 addons、name 和 module: dzhgo gen -a user -n user -m admin
9. 在 addons 中单独生成模型: dzhgo gen -a user -n user -m admin -M custom
10. 在 addons 中单独生成控制器: dzhgo gen -a user -n user -m admin -C custom
11. 在 addons 中单独生成逻辑: dzhgo gen -a user -n user -m admin -L custom
12. 在 addons 中组合生成多个文件: dzhgo gen -a user -n user -m admin -M custom -C custom -L custom
13. 只指定 addons 和 controller，同时生成 admin 和 app: dzhgo gen -a user -C comm
14. 只指定 addons 和 model，生成模型: dzhgo gen -a user -M model
15. 只指定 addons 和 logic，生成逻辑: dzhgo gen -a user -L logic

**使用规则：**
- 有 addons 参数时，name 和 module 参数可以搭配使用，如果name为空且没有指定特定的生成参数则用addons名称，如果module为空则同时生成admin和app
- 有 addons 参数时，可以只指定 addons 和 controller/model/logic，分别生成对应的文件
- 没有 addons 参数时，只能使用 model、controller 或 logic 参数生成 internal 下的文件
- model、controller、logic 可以单独使用，只生成对应的逻辑模板
`,
		Additional: `
安装和更新：
go install github.com/gzdzh-cn/dzhcore/dzhgo@latest

运行 'dzhgo COMMAND -h' 获取更多命令帮助信息。
`,
	}
)

func init() {
	// 添加所有命令到根命令
	Root.AddCommand(GenCode)
	Root.AddCommand(InitProject)
}
