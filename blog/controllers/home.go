package controllers

import (
	"fmt"
	"blog/models"
	"github.com/astaxie/beego"

)

type HomeController struct {
	beego.Controller
}

func (this *HomeController) Get() {
	this.Data["IsHome"]=true
	//this.TplName="login.html"
	 this.Data["IsLogin"]=checkAccount(this.Ctx)
	 fmt.Println(" checkAccount ", checkAccount(this.Ctx))
	if this.Data["IsLogin"]==false{
		this.Redirect("/login",301)
		return 
	}
	fmt.Println(len(this.Input().Get("cate")))
	topics, err:=models.GetAllTopic(this.Input().Get("cate"),false)//???????
	 if err!=nil{
		 beego.Error(err.Error())
	 }else{
		 this.Data["Topics"]=topics
	 }
	 categories,err:=models.GetAllCategory()
	 if err!=nil{
		 beego.Error(err)
	 }
	 //comments,err:=models.GetAllRelies(tid string)
	 this.Data["Categories"]=categories
	 this.TplName = "home.html"
}
