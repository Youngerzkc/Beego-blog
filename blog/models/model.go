package models

import(
	// "os"
	"time"
	//	"path"
	"fmt"
	"strconv"
	// "github.com/Unknown/com"
	"github.com/astaxie/beego/orm"
	_"github.com/mattn/go-sqlite3"
)
const (
	_DB_NAME="data/beeblog.db"
	_SQLITE3_DRIVER="sqlite3"
)
type Category struct{
	Id int64  //主键
	Title string   
	//Created time.Time `orm:"index"` //创建时间，反射必须是字段导出
	Views int64 `orm:"index"` //浏览数量，索引，建表的时候建个索引
	// tag 只有orm可以读
	// TopicTime time.Time 
	TopicCount int64
	Category string
	TopicLastUserId int64 //最后一个操作者Id
}
type Topic struct{
	Id int64  
	Uid int64 //z作者
	Title string
	Category string
	Content string `orm:"size(5000)"`
	Attachment string //附件
	Created time.Time`orm:"index"`
	//Created  string `orm:"index"`
	Updated time.Time `orm:"index"`
	Views int64 `orm:"index"`//
	Author string //Uid 链接到作者信息
	ReplyTime time.Time `orm:"index"`//最后回复时间
	ReplyCount int64  //评论个数
	ReplyLastUserId int64
}
type Comment struct{
	Id int64
	Tid int64
	Name string
	Content string `orm:"size(10000)"`
	Created time.Time `orm:"index"`
	CommentCount int64 `orm:"index"`
}
type Account struct{
	Id int64 
	Name string `xorm:"unique"`
	Passwd string 
}
func RegisterDB()  {
	//检查数据库存在
	// os.MkdirAll(path.Dir(_DB_NAME),os.ModePerm)
	// os.Create(_DB_NAME)
	orm.RegisterModel(new(Account),new(Category),new(Topic),new(Comment))//注册模型(struct)
	orm.RegisterDriver(_SQLITE3_DRIVER,orm.DRSqlite)//注册数据库驱动,可有可无
	orm.RegisterDataBase("default",_SQLITE3_DRIVER,_DB_NAME+"?charset=utf8&loc=Asia%2FShanghai",10)//default,规定参数，10是最大数据库链接个数
	
}
func AddCategory(title string) error  {
	o:=orm.NewOrm()
	 cate:=&Category{
		 Title:title,
		}
	// err := o.QueryTable("category").Filter("title", "title").One(cate)
	// // 大小写敏感，注意可另外进行设
	// fmt.Println("err is ",err)	
	// if err == orm.ErrMultiRows {
    // // 多条的时候报错
    // 	fmt.Printf("Returned Multi Rows Not One")
	// 	}else if err == orm.ErrNoRows {
    //       // 没有找到记录
   	// 	 fmt.Printf("Not row found")
	// 	} else {
	// 	fmt.Println("查到了")
	// 	return nil
	// 	}	
	err:=o.Read(cate,"title")	
		// 已有分类
	 if err==orm.ErrNoRows{
		//  查询不到
		fmt.Println("查询不到")
		_,err =o.Insert(cate)
		if err!=nil{
	 	fmt.Println("插入失败")
			return err
		}
	 }else{
		//  查到了
		fmt.Println("查到栏")
		return nil;
	 }	
	 return nil
}
func GetAllCategory()([]*Category,error){	
	o:=orm.NewOrm()
	cates:=make([]*Category,0)
	qs:=o.QueryTable("category")
	fmt.Println("查询数据库")
	_,err:=qs.All(&cates)
	return cates,err
}
func DelCategory(id string ) error  {
	cid ,err :=strconv.ParseInt(id, 10, 64)
	if err!=nil{
		return err
	}
	o:=orm.NewOrm()
	cate:=&Category{Id:cid}
	_,err=o.Delete(cate)
	if err!=nil{
		return err
	}
	return nil
}
func  AddTopic( title ,category,content string ) error  {
	o:=orm.NewOrm()
	topic:=&Topic{
		Title:title,
		Category:category,//后增的
		Content:content,
		Created:time.Now(),
		Updated:time.Now(),
		ReplyTime:time.Now(),
	}
	_,err:=o.Insert(topic)
	if err!=nil{
		fmt.Println("insert topic failed.")
		return err
	}
	return nil
}
func GetTopic(tid string ) (*Topic,error)  {
	tidNum,err:=strconv.ParseInt(tid, 10,64)
	if err!=nil{
		return nil,err
	}
	o:=orm.NewOrm()
	topic:=new(Topic)
	qs:=o.QueryTable("topic")
	err=qs.Filter("id",tidNum).One(topic)
	if err!=nil{
	return nil,err	
		}
	topic.Views++
	_,err=o.Update(topic)
		return topic,err
}
func  GetAllTopic(cate string,isDesc bool ) ([]*Topic,error) {
	o:=orm.NewOrm()
	topics:=make([]*Topic,0)
	qs:=o.QueryTable("topic")
	var err error
	fmt.Println("isDesc+++++",isDesc);
	if isDesc==false {
		if len(cate)>0{
			qs=qs.Filter("category",cate)
		}
		fmt.Println("进入查询了")
		_,err=qs.OrderBy("-created").All(&topics)//注意嵌入字段的方式，可能会不同
		return topics,err
	}else{
	_,err =qs.All(&topics)//不涉及分表问题
	}
	return topics,err
}
func  ModifyTopic(tid , title ,category,content string) error  {
	tidNum ,err:=strconv.ParseInt(tid,10,64)
	if err!=nil{
		return err
	}
	o:=orm.NewOrm()
	// cate:=&Category{Title:category}
	topic:=&Topic{Id:tidNum}
	if o.Read(topic)==nil{
		topic.Title=title
		topic.Content=content
		topic.Id=tidNum
		topic.Category=category
		topic.Updated=time.Now()
		o.Update(topic)
	}
	return nil
}
func DeleteTopic(tid string) error  {
	fmt.Println("删除文章数")
	cid ,err :=strconv.ParseInt(tid, 10, 64)
	if err!=nil{
		return err
	}
	o:=orm.NewOrm()
	topic:=&Topic{Id:cid}
	err =o.Read(topic)
	if err!=nil{
		return err
	}
	fmt.Println("查到文章数椐了")
	_,err=o.Delete(topic)
	if err!=nil{
		fmt.Println("删除文章数椐失败了")
		return err
	}
	return nil
}
var i int64 =0
func AddReply(tid,name,content string) error  {
	 i++;
	tidNum,err:=strconv.ParseInt(tid,10,64)
	if err!=nil{
		return err
	}
	reply:=&Comment{
		Tid:tidNum,
		Name:name,
		Created:time.Now(),
		Content:content,
	}
	reply.CommentCount=i
	o:=orm.NewOrm()
	_,err=o.Insert(reply)
	return err
}
func GetAllRelies(tid string)(replies []*Comment,err error )  {
	tidNum,err:=strconv.ParseInt(tid,10,64)
	if err!=nil{
		return nil,err
	}
	replies=make([]*Comment,0)
	o:=orm.NewOrm()
	qs:=o.QueryTable("comment")
	_,err=qs.Filter("tid", tidNum).All(&replies)
	return replies,err
}
func DeleteReply(rid string) error {
	ridNum,err:=strconv.ParseInt(rid,10,64)
	if err!=nil{
		return err
	}
	o:=orm.NewOrm()
	i--
	reply:=&Comment{Id:ridNum}
	_,err=o.Delete(reply)
	return err
}
func RegisterAccount(name string,pwd string) error{ 
	o:=orm.NewOrm()
	account:=&Account{
		Name:name,
		Passwd:pwd,
	}
	err:=o.Read(account,"Name")
	if err==orm.ErrNoRows{
		fmt.Println("该帐号不存在，可以进行注册！")
		_,err=o.Insert(account)
		if err!=nil{
			fmt.Println("注册失败")
			return err
		}
		return nil
	}else{
		fmt.Println("该帐号已存在，请重新注册")
	}
		return err
}
func CheckLogin(name string,pwd string) error{
	o:=orm.NewOrm()
	account:=&Account{
		Name:name,
		Passwd:pwd,
	}
	err:=o.Read(account,"name")
	if err==orm.ErrNoRows{
		fmt.Println("该帐号不存在，请注册");
		RegisterAccount(name,pwd);
	}else{
		if account.Passwd==pwd{
			fmt.Println("登陆成功")
			return nil
		}
	}
	return nil
}








