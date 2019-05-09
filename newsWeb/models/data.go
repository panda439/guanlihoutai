package models

import (
	"github.com/astaxie/beego/orm"
	_"github.com/go-sql-driver/mysql"
	"time"
)

type User struct {
	Id int
	Name string
	Pwd string
	Articles []*Article `orm:"reverse(many)"`
}
type Article struct {
	Id int `orm:"pk;auto"`
	Title string `orm:"size(40),unique"`
	Content string `orm:"size(500)"`
	Img string  `orm:"null"`
	Time time.Time `orm:"type(datetime);auto_now_add"`
	ReadCount int `orm:"default(0)"`
	//外键
	ArticleType *ArticleType	`orm:"rel(fk);null;on_delete(set_null)"`
	Users []*User `orm:"rel(m2m)"`
}

type ArticleType struct {
	Id int
	TypeName string		`orm:"unique"`
	Articles []*Article	`orm:"reverse(many)"`
}

func init() {



	//注册数据库
	orm.RegisterDataBase("default","mysql","root:lsr88778877@tcp(127.0.0.1:3306)/newWeb")
	//注册表
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	//跑起来
	orm.RunSyncdb("default",false,true)
}


