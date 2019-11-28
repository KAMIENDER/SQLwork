package database

import (
	Models "SQLwork/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type DatabaseController struct {
	beego.Controller
}

// /*这一部分是(批量)删除测试
func (this *DatabaseController) Get() {
	o := orm.NewOrm()
	o.Using("user")
	var user Models.User
	user.Name = "test"
	err := o.Read(&user, "name")
	for {
		if err != nil {
			println(err)
			break
		}
		println("in")
		count, _ := o.Delete(&user)
		println(count)
		err = o.Read(&user, "name")
	}
	this.Ctx.WriteString("test")
	// 以上演示了如何进行数据库插入，需要注意的是，这里没有写reponses body因此访问链接的时候会报错
}

// */

/* 以下部分是插入的代码
func (this *DatabaseController) Get() {
	o := orm.NewOrm()
	o.Using("user")
	var user Models.User
	user.Name = "test"
	user.Password = "1234"
	user.State = 0
	user.Token = "00"
	id, err := o.Insert(&user)

	if err == nil {
		println(id)
	}
	this.Ctx.WriteString("test")
	// 以上演示了如何进行数据库插入，需要注意的是，这里没有写reponses body因此访问链接的时候会报错
}

*/
