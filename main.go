package main

import (
	_ "SQLwork/routers"

	"github.com/astaxie/beego/orm"

	"github.com/astaxie/beego"
)

// 基本的文件目录结构：routers声明匹配的访问路径（url）
// main.go 声明使用的全局结构以及项目的初始化
// models 包含了数据库建立连接时，所需的相同结构，以及处理函数
// controller 默认的初始路径

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	// @/ 后面跟的是数据库名
	orm.RegisterDataBase("default", "mysql", "root:root@/sql_test?charset=utf8")
}
func main() {
	beego.Run()
}
