package controllers

import(
	"github.com/astaxie/beego"
 	"github.com/astaxie/beego/context"
//	 "net/http"
	// "net/url"
	 "fmt"
)

type LoginController struct{
	beego.Controller
}
func (this *LoginController)Get()  {

    this.Data["IsLogin"]=true
	isExit:=this.Input().Get("exit")
	fmt.Println("isExit",isExit)
	if isExit=="true"{
		//??cookie,-1????
		this.Ctx.SetCookie("uname", "uname",-1,"/")//??beego??cookie???
		this.Ctx.SetCookie("pwd","pwd",-1,"/")
		this.Data["IsLogin"]=false
		fmt.Println("跳转呀!")
		this.TplName="home.html"
		return //?????	
	}
	this.TplName="login.html"
	return 
}

//post 登录
func (this *LoginController)Post()  {
//	this.Ctx.WriteString(fmt.Sprint(this.Input()))
	uname:=this.Input().Get("uname")
	pwd:=this.Input().Get("pwd")
	autoLogin:=this.Input().Get("autoLogin")=="on"
	// var w http.ResponseWriter
	// var r *http.Request
	if beego.AppConfig.String("uname")==uname &&
		beego.AppConfig.String("pwd")==pwd{
				maxAge :=0//???????
			if autoLogin{
				maxAge =1 <<31-1
			}
			//??Cookie
			this.Ctx.SetCookie("uname",uname, maxAge,"/")//??beego??cookie???
			this.Ctx.SetCookie("pwd",pwd,maxAge,"/")
			this.Redirect("/", 301)//???
			return 
		}else{
			this.Redirect("/login",301)
			return 
		} 
	
}
//cookie ???????????
func checkAccount(ctx *context.Context)  bool{
	ck,err:=ctx.Request.Cookie("uname")
	if err!=nil{
		return false
	}
	uname:=ck.Value
	ck,err =ctx.Request.Cookie("pwd")
	if err!=nil{
		return false
	}
	pwd:=ck.Value
	return beego.AppConfig.String("uname")==uname &&
		beego.AppConfig.String("pwd")==pwd
}