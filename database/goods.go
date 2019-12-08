package database

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"fmt"
)

type GetGoodsController struct {
	beego.Controller
}

type GetLabelGoodsController struct {
	beego.Controller
}

func (this *GetGoodsController) GetAllGoods() {
	/*********************************************
	用户注册函数：
	从前端接收的数据：
	1、传入需要的记录数num (int)
	返回给前端的数据：用JsonResponse封装
	1、json数组为查找成功的所有goods的相关信息
	*********************************************/
	num,_ := this.GetInt("num")
	o := orm.NewOrm()
	var goods []orm.Params
	o.Raw("select * from user").Values(&goods)
	t := 0
	for _, term := range goods {
		t = t + 1
		fmt.Println(term["id"],":",term["name"])
		if t>=num {
			break
		}
	}
	this.Ctx.WriteString("test")
}
