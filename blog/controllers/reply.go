package controllers

import (
	"github.com/astaxie/beego"
	"blog/models"
	"fmt"
)
type ReplyController struct{
	beego.Controller
}
func  (this *ReplyController) Add()  {
	tid:=this.Input().Get("tid")
	name:=this.Input().Get("nickname")
	content:=this.Input().Get("content")
	err:=models.AddReply(tid,name,content)
	if err!=nil{
		beego.Error(err)
		return 
	}
	this.Redirect("/topic/cat/"+tid, 301)
	return 
}
func (this *ReplyController) Delete()  {
	  fmt.Println("欢迎进入删除评论")
	if !checkAccount(this.Ctx){
		this.Redirect("/login", 302)
		return 
	}
	rid:=this.Input().Get("rid")
	tid:=this.Input().Get("tid")
	// fmt.Println("rid is ",rid )
	err:=models.DeleteReply(rid)
	if err!=nil{
		beego.Error(err)
		return 
	}
	this.Redirect("/topic/cat/"+tid,301)
//	this.Redirect("/", 301)
	return 
}