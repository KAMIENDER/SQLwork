package database

import (
	models "SQLwork/models"
	"fmt"
	"path"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type GetGoodsController struct {
	beego.Controller
}

type GetLabelGoodsController struct {
	beego.Controller
}

type PostGoodsController struct {
	beego.Controller
}

func (this *GetGoodsController) GetAllGoods() {
	/*********************************************
	获取所有物品函数：
	从前端接收的数据：
	1、传入需要的记录数num (int)
	返回给前端的数据：用JsonResponse封装
	1、json数组为查找成功的所有goods的相关信息
	*********************************************/
	num, _ := this.GetInt("num")
	o := orm.NewOrm()
	var goods []orm.Params
	var out []orm.Params
	o.Raw("select * from goods").Values(&goods)
	t := 0
	fmt.Println(num)
	for _, term := range goods {
		if t >= num {
			break
		}
		t = t + 1
		out = append(out, term)
	}
	this.Data["json"] = out
	this.ServeJSON()
	return
}

func (this *GetLabelGoodsController) GetLabelGoods() {
	/*********************************************
	获取对应标签物品函数：
	从前端接收的数据：
	1、传入标签kind(string)，以及需要的记录数num (int)
	返回给前端的数据：用JsonResponse封装
	1、json数组为查找成功的所有goods的相关信息
	*********************************************/
	beego.Info("Get Label Goods")
	label := this.GetString("kind")
	num, _ := this.GetInt("num")
	o := orm.NewOrm()
	var out []models.Goods
	if label=="Household" {
		tarlabel := []*models.Household{}
		o.QueryTable(label).Limit(num).All(&tarlabel)
		t := 0
		for _,good := range tarlabel {
			if t >= num {
				break
			}
			t = t + 1
			temp:=models.Goods{}
			o.QueryTable("goods").Filter("id", good.Goodsid).One(&temp)
			out = append(out, temp)
		}
	} else if label=="Digital" {
		tarlabel := []*models.Digital{}
		o.QueryTable(label).Limit(num).All(&tarlabel)
		t := 0
		for _,good := range tarlabel {
			if t >= num {
				break
			}
			t = t + 1
			temp:=models.Goods{}
			o.QueryTable("goods").Filter("id", good.Goodsid).One(&temp)
			out = append(out, temp)
		}
	} else if label=="Study" {
		tarlabel := []*models.Study{}
		o.QueryTable(label).Limit(num).All(&tarlabel)
		t := 0
		for _,good := range tarlabel {
			if t >= num {
				break
			}
			t = t + 1
			temp:=models.Goods{}
			o.QueryTable("goods").Filter("id", good.Goodsid).One(&temp)
			out = append(out, temp)
		}
	} else if label=="Life" {
		tarlabel := []*models.Life{}
		o.QueryTable(label).Limit(num).All(&tarlabel)
		t := 0
		for _,good := range tarlabel {
			if t >= num {
				break
			}
			t = t + 1
			temp:=models.Goods{}
			o.QueryTable("goods").Filter("id", good.Goodsid).One(&temp)
			out = append(out, temp)
		}
	} else {
		tarlabel := []*models.Other{}
		o.QueryTable(label).Limit(num).All(&tarlabel)
		t := 0
		for _,good := range tarlabel {
			if t >= num {
				break
			}
			t = t + 1
			temp:=models.Goods{}
			o.QueryTable("goods").Filter("id", good.Goodsid).One(&temp)
			out = append(out, temp)
		}
	}
	

	this.Data["json"] = out

	this.ServeJSON()
}

func (this *PostGoodsController) PostGoods() {
	/*********************************************
	物品注册函数：
	从前端接收的数据：
	1、price价格 int
	2、describe商品描述 string
	3、photo图片文件 file
	4、quantity数量 int
	5、name名称 string
	6、label商品标签 string 
	返回给前端的数据：用JsonResponse封装
	1、status是否注册成功———— 0：失败，1：成功
	2、若注册失败的说明信息，成功返回ok
	*********************************************/
	JsonResponse:=make(map[string]interface{})
	userid:=this.Ctx.Input.GetData("user").(int64)
	price, _ := this.GetFloat("price")
	describe := this.GetString("describe")
	photo,_,err := this.GetFile("photo")
	if err!=nil {
		JsonResponse["status"] = 0
		JsonResponse["msg"] = "图片上传失败"
		this.Data["json"] = JsonResponse
		this.ServeJSON()
		return
	}
	defer photo.Close()
	quantity, _ := this.GetInt("quantity")
	name := this.GetString("name")
	label := this.GetString("label")
	status := 0
	var msg string

	var toinsert models.Goods

	toinsert.Price = price
	toinsert.Describe = describe
	toinsert.Userid = int64(userid)
	toinsert.Quantity = int64(quantity)
	toinsert.Name = name

	o := orm.NewOrm()

	var temp []*models.User
	num, _ := o.QueryTable("user").Filter("id", userid).All(&temp)
	if num <= 0 {
		status = 0
		msg = "目标用户不存在"
		JsonResponse["status"] = status
		JsonResponse["msg"] = msg
		this.Data["json"] = JsonResponse
		this.ServeJSON()
		return
	}

	toinsert.Photo = " "
	idd, _ := o.Insert(&toinsert)
	filename := strconv.Itoa(int(toinsert.Userid)) + "_" + strconv.Itoa(int(idd)) + "_" + toinsert.Name + ".jpg"
	toinsert.Photo = path.Join("static/photo", filename)
	this.SaveToFile("photo", toinsert.Photo)
	o.Update(&toinsert)
	beego.Info(label)

	if label=="Household" {
		GoodLabel:=models.Household{Goodsid:idd}
		o.Insert(&GoodLabel)
	} else if label=="Digital" {
		GoodLabel:=models.Digital{Goodsid:idd}
		o.Insert(&GoodLabel)
	} else if label=="Study" {
		GoodLabel:=models.Study{Goodsid:idd}
		_,err:=o.Insert(&GoodLabel)
		if err!=nil {
			beego.Info(err.Error())
		}
	} else if label=="Life" {
		GoodLabel:=models.Life{Goodsid:idd}
		o.Insert(&GoodLabel)
	} else {
		GoodLabel:=models.Other{Goodsid:idd}
		o.Insert(&GoodLabel)
	}
	

	JsonResponse["status"] = 1
	JsonResponse["msg"] = "ok"
	this.Data["json"] = JsonResponse
	this.ServeJSON()
	return

}
