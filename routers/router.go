package routers

import (
	"SQLwork/controllers"
	"SQLwork/database"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/database", &database.DatabaseController{})

	//测试
	beego.Router("/test", &database.TestController{})
	//用户相关
	beego.Router("/register", &database.RegisterController{}, "get:Get;post:Register")
	beego.Router("/login", &database.LoginController{}, "get:Get;post:Login")
	beego.Router("/logout", &database.LogoutController{}, "get:Logout")
	beego.Router("/active/?:id", &database.RegisterController{}, "get:Active")
	beego.Router("/userinfo",&database.UserInfoController{})

	// 商品相关
	beego.Router("/goodget", &database.GetGoodsController{}, "post:GetAllGoods")
	beego.Router("/goodlabelget", &database.GetLabelGoodsController{}, "post:GetLabelGoods")
	beego.Router("/postgoods", &database.PostGoodsController{}, "post:PostGoods")
	
	//编辑商品信息
	beego.Router("/edit/:username[\\w]+/?:id",&database.EditController{})


	//设置过滤函数
	beego.InsertFilter("*",beego.BeforeRouter,database.RestfulHandler)
	beego.InsertFilter("/test",beego.BeforeRouter,database.Authenticate)
	beego.InsertFilter("/logout",beego.BeforeRouter,database.Authenticate)
	beego.InsertFilter("/userinfo",beego.BeforeRouter,database.Authenticate)
	beego.InsertFilter("/goodget",beego.BeforeRouter,database.Authenticate)
	beego.InsertFilter("/goodlabelget",beego.BeforeRouter,database.Authenticate)
	beego.InsertFilter("/postgoods",beego.BeforeRouter,database.Authenticate)
	beego.InsertFilter("/edit/:username[\\w]+/?:id",beego.BeforeRouter,database.Authenticate)

}
