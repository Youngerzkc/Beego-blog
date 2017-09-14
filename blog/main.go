package main

import (
			_ "blog/routers"
			"github.com/astaxie/beego"
			"github.com/astaxie/beego/orm"
			"blog/models"
			"blog/controllers"
)
func init(){
	models.RegisterDB()
}
func main() {
	//默认情况下不会建表,
	orm.Debug=true //打印所有调试信息
	orm.RunSyncdb("default", false,false)//数据库自动建表
	// 1.默认2.true会删掉重新建表，改为false，3.true，打印建表所有信息
	
	beego.Router("/",&controllers.HomeController{})
	beego.Router("/login",&controllers.LoginController{})
	
	beego.Router("/cate",&controllers.CategoryController{})
	
	beego.Router("/topic", &controllers.TopicController{})
	// beego自动路由
	beego.AutoRouter(&controllers.TopicController{})
	//注意格式
	 beego.Router("/reply/add",&controllers.ReplyController{},"post:Add")
	// beego.Router("/reply/delete",&controllers.ReplyController{})
	 beego.Router("/reply/delete",&controllers.ReplyController{},"get:Delete")
	beego.Run()

}

