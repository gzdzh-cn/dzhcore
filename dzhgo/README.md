# DzhGO 代码生成工具 

这是一个独立的 DzhGO 代码生成工具，用于快速生成控制器、模型和服务代码。

## 功能特点

- 初始化项目基础结构
- 生成插件模块代码
- 生成模型和实体代码
- 生成服务接口和实现代码
- 支持 admin 和 app 两种模块类型
- 生成项目基础目录

## 安装

```bash
# 首次安装
go install github.com/gzdzh-cn/dzhcore/dzhgo@latest

# 更新到最新版本
go install github.com/gzdzh-cn/dzhcore/dzhgo@latest
```


## 快速开始

### 1. 初始化项目

```bash
# 创建新项目
dzhgo init -n myproject

# 进入项目目录
cd myproject

# 初始化依赖
go mod tidy
```

### 2. 生成代码

```bash
# 生成 internal 下的模型
dzhgo gen -M user

# 生成 addons 插件
dzhgo gen -a dict -n user -m admin
```

### 3. 运行项目

```bash
# 启动服务器
go run main.go
```

## 命令说明

### init 命令 - 初始化项目结构

`init` 命令用于快速创建一个基于 GoFrame 的项目基础结构，包含完整的目录结构和基础文件。

#### 命令格式

```bash
dzhgo init [参数]
```

#### 支持参数

| 参数名   | 简写 | 说明                     | 示例                |
| -------- | ---- | ------------------------ | ------------------- |
| --name   | -n   | 项目名称（必填）         | -n myproject        |
| --path   | -p   | 项目路径（可选，默认为当前目录） | -p /path/to/project |

#### 使用示例

**示例一：在当前目录创建项目**
```bash
dzhgo init -n myproject
```

**示例二：在指定路径创建项目**
```bash
dzhgo init -n myproject -p /path/to/project
```

**示例三：使用位置参数**
```bash
dzhgo init myproject
```

#### 生成内容

- **目录结构**：
  ```
  myproject/
  ├── cmd/
  │   └── cmd.go
  ├── internal/
  │   ├── controller/
  │   │   ├── admin/
  │   │   └── app/
  │   ├── model/
  │   │   └── entity/
  │   ├── service/
  │   ├── logic/
  │   │   ├── admin/
  │   │   └── app/
  │   └── init.go
  ├── api/
  ├── main.go
  ├── go.mod
  └── README.md
  ```

- **基础文件**：
  - `main.go`: 项目入口文件
  - `go.mod`: Go 模块文件
  - `cmd/cmd.go`: 命令行入口
  - `internal/init.go`: 应用初始化文件
  - `README.md`: 项目说明文档

### gen 命令 - 生成代码文件

`gen` 命令用于一键生成 Go 项目的插件模块目录结构及基础代码，包括控制器、模型、逻辑、服务等，极大提升开发效率，适合中大型项目的模块化、插件化开发。

### version 命令 - 显示版本信息

`version` 命令用于显示当前工具的版本信息。

```bash
dzhgo version
```

### 命令格式

```bash
dzhgo gen [参数]
```

### 支持参数

| 参数名         | 简写 | 说明                                 | 示例           |
| -------------- | ---- | ------------------------------------ | -------------- |
| --addons       | -a   | 插件名称，生成到 /addons 目录下      | -a dict        |
| --name         | -n   | 模块名称（如 user）                  | -n user        |
| --module       | -m   | 所属模块类型（admin 或 app）         | -m admin       |
| --model        | -M   | 单独生成模型                         | -M user        |
| --controller   | -C   | 单独生成控制器                       | -C user        |
| --logic        | -L   | 单独生成逻辑                         | -L user        |

### 使用场景

#### 1. 生成 addons 插件代码

**场景一：只指定 addons，name 和 module 都为空**
```bash
dzhgo gen -a user
```
- 作用：在 `addons/user` 目录下生成插件的基础目录结构
- name 自动使用 addons 名称（user）
- module 为空时同时生成 admin 和 app 两种模块

**场景二：指定 addons 和 name，module 为空**
```bash
dzhgo gen -a dict -n user
```
- 作用：在 `addons/dict` 下生成 `user` 模块的代码
- module 为空时同时生成 admin 和 app 两种模块

**场景三：指定 addons 和 module，name 为空**
```bash
dzhgo gen -a dict -m admin
```
- 作用：在 `addons/dict` 下生成 admin 模块的代码
- name 自动使用 addons 名称（dict）

**场景四：指定 addons、name 和 module**
```bash
dzhgo gen -a dict -n user -m admin
```
- 作用：在 `addons/dict` 下生成 `user` 模块的 admin 控制器、模型、logic/sys/user.go 等代码

