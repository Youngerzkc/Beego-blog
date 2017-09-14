package controllers

import(
	"github.com/astaxie/beego"
	"blog/models"
	"fmt"
)
type CategoryController struct{
	beego.Controller
}
func (this *CategoryController) Get() {

	this.Data["IsLogin"]=checkAccount(this.Ctx)
	this.Data["IsCategory"]=true
	 fmt.Println("分类列表get方法")
	 op:=this.Input().Get("op")
	 fmt.Println("op is ",op)
	switch op {
	case "add":
	name:=this.Input().Get("name")
	fmt.Println("name is ",name)
	if len(name)==0{
		break
	}
	err:=models.AddCategory(name)
	if err!=nil{
		beego.Error(err)
	}
	// this.Redirect("/category", 301)
	//	 return 
	case "del":	
		 id:=this.Input().Get("id")
		 if len(id)==0{
	 		break
	 	  }
		err:=models.DelCategory(id)
		if err!=nil{
			beego.Error(err)
		}
	//	this.Redirect("/category", 301)
	//	return 	
}
	fmt.Println("执行到这了")	
	var err error
	
	this.Data["Categories"],err=models.GetAllCategory()
	if err!=nil{
		beego.Error(err)
	}
	this.TplName= "category.html"
}