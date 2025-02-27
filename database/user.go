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
	"regexp"
	"path"
	"os"
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

type UserInfoController struct {
	beego.Controller
}

type EditController struct {
	beego.Controller
}

type TestController struct {
	beego.Controller
}


var (
	//NmaeFormat=regexp.MustCompile(`^(\w+)$`)
	PasswordFormat=regexp.MustCompile(`^[A-Za-z0-9]+$`)
	// \w 匹配 A-Z 或 a-z 或 0-9
	//\w+ 匹配1个或多个; ()表示将表达式分组; ?表示匹配0次或1次 
	//{m,n}表示匹配前一个字符或分组m至n次
	//^表示匹配字符串开头   $表示匹配字符串结尾
	EmailFormat=regexp.MustCompile(`^(\w+)?(\.\w+)?@(\w+)?\.(\w{2,5})(\.\w{2,3})?$`)
	PhoneFormat=regexp.MustCompile(`^(1[356789]\d)(\d{4})(\d{4})$`)
)

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
	email.HTML="请点击该链接进行激活：http://112.125.88.184:8090/active/"+ActiveSeq
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

	//检查输入的格式

	if len(username)>20 || len(password)>16 {
		JsonResponse["msg"]="用户名或密码过长"
		this.Data["json"]=JsonResponse
		this.ServeJSON()
		return 
	}
	
	// if NmaeFormat.MatchString(username)==false {
	// 	JsonResponse["msg"]="用户名格式错误"
	// 	this.Data["json"]=JsonResponse
	// 	this.ServeJSON()
	// 	return 
	// }

	if PasswordFormat.MatchString(password)==false {
		JsonResponse["msg"]="密码格式错误"
		this.Data["json"]=JsonResponse
		this.ServeJSON()
		return 
	}

	if EmailFormat.MatchString(email)==false {
		JsonResponse["msg"]="邮箱格式错误"
		this.Data["json"]=JsonResponse
		this.ServeJSON()
		return 
	}

	if PhoneFormat.MatchString(wechat)==false {
		JsonResponse["msg"]="微信号格式错误"
		this.Data["json"]=JsonResponse
		this.ServeJSON()
		return 
	}

	if PhoneFormat.MatchString(phone)==false {
		JsonResponse["msg"]="手机号码格式错误"
		this.Data["json"]=JsonResponse
		this.ServeJSON()
		return 
	}

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
				orm.Delete(&newuser)
			} else {
				msg="邮件发送成功"
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

func GenerateToken(userid int64) string {
	/*********************************************************
	生成token的函数：token的形式为 用户名.加密字符串
	*********************************************************/
	id:=strconv.FormatInt(userid,10)
	seq:=GenereteActiveSeq()
	token:=Encrypt(id+seq)
	token=id+"."+token
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
	//检查输入的格式
	if len(username)>20 || len(password)>16 {
		JsonResponse["msg"]="用户名或密码过长"
		this.Data["json"]=JsonResponse
		this.ServeJSON()
		return 
	}

	// if NmaeFormat.MatchString(username)==false {
	// 	JsonResponse["msg"]="用户名格式错误"
	// 	this.Data["json"]=JsonResponse
	// 	this.ServeJSON()
	// 	return 
	// }

	if PasswordFormat.MatchString(password)==false {
		JsonResponse["msg"]="密码格式错误"
		this.Data["json"]=JsonResponse
		this.ServeJSON()
		return 
	}

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
	if password!=user.Password {
		JsonResponse["status"]=status
		JsonResponse["msg"]="密码不正确"
		this.Data["json"]=JsonResponse
		this.ServeJSON()
		return
	}
	//签发token，设置用户状态为登录状态
	msg=GenerateToken(user.Id)	//msg存的是token
	beego.Info(msg)
	this.SetSession(user.Id,msg)
	JsonResponse["status"]=1
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
	userid:=stringBuilder.String()
	user,_:= strconv.ParseInt(userid, 10, 64)
	SaveToken,err:=ctx.Input.Session(user).(string)
	if err==false || token!=SaveToken {
		response:=make(map[string]interface{})
		response["status"]=0
		if err==false {
			response["msg"]="token认证失效"
		} else {
			response["msg"]="token不一致"
		}
		ctx.Output.JSON(response,true,true)
		return
	}
	//修改请求体的信息 向控制器传递新信息
	ctx.Input.SetData("user",user)
}

func (this *TestController) Get() {
	userid:=this.Ctx.Input.GetData("user").(int64)
	user:=strconv.FormatInt(userid,10)
	this.Ctx.WriteString(user)
}

func (this *LogoutController) Logout() {
	/*******************************************
	用户注销函数：删除用户对应的session
	*******************************************/
	userid:=this.Ctx.Input.GetData("user").(int64)
	_,err:=this.GetSession(userid).(string)
	JsonResponse:=make(map[string]interface{})
	if err==false {
		JsonResponse["status"]=0
		JsonResponse["msg"]="用户未登录"
	} else {
		this.DelSession(userid)
		JsonResponse["status"]=1
		JsonResponse["msg"]="注销成功"
	}
	this.Data["json"]=JsonResponse
	this.ServeJSON()
}

func (this *UserInfoController) Get() {
	/****************************************
	获取用户个人信息的函数：
	前端传入值：无
	返回值：
	1、用户名username
	2、商品列表goods：包括商品名，商品图片(路径)，商品编辑超链接
	若用户无商品，则goods返回的是空字符串
	****************************************/
	beego.Info("UserInfo")
	userid:=this.Ctx.Input.GetData("user").(int64)
	JsonResponse:=make(map[string]interface{})
	user:=models.User{Id:userid}
	OrmQuery:=orm.NewOrm()
	OrmQuery.Read(&user)
	JsonResponse["username"]=user.Name
	GoodsTable:=OrmQuery.QueryTable("goods")
	var goods=[]*models.Goods{}
	n,err:=GoodsTable.Filter("userid",userid).All(&goods)
	if err==nil && n>0 {
		var GoodsResponse []map[string]string
		for i:=0;i<len(goods);i++ {
			GoodName:=goods[i].Name
			GoodId:=goods[i].Id
			//112.125.88.184
			EditUrl:="112.125.88.184:8090/edit/"+user.Name+"/"+strconv.FormatInt(GoodId,10)	//10进制形式转为GoodId
			good:=make(map[string]string)
			good[GoodName]=EditUrl
			GoodsResponse=append(GoodsResponse,good)
		}
		JsonResponse["goods"]=GoodsResponse
	} else {
		JsonResponse["goods"]=""	//无商品
	}
	this.Data["json"]=JsonResponse
	this.ServeJSON()
}

func (this *EditController) Get() {
	/*****************************************************
	获取指定商品信息的函数
	前端传入值：无 （直接点击超链接访问的）
	返回值：
	1、status：status=0表示没有该商品或当前用户没有修改权限，不返回下面其他数据；status=1表示商品存在
	2、name：商品名称
	3、price：商品价格
	4、describe：商品描述
	5、photo：商品照片路径
	6、quantity：商品数量
	*****************************************************/
	userid:=this.Ctx.Input.GetData("user").(int64)
	goodid:=this.Ctx.Input.Param(":id")
	id,_:= strconv.ParseInt(goodid,10,64)
	JsonResponse:=make(map[string]interface{})
	JsonResponse["status"]=0
	good:=models.Goods{Id:id}
	orm:=orm.NewOrm()
	err:=orm.Read(&good)
	if err!=nil {		//无该商品
		this.Data["json"]=JsonResponse
		this.ServeJSON()
		return 
	}
	if good.Userid!=userid {
		this.Data["json"]=JsonResponse
		this.ServeJSON()
		return 
	}
	JsonResponse["status"]=1
	JsonResponse["name"]=good.Name
	JsonResponse["price"]=good.Price
	JsonResponse["describe"]=good.Describe
	JsonResponse["photo"]=good.Photo
	JsonResponse["quantity"]=good.Quantity
	this.Data["json"]=JsonResponse
	this.ServeJSON()
}

func (this *EditController) Post() {
	/**********************************************
	修改指定商品信息的函数
	前端传入值:
	1、name：商品名称
	2、price：商品价格
	3、describe：商品描述
	4、photo：商品照片路径
	5、quantity：商品数量
	返回值：
	1、status：0表示编辑失败，1表示编辑成功
	2、msg：失败或者成功的信息
	**********************************************/
	JsonResponse:=make(map[string]interface{})
	JsonResponse["status"]=0
	userid:=this.Ctx.Input.GetData("user").(int64)
	goodid:=this.Ctx.Input.Param(":id")
	id,_:= strconv.ParseInt(goodid,10,64)
	name:=this.GetString("name")
	price,err1:= this.GetFloat("price")
	if err1!=nil {
		JsonResponse["msg"]="输入格式不正确"
		this.Data["json"]=JsonResponse
		this.ServeJSON()
	}
	describe:=this.GetString("describe")
	quantity,err2:= this.GetInt64("quantity")
	if err2!=nil {
		JsonResponse["msg"]="输入格式不正确"
		this.Data["json"]=JsonResponse
		this.ServeJSON()
	}

	photo,_,err:=this.GetFile("photo")
	if err!=nil {
		JsonResponse["status"] = 0
		JsonResponse["msg"] = "图片上传失败"
		this.Data["json"] = JsonResponse
		this.ServeJSON()
		return
	}
	defer photo.Close()
	orm:=orm.NewOrm()
	good:=models.Goods{Id:id}
	err=orm.Read(&good)
	if err!=nil {
		JsonResponse["msg"]="商品不存在"
		this.Data["json"]=JsonResponse
		this.ServeJSON()
		return
	}
	if good.Userid!=userid {
		JsonResponse["msg"]="没有权限编辑该商品"
		this.Data["json"]=JsonResponse
		this.ServeJSON()
		return 
	}
	//删除旧照片
	filename:=strconv.FormatInt(userid,10)+"_"+goodid+"_"+good.Name+".jpg"
	filename=path.Join("static/photo", filename)
	os.Remove(filename)
	//更新商品信息
	filename=strconv.FormatInt(userid,10)+"_"+goodid+"_"+name+".jpg"
	filename=path.Join("static/photo", filename)
	good.Name=name
	good.Price=price
	good.Describe=describe
	good.Photo=filename
	good.Quantity=quantity
	_,err=orm.Update(&good)
	if err!=nil {
		JsonResponse["msg"]="商品不存在"
	} else {
		JsonResponse["msg"]="商品更新成功"
		this.SaveToFile("photo",filename)
	}
	this.Data["json"]=JsonResponse
	this.ServeJSON()
}

func (this *EditController) Delete() {
	/************************************************
	用于删除商品的函数
	前端传入：无
	后端返回：
	1、status:0/1 0表示删除失败 1表示删除成功
	2、msg：表示失败或成功的消息
	************************************************/
	userid:=this.Ctx.Input.GetData("user").(int64)
	goodid:=this.Ctx.Input.Param(":id")
	id,_:= strconv.ParseInt(goodid,10,64)
	JsonResponse:=make(map[string]interface{})
	orm:=orm.NewOrm()
	good:=models.Goods{Id:id}
	err:=orm.Read(&good)
	if err!=nil {
		JsonResponse["status"]=0
		JsonResponse["msg"]="该商品不存在"
	} else {
		if good.Userid!=userid {
			JsonResponse["msg"]="没有权限删除该商品"
			this.Data["json"]=JsonResponse
			this.ServeJSON()
			return 
		}
		filename:=strconv.Itoa(int(userid))+"_"+goodid+"_"+good.Name+".jpg"
		photopath:=path.Join("static/photo",filename)
		_,err=orm.Delete(&good)
		if err!=nil {
			JsonResponse["status"]=0
			JsonResponse["msg"]="商品删除失败"
		} else {
			os.Remove(photopath)
			JsonResponse["status"]=1
			JsonResponse["msg"]="商品删除成功"
		}
	}
	this.Data["json"]=JsonResponse
	this.ServeJSON()
}

var RestfulHandler=func(ctx *context.Context) {
	/*********************************************
	处理前端请求方法的过滤器，使得beego的路由可以匹配到
	除PUT和POST之外的方法
	前端传入：前端需要以POST方法发起请求，但是携带一个
	"_method"的表单字段表明实际的操作，如：
	"_method":"DELETE"
	后端返回：
	1、status：0表示_method表示的方法是非法的方法
	当_method表示的方法是正常方法时不返回status
	*********************************************/
	RequestMethod:=ctx.Input.Query("_method")
	if RequestMethod=="" {	//正常请求 PUT或者POST
		RequestMethod=ctx.Input.Method()
	}
	var SupportMethod=[6]string{"GET","POST","PUT","PATCH","DELETE","OPTIONS"}
	IsSupport:=false
	for _,method:=range SupportMethod {
		if method==RequestMethod {
			IsSupport=true
			break
		}
	}
	if IsSupport==false {	//不是本应用支持的请求方法
		ctx.ResponseWriter.WriteHeader(405)
		response:=make(map[string]int)
		response["status"]=0
		ctx.Output.JSON(response,true,true) 
	} else {
		//伪造请求方法
		ctx.Request.Method=RequestMethod
	}
	beego.Info(RequestMethod)
}
