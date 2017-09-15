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
	cate:=new(Category)
	cate.Title=category
	err=o.Read(cate,"Title")
	if err==nil{
		fmt.Println("查到这个分类了")
	}else{
		// 查询不到分类，插入分类。
		AddCategory(category)
		o.Read(cate,"Title")
	}
	cate.TopicCount++
	_,err=o.Update(cate)
	return err
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
	cate:=&Category{
		Title:category,
	}
	errCate:=o.Read(cate,"title")
	topic:=&Topic{Id:tidNum}
	if o.Read(topic)==nil{
		// 对应的文章标题，内容进行修改了。若只修改了分类呢？
		if topic.Category!=category{
			// 只改修改分类的情况处理！
			// 1.若该分类不存在 2.该文类已存在
			if errCate!=nil{
				AddCategory(category)
				o.Read(cate,"title")
				cate.TopicCount++
			}else{
				cate.TopicCount++
			}
			// 原分类文章数减1
				o.Update(cate)
				cate.Title=topic.Category
				o.Read(cate,"title")
				cate.TopicCount--
		}else if topic.Title!=title||topic.Content!=content && topic.Category==category{
			// 类别不变，文章标题和内容改变！
			cate.TopicCount++
		}else  {
			// 分章未进行修改
		}
		topic.Title=title
		topic.Content=content
		topic.Id=tidNum
		topic.Category=category
		topic.Updated=time.Now()
		o.Update(topic)
	}
	_,err=o.Update(cate)
	return err
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
	cate:=new(Category)
	cate.Title=topic.Category
	err=o.Read(cate,"Title")
	if err==nil{
		cate.TopicCount--
		o.Update(cate)
	}else{
		// 不存在这个分类
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
func AddReply(tid,name,content string) error  {	 
	 o:=orm.NewOrm()
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
	_,err=o.Insert(reply)
	if err==nil{
	  err=o.Read(reply,"Tid")
		if err==nil{
			fmt.Println("找到这篇文章了")
		reply.CommentCount++;
		}
		_,err=o.Update(reply)
	}
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
	reply:=&Comment{Id:ridNum}
	o.Read(reply,"id")
	reply.CommentCount--
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
