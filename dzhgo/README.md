# DzhGO 代码生成工具 

这是一个独立的 DzhGO 代码生成工具，用于快速生成控制器、模型和服务代码。

## 功能特点

- 生成插件模块代码
- 生成模型和实体代码
- 生成服务接口和实现代码
- 支持 admin 和 app 两种模块类型
- 生成项目基础目录

## 安装

```bash
# 构建
go install github.com/gzdzh-cn/dzhcore/dzhgo@latest
```

## gen 命令操作说明

`gen` 命令用于一键生成 Go 项目的插件模块目录结构及基础代码，包括控制器、模型、逻辑、服务等，极大提升开发效率，适合中大型项目的模块化、插件化开发。

### 命令格式

```bash
dzhgo gen [参数]
```

### 支持参数

| 参数名         | 简写 | 说明                                 | 示例           |
| -------------- | ---- | ------------------------------------ | -------------- |
| --name         | -n   | 模块名称（如 user）                  | -n user        |
| --module       | -m   | 所属模块类型（admin 或 app）         | -m admin       |
| --addons       | -a   | 插件名称，生成到 /addons 目录下      | -a dict        |

### 使用示例

#### 1. 生成插件模块基础目录

```bash
dzhgo gen -a dict
```

- 作用：在 `addons/dict` 目录下生成插件的基础目录结构和部分基础代码。

#### 2. 生成插件模块下的控制器、模型、逻辑等代码

```bash
dzhgo gen -a dict -n user -m admin
```

- 作用：在 `addons/dict` 下，生成 `user` 模块的 admin 控制器、模型、logic/sys/user.go 等代码。

#### 3. 生成 app 模块的代码

```bash
dzhgo gen -a dict -n user -m app
```

- 作用：在 `addons/dict` 下，生成 `user` 模块的 app 控制器、模型、logic/sys/user.go 等代码。

### 参数说明

- `--addons`（必填）：指定插件名称，所有生成的内容会放在 `addons/插件名` 目录下。
- `--name`（必填）：指定模块名称，如 user、role 等。
- `--module`（可选）：指定模块类型，支持 `admin` 或 `app`，分别生成后台或前台的控制器代码。不填则同时生成 admin 和 app 两种类型。

### 生成内容说明

- 目录结构示例（以 `-a dict -n user -m admin` 为例）：

  ```
  addons/dict/
    ├── controller/
    │   └── admin/
    │       └── user.go
    ├── model/
    │   └── user.go
    ├── logic/
    │   └── sys/
    │       └── user.go
    ├── config.go
    └── dict.go
  ```

- 生成的代码均为基础模板，开发者可直接在此基础上补充业务逻辑。

### 常见问题

1. **未指定 --addons 参数会报错**
   - 必须通过 `-a` 或 `--addons` 指定插件名称。

2. **--module 仅支持 admin 或 app**
   - 其他值会提示错误。

3. **已存在的文件不会被覆盖**
   - 若目标文件已存在，命令会报错提示，避免误覆盖。

---

