package main

import (
	_ "SQLwork/routers"

	"github.com/astaxie/beego/orm"

	"github.com/astaxie/beego"

	//"github.com/astaxie/beego/plugins/cors"

	"github.com/astaxie/beego/context"

	_ "github.com/go-sql-driver/mysql"
)

// 基本的文件目录结构：routers声明匹配的访问路径（url）
// main.go 声明使用的全局结构以及项目的初始化
// models 包含了数据库建立连接时，所需的相同结构，以及处理函数
// controller 默认的初始路径

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	// @/ 后面跟的是数据库名
	orm.RegisterDataBase("default", "mysql", "root:root@tcp(112.125.88.184:3306)/sql_test?charset=utf8")
	//orm.RunSyncdb("default",true,false)
	var FilterGateWay = func(ctx *context.Context) {
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
		//允许访问源
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")
		//允许post访问
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin,ContentType,Authorization,accept,accept-encoding, authorization, content-type") //header的类型
		ctx.ResponseWriter.Header().Set("Access-Control-Max-Age", "1728000")
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	}
	beego.InsertFilter("*", beego.BeforeRouter, FilterGateWay)
	//路由设置
}
func main() {
	beego.Run()
}
