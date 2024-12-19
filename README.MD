# dzhcore

> 基于 cooladmin 二开版本
> 这里要感谢作者李栋

## 依赖
```shell
go get -u github.com/gzdzh-cn/dzhcore
```

## Documentation

https://pkg.go.dev/github.com/gzdzh-cn/dzhcore

### 计划更新
- [父类使用dao]
- [增加前置和后置函数可以影响req传参，从而影响update，add，delete]
- [父类入参增加自定义验证]

## 更新日志

v1.1.7 
- 更改依赖顺序

v1.1.6 2024年8月29日
- 更新工具函数

v1.1.5
- 增加工具函数
- 接口查询，id转换为数字
- service的PageQueryOp和ListQueryOp增加As参数

v1.1.4
- 自增id改为雪花ID

v1.1.3
- 修改bug

v1.1.0
- 修改ok返回自定义msg
- 修改FillInitData的逻辑
- 增加版本号打印
- service数据库调用改为model封装DBM

v1.0.1
- 加入了 mysql,本地上传驱动，oss驱动
- 修改为gf最新的dao层方式调用