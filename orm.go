
/*这是关于orm的操作实例,在初期时，数据库创建成功，但是数据表创建失败，在其加入orm.RunSyncdb("default",false,true)
即可成功创建数据表。还有就是关于操作的默认字段为主建，可以进行相应的字段修改，进行下一步操作。
*/
package main
import (
	"github.com/astaxie/beego/orm"
	"fmt"
	_"github.com/mattn/go-sqlite3"
)
type User struct{
	Id int64 
	Name string
	Number int64
}
func init(){
	
	orm.RegisterDriver("sqlite3",orm.DRSqlite)
	orm.RegisterDataBase("default","sqlite3","base.db?charset=utf8&loc=Asia%2FShanghai")
	orm.RegisterModel(new(User))	
	orm.RunSyncdb("default",false,true)
	// 数据表(user)不存在，就进行创建数据表(user)，存在则跳过。
}
 
func main(){
		// insertUser("Younger",1206)
		// insertUser("ran",1202)
		// err:=getUser("Younger")
		// if err!=nil{
		// 	fmt.Println("不存在这个子段")
		// 	insertUser("Younger",1206)
		// }
		getUser("tremble")
		updateUser("tremble",1209)
}
func insertUser( name string,number int64 ) error{
	o:=orm.NewOrm()
	user:=&User{
		Name:name,
		Number:number,
	}
	_,err:=o.Insert(user)
	return err
}
func getUser(name string) error{
	o:=orm.NewOrm()
	user:=&User{
		Name:name,
	}
	err:=o.Read(user,"Name")
	// Read查询默认是主建查询，可以进行更改子段进行查询。
	if err == orm.ErrNoRows {
    fmt.Println("查询不到")
	} else if err == orm.ErrMissPK {
    fmt.Println("找不到主键")
	} else {
	 fmt.Println("查找这个子段，终于找到了")	
    fmt.Println("%#v",user)
	}
	return err
}
func updateUser(name string,number int64) error{
	o:=orm.NewOrm()
	user:=&User{
		Name:name,
	}
	if o.Read(user,"name") == nil {
		// 根据这个子段找到记录，
    	user.Number=number
	fmt.Println("Number is ",user.Number)
	// 在上面的一步基础上更新另一个子段。
    if num, err := o.Update(user,"Number"); err == nil {
        fmt.Println(num)
    	}
	}
	// _,err:=o.Update(user,"Number")
	fmt.Printf("%#v",user)
	return nil
}
