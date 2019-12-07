package database

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils"
	"github.com/astaxie/beego/context"
	models "SQLwork/models"
	"strconv"
	"crypto/md5"
	"math/rand"
	"time"
	"encoding/hex"
	"bytes"
)

type RegisterController struct {
	beego.Controller
}

type LoginController struct {
	beego.Controller
}

type LogoutController struct {
	beego.Controller
}

type TestController struct {
	beego.Controller
}

// type AuthController struct {	//认证相关的控制器
// 	beego.Controller
// }

func GenereteActiveSeq() string {
	//生成随机码
	rand.Seed(time.Now().UnixNano())
	seq:=0
	for i:=0;i<10;i++ {
		seq=seq*10+rand.Intn(10)
	}
	ActiveSeq:=strconv.Itoa(seq)
	return ActiveSeq
}

func Encrypt(original string) string {
	/***********************************************
	用AES算法进行加密的函数
	参数：
	1、original：原始字符串
	返回值：
	2、ciphertext：加密后的字符串
	***********************************************/
	key:=beego.AppConfig.String("key")
	crypto:=md5.New()
	crypto.Write([]byte(original))
	ciphertext:=hex.EncodeToString(crypto.Sum(nil))
	crypto.Write([]byte(ciphertext+key))
	ciphertext=hex.EncodeToString(crypto.Sum(nil))
	return ciphertext
}


func SendMail(to,ActiveSeq string) error {
	/*************************************
	发送激活邮件的函数：
	参数：
	1、to：接收对象的邮箱
	2、ActiveSeq：随机激活码
	返回值：
	1、err：是否发送成功的错误信息
	*************************************/
	config:=`{"username":"986439206@qq.com","password":"ashrmqfkwptxbfdi","host":"smtp.qq.com","port":587}`
	email:=utils.NewEMail(config)
	email.To=[]string{to}
	email.From="986439206@qq.com"
	email.Subject="校内咸鱼激活邮件"
	email.HTML="请点击该链接进行激活：http://127.0.0.1:8090/active/"+ActiveSeq
	err:=email.Send()
	return err
}

func (this *RegisterController) Get() {
	this.ServeJSON()
}

func (this *RegisterController) Register(){
	/*********************************************
	用户注册函数：
	从前端接收的数据：
	1、username账号 要求不能带'.'号
	2、password密码
	3、email邮箱
	4、wechat微信号
	5、phone手机号码
	返回给前端的数据：用JsonResponse封装
	1、status是否注册成功———— 0：失败，1：成功
	2、若注册失败的说明信息
	*********************************************/
	username:=this.GetString("username")
	password:=this.GetString("password")
	email:=this.GetString("email")
	wechat:=this.GetString("wechat")
	phone:=this.GetString("phone")
	status:=0
	var msg string 
	JsonResponse:=make(map[string]interface{})
	//查询用户名、邮箱、手机是否已注册 是的话就不能再注册
	orm:=orm.NewOrm()
	newuser:=models.User{Name:username}
	err:=orm.Read(&newuser,"Name")
	if err==nil {
		//用户名已存在
		msg="用户名已存在"
		JsonResponse["status"]=status
		JsonResponse["msg"]=msg
		this.Data["json"]=JsonResponse
		this.ServeJSON()
		return 
	}
	comm:=models.Communications{Email:email}
	err=orm.Read(&comm,"Email")
	if err==nil {
		//邮箱已存在
		msg="邮箱已存在"
		JsonResponse["status"]=status
		JsonResponse["msg"]=msg
		this.Data["json"]=JsonResponse
		this.ServeJSON()
		return
	}
	comm=models.Communications{Phone:phone}
	if err==nil {
		//手机号码已存在
		msg="手机号码已存在"
		JsonResponse["status"]=status
		JsonResponse["msg"]=msg
		this.Data["json"]=JsonResponse
		this.ServeJSON()
		return
	}
	//密码加密存储
	password=Encrypt(password)
	//State=0表示未激活 激活后State只有两种状态 一种是登录一种是未登录
	newuser=models.User{Name:username,Password:password,State:0}
	count,err:=orm.Insert(&newuser)
	if err==nil&&count>0 {
		//插入新用户的记录
		_=orm.Read(&newuser,"Name")
		userid:=newuser.Id
		ActiveSeq:=GenereteActiveSeq()
		comm=models.Communications{Email:email,Phone:phone,Wechat:wechat,Userid:userid,Activeid:ActiveSeq}
		count,err:=orm.Insert(&comm)
		if err==nil&&count>0 {
			//成功插入新用户对应的联系方式的记录 发送激活邮件
			res:=SendMail(email,ActiveSeq)
			if res!=nil {
				beego.Error("邮件发送失败",res)
				msg="邮件发送失败"
			} else {
				msg="发送成功"
				status=1
			}
		} else {
			//插入记录失败 也要删除已经插入的新用户记录
			beego.Error("插入数据失败",err)
			orm.Delete(&newuser)
			msg="无法添加联系方式"
		}
	} else {
		//注册失败
		msg="无法添加用户"
	}
	JsonResponse["status"]=status
	JsonResponse["msg"]=msg
	this.Data["json"]=JsonResponse
	this.ServeJSON()
	return
}

