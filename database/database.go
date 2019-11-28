package database

import (
	Models "SQLwork/models"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type DatabaseController struct {
	beego.Controller
}

func (this *DatabaseController) Get() {
	o := orm.NewOrm()
	o.Using("user")
	var user Models.User
	user.Id = 0
	user.Name = "test"
	user.Password = "1234"
	user.State = 0
	user.Token = "00"
	id, err := o.Insert(&user)

	if err == nil {
		fmt.Print(id)
	}

	// 以上演示了如何进行数据库插入，需要注意的是，这里没有写reponses body因此访问链接的时候会报错
}
