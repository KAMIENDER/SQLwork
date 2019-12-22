package routers

import (
	"SQLwork/database"

	"github.com/astaxie/beego"
	// "github.com/astaxie/beego/plugins/cors"
)

func init() {
	//跨域

	// beego.Router("/", &controllers.MainController{})
	beego.Router("/database", &database.DatabaseController{})

	//测试
	//beego.Router("/test", &database.TestController{})
	//用户相关
	beego.Router("/register", &database.RegisterController{}, "post:Get;get:Register")
	beego.Router("/login", &database.LoginController{}, "post:Get;get:Login")
	beego.Router("/logout", &database.LogoutController{}, "get:Logout")
	beego.Router("/active/?:id", &database.RegisterController{}, "get:Active")

	// 商品相关
	beego.Router("/goodget", &database.GetGoodsController{}, "get:GetAllGoods")
	beego.Router("/goodlabelget", &database.GetLabelGoodsController{}, "get:GetLabelGoods")
	beego.Router("/postgoods", &database.PostGoodsController{}, "get:PostGoods")
}