func (this *RegisterController) Active() {
	/***************************************************
	响应邮箱激活链接的函数，无返回值， 可以不能写view
	***************************************************/
	ActiveSeq:=this.Ctx.Input.Param(":id")	//从URL参数中获取激活码
	beego.Info("激活码："+ActiveSeq)
	orm:=orm.NewOrm()
	comm:=models.Communications{Activeid:ActiveSeq}
	err:=orm.Read(&comm,"Activeid")
	if err != nil {
		//激活码不存在
		this.Ctx.WriteString("激活码不存在")
		return 
	}
	userid:=comm.Userid
	comm.Activeid=""
	user:=models.User{Id:userid}
	err=orm.Read(&user)
	if err != nil {
		this.Ctx.WriteString("用户不存在")
		return 
	}
	user.State=1
	orm.Update(&user)
	orm.Update(&comm)
	this.Ctx.WriteString("激活成功")
}

func GenerateToken(username string) string {
	/*********************************************************
	生成token的函数：token的形式为 用户名.加密字符串
	*********************************************************/
	seq:=GenereteActiveSeq()
	token:=Encrypt(username+seq)
	token=username+"."+token
	return token
}

func (this *LoginController) Get() {
	//获取登录界面
	this.Ctx.WriteString("登录界面")
}

func (this *LoginController) Login(){
	/********************************************
	用户登录函数：
	从前端接收的数据：
	1、username（唯一性）
	2、password
	返回给前端的数据：用JsonResponse封装
	1、status：表示是否成功登录————0：失败，1：成功
	2、msg：若登录失败则记录的是失败的信息，否则记录的是token
	********************************************/
	username:=this.GetString("username")
	password:=this.GetString("password")
	status:=0
	var msg string
	JsonResponse:=make(map[string]interface{})
	orm:=orm.NewOrm()
	user:=models.User{Name:username}
	err:=orm.Read(&user,"Name")
	if err!=nil {
		JsonResponse["status"]=status
		JsonResponse["msg"]="用户不存在"
		this.Data["json"]=JsonResponse
		this.ServeJSON()
		return
	}
	if user.State==0 {
		JsonResponse["status"]=status
		JsonResponse["msg"]="用户未激活"
		this.Data["json"]=JsonResponse
		this.ServeJSON()
		return
	}
	password=Encrypt(password)
	beego.Info(password)
	beego.Info(user.Password)
	if password!=user.Password {
		JsonResponse["status"]=status
		JsonResponse["msg"]="密码不正确"
		this.Data["json"]=JsonResponse
		this.ServeJSON()
		return
	}
	//签发token，设置用户状态为登录状态
	msg=GenerateToken(username)	//msg存的是token
	this.SetSession(username,msg)
	JsonResponse["status"]=status
	JsonResponse["msg"]=msg
	this.Data["json"]=JsonResponse
	this.ServeJSON()
}

var Authenticate=func(ctx *context.Context) {
	/**********************************************
	过滤函数，用于访问页面前验证用户身份
	要求请求头包含“Authotization”字段
	**********************************************/
	token:=ctx.Input.Header("Authotization")
	var stringBuilder bytes.Buffer
	for i:=0;i<len(token);i++ {
		if token[i] == '.' {
			break
		} else {
			stringBuilder.WriteByte(token[i])
		}
	}
	username:=stringBuilder.String()
	SaveToken,err:=ctx.Input.Session(username).(string)
	beego.Info(SaveToken)
	if err==false || token!=SaveToken {
		response:=make(map[string]interface{})
		response["status"]=0
		response["msg"]="token认证失效"
		ctx.Output.JSON(response,true,true)
		return
	}
	//修改请求体的信息 向控制器传递新信息
	ctx.Input.SetData("user",username)
}

func (this *TestController) Get() {
	username:=this.Ctx.Input.GetData("user").(string)
	this.Ctx.WriteString(username)
}

func (this *LogoutController) Logout() {
	/*******************************************
	用户注销函数：删除用户对应的session
	*******************************************/
	username:=this.Ctx.Input.GetData("user").(string)
	_,Nothing:=this.GetSession(username).(string)
	JsonResponse:=make(map[string]interface{})
	if Nothing==false {
		JsonResponse["status"]=0
		JsonResponse["msg"]="用户未登录"
	} else {
		this.DelSession(username)
		JsonResponse["status"]=1
		JsonResponse["msg"]="注销成功"
	}
	this.Data["json"]=JsonResponse
	this.ServeJSON()
}

func init() {
	//设置过滤函数
	beego.InsertFilter("/test",beego.BeforeRouter,Authenticate)
	beego.InsertFilter("/logout",beego.BeforeRouter,Authenticate)
}