**场景五：在 addons 中单独生成特定文件**
```bash
# 单独生成模型
dzhgo gen -a dict -n user -m admin -M custom

# 单独生成控制器
dzhgo gen -a dict -n user -m admin -C custom

# 单独生成逻辑
dzhgo gen -a dict -n user -m admin -L custom

# 组合生成多个文件
dzhgo gen -a dict -n user -m admin -M custom -C custom -L custom
```

**场景六：只指定 addons 和 controller，同时生成 admin 和 app**
```bash
# 只指定 addons 和 controller，module 为空时同时生成 admin 和 app
dzhgo gen -a user -C comm
```

**场景七：只指定 addons 和 model，生成模型**
```bash
# 只指定 addons 和 model，在 addons/user 下生成 model 模型
dzhgo gen -a user -M model
```

**场景八：只指定 addons 和 logic，生成逻辑**
```bash
# 只指定 addons 和 logic，在 addons/user 下生成 logic 逻辑
dzhgo gen -a user -L logic
```

#### 2. 生成 internal 目录下的单独文件

**场景一：单独生成模型**
```bash
dzhgo gen -M user
```
- 作用：在 `internal/model` 目录下生成 `user.go` 模型文件

**场景二：单独生成控制器**
```bash
dzhgo gen -C user
```
- 作用：在 `internal/controller/admin` 目录下生成 `user.go` 控制器文件

**场景三：单独生成逻辑**
```bash
dzhgo gen -L user
```
- 作用：在 `internal/logic/sys` 目录下生成 `user.go` 逻辑文件

**场景四：组合生成多个文件**
```bash
dzhgo gen -M user -C user -L user
```
- 作用：同时生成模型、控制器和逻辑文件

### 参数说明

#### addons 相关参数
- `--addons`（-a）：指定插件名称，所有生成的内容会放在 `addons/插件名` 目录下
- `--name`（-n）：指定模块名称，如 user、role 等。当 addons 存在且 name 为空时，自动使用 addons 名称
- `--module`（-m）：指定模块类型，支持 `admin` 或 `app`。当 addons 存在且 module 为空时，同时生成 admin 和 app 两种类型

#### 单独生成参数
- `--model`（-M）：单独生成模型文件
- `--controller`（-C）：单独生成控制器文件
- `--logic`（-L）：单独生成逻辑文件

### 生成内容说明

#### addons 插件目录结构示例（以 `-a dict -n user -m admin` 为例）

```
addons/dict/
  ├── controller/
  │   ├── admin/
  │   │   └── user.go
  │   └── app/
  │       └── user.go
  ├── model/
  │   └── user.go
  ├── logic/
  │   └── sys/
  │       └── user.go
  ├── service/
  ├── middleware/
  ├── funcs/
  ├── config/
  ├── consts/
  ├── packed/
  ├── resource/
  │   └── initjson/
  ├── api/
  │   └── v1/
  │       └── dict.go
  ├── config.go
  └── dict.go
```

#### internal 目录结构示例（以 `-M user -C user -L user` 为例）

```
internal/
  ├── model/
  │   └── user.go
  ├── controller/
  │   └── admin/
  │       └── user.go
  └── logic/
      └── sys/
          └── user.go
```

### 使用规则

1. **有 addons 参数时**：
   - name 和 module 参数可以搭配使用
   - 如果 name 为空且没有指定特定的生成参数（model/controller/logic），则使用 addons 名称
   - 如果 module 为空，则同时生成 admin 和 app
   - 如果 module 指定为 "admin" 或 "app"，则只生成对应的模块
   - 可以只指定 addons 和 controller，同时生成 admin 和 app
   - 可以只指定 addons 和 model，生成模型文件
   - 可以只指定 addons 和 logic，生成逻辑文件

2. **没有 addons 参数时**：
   - name 和 module 参数都不能使用
   - 只能使用 model、controller 或 logic 参数生成 internal 下的文件

3. **model、controller、logic 可以单独使用**：
   - 只生成对应的逻辑模板
   - 可以与 addons 参数组合使用

### 常见问题

#### init 命令相关

1. **未指定项目名称会报错**
   - 必须通过 `-n` 或 `--name` 指定项目名称

2. **项目目录已存在会报错**
   - 若目标目录已存在，命令会报错提示，避免误覆盖

3. **路径参数支持相对路径和绝对路径**
   - 可以使用相对路径如 `./projects` 或绝对路径如 `/home/user/projects`

#### gen 命令相关

4. **未指定任何参数会报错**
   - 必须至少提供一个参数

5. **没有 addons 参数时使用 name 或 module 会报错**
   - 没有 addons 参数时，name 和 module 参数都不能使用

6. **已存在的文件不会被覆盖**
   - 若目标文件已存在，命令会报错提示，避免误覆盖

7. **module 参数仅支持 admin 或 app**
   - 其他值会提示错误

---

