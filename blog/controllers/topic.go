package controllers

import(
	"github.com/astaxie/beego"
	"blog/models"
	"fmt"
)
type TopicController struct{
	beego.Controller
}
func (this *TopicController)Get()  {
	this.Data["IsLogin"]=checkAccount(this.Ctx)
	this.Data["IsTopic"]=true
	 this.TplName="topic.html"
	 topics, err:=models.GetAllTopic("",false)
	 if err!=nil{
		 beego.Error(err.Error())
	 }else{
		 this.Data["Topics"]=topics
	 }
}
func (this *TopicController)Post(){
	if !checkAccount(this.Ctx){
		this.Redirect("/login",302)
		return 
	}
	tid:=this.Input().Get("tid")
	title:=this.Input().Get("title")
	content:=this.Input().Get("content")
	category:=this.Input().Get("category")
	var err error

	if len(tid) ==0{
		err=models.AddTopic(title,category,content)
		if err!=nil{
			beego.Error(err)
		}
		err=models.AddCategory(category)
		if err!=nil{
			beego.Error(err)
		}
	}else{
		err=models.ModifyTopic(tid, title ,category, content )
		if err!=nil{
			beego.Error(err)
		}
	    err=models.AddCategory(category)
		if err!=nil{
			beego.Error(err)
		}
	}
	
	//err=models.AddTopic(title,content)

	this.Redirect("/topic",301)
}
func (this *TopicController)Add(){	

	this.TplName="topic_add.html"

}
//func (this *TopicController)View()
//该方法名匹配，存在
func (this *TopicController) Cat(){
	 this.Data["IsLogin"]=checkAccount(this.Ctx)
	//  fmt.Println("CAT+++++++++++")
	mm:=this.Ctx.Input.Params()
	if mm["0"]==""{
		//fmt.Println("没有参数o")
		this.Redirect("/topic", 301)
		return 
	}
	   for k,_:=range  mm{
		   fmt.Println("k is ",k,mm[k])
	   }
	 topic,err:=models.GetTopic(mm["0"])
	 if err!=nil{
		 beego.Error(err)
		 return 
	 }
	 this.Data["Topic"]=topic
 	 this.Data["Tid"]=mm["0"]//文章ID
	replies,err:=models.GetAllRelies(mm["0"]) //取文章失败?
	if err!=nil{
		beego.Error(err)
		return
	}
	 this.Data["Replies"]=replies 
	 this.TplName="topic_view.html"
}
func  (this *TopicController) Modify()  {
	this.TplName="topic_modify.html"
	tid :=this.Input().Get("tid")
	fmt.Println("hello ++++")
	topic,err:=models.GetTopic(tid)
	if err!=nil{
		beego.Error(err)
		this.Redirect("/", 302)
		return 
	}
	this.Data["Topic"]=topic
	this.Data["Tid"]=tid 
	return 
}
func (this *TopicController) Delete()  {	
	if !checkAccount(this.Ctx){
		this.Redirect("/login",302)
		return 
	}
  tid :=this.Input().Get("tid")

  if len(tid)==0{
	  this.Redirect("/",302)
	  return 
  }
  err:=models.DeleteTopic(tid)
  if err!=nil{
	this.Redirect("/",301)
	return 
   }
   this.Redirect("/topic", 302)
   return 
}