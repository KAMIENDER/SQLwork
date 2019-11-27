package routers

import (
	"SQLwork/controllers"
	"SQLwork/database"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/database", &database.DatabaseController{})
}
