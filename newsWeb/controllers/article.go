package controllers

import (
	"github.com/astaxie/beego"
	"path"
	"time"
	"github.com/astaxie/beego/orm"
	"newsWeb/models"
	"math"
	"strconv"
	"github.com/gomodule/redigo/redis"
	"encoding/gob"
	"bytes"
)

type ArticleController struct {
	beego.Controller
}
//主页
func (this *ArticleController)ShowIndex()  {
	//获取所有文章数据，展示到页面

	//校验登陆状态
	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/login",302)
		return
	}
	this.Data["userName"] = userName.(string)


	o := orm.NewOrm()
	qs := o.QueryTable("Article")
	var articles []models.Article
	//qs.All(&articles)

	//获取选中的类型
	typeName := this.GetString("select")
	var count int64

	if typeName == "" {
		//获取总记录数
		count,_ = qs.RelatedSel("ArticleType").Count()
	}else {
		count,_ = qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).Count()
	}

	//获取总页数
	pageIndex := 2


	pageCount := math.Ceil(float64(count) / float64(pageIndex))
	//获取首页和末页数据
	//获取页码
	pageNum ,err := this.GetInt("pageNum")
	if err != nil {
		pageNum = 1
	}
	beego.Info("数据总页数未:",pageNum)

	//获取对应页的数据   获取几条数据     起始位置
	//ORM多表查询的时候默认是惰性查询 关联查询之后，如果关联的字段为空，数据查询不到


	//where ArticleType.typeName = typename   filter相当于where
	if typeName == "" {
		qs.Limit(pageIndex,pageIndex * (pageNum - 1)).RelatedSel("ArticleType").All(&articles)
	}else {
		qs.Limit(pageIndex,pageIndex * (pageNum - 1)).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)
	}


	//查询所有文章类型，并展示
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)
	this.Data["articleTypes"] = articleTypes

	//把数据存储到redis中
	conn , err := redis.Dial("tcp",":6379")
	if err != nil {
		beego.Error("redis数据库连接失败")
		return
	}
	defer  conn.Close()

	resp,err := conn.Do("get","newWeb")
	result,_ := redis.Bytes(resp,err)
	if len(result) == 0 {
		o.QueryTable("ArticleType").All(&articleTypes)
		this.Data["articleTypes"] = articleTypes
		//把数据存储到redis中	转码处理
		var buffer bytes.Buffer//容器
		enc := gob.NewEncoder(&buffer)//编码器
		_ = enc.Encode(articleTypes)//编码

		conn.Do("set","newsWeb",buffer.Bytes())
		resp,_ := conn.Do("get","newWeb")
		result,_ := redis.Bytes(resp,err)
		dec := gob.NewDecoder(bytes.NewReader(result))
		dec.Decode(articleTypes)
		beego.Info("获取到的类型resp",articleTypes)
	}else {
		//解码操作
		dec := gob.NewDecoder(bytes.NewReader(result))
		//定义一个容器，接收解码之后的数据

		dec.Decode(articleTypes)
	}





	this.Data["articles"] = articles
	this.Data["count"] = count
	this.Data["pageCount"] = pageCount
	this.Data["pageNum"] = pageNum

	//将选中的Type类型传递给前端
	this.Data["TypeName"] = typeName

	//
	//this.Layout = "layout.html"
	//this.LayoutSections = make(map[string]string)
	//this.LayoutSections["indexJs"] = "indexJs.html"

	this.TplName = "index.html"

}





//展示添加文章页面
func (this *ArticleController)ShowAddArticle()  {
	//获取所有类型并绑定下拉框
	o := orm.NewOrm()
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)

	this.Data["ArticleTypes"] = articleTypes
	this.TplName = "add.html"
}

func (this *ArticleController)HandleAddArticle()  {
	//获取数据
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	typeName := this.GetString("select")

	//校验数据
	if articleName == "" || content == "" {
		beego.Error("获取数据错误")
		this.Data["errmsg"] = "获取数据错误"
		this.TplName = "add.html"
		return
	}

	//获取图片
	//返回值 文件二进制流  文件头    错误信息
	file,head,err := this.GetFile("uploadname")
	if err != nil {
		beego.Error("获取数据错误")
		this.Data["errmsg"] = "图片上传失败"
		this.TplName = "add.html"
		return
	}
	defer file.Close()
	//校验文件大小
	if head.Size >5000000{
		beego.Error("获取数据错误")
		this.Data["errmsg"] = "图片数据过大"
		this.TplName = "add.html"
		return
	}

	//校验格式 获取文件后缀
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		beego.Error("获取数据错误")
		this.Data["errmsg"] = "上传文件格式错误"
		this.TplName = "add.html"
		return
	}

	//防止重名
	fileName := time.Now().Format("200601021504052222")


	//jianhuangcaozuo

	//把上传的文件存储到项目文件夹
	this.SaveToFile("uploadname","./static/img/"+fileName+ext)

	//处理数据
	//把数据存储到数据库
	//获取orm对象
	o := orm.NewOrm()
	//获取插入独享
	var article models.Article
	//给插入对象赋值
	article.Title = articleName
	article.Content = content
	article.Img = "/static/img/"+fileName+ext

	//获取一个类型对象，并插入到文章中
	var articleType models.ArticleType
	articleType.TypeName = typeName
	o.Read(&articleType,"TypeName")

	beego.Info(articleType)
	article.ArticleType = &articleType

	//插入数据
	_,err = o.Insert(&article)
	if err != nil {
		beego.Error("获取数据错误",err)
		this.Data["errmsg"] = "数据插入失败"
		this.TplName = "add.html"
		return
	}

	//返回数据  跳转页面
	this.Redirect("/article/index",302)

}

func (this *ArticleController)ShowContent(){
	//获取数据
	id,err := this.GetInt("id")
	if err != nil {
		beego.Error("获取文章id错误")
		this.Redirect("/article/index",302) 		//渲染 如果页面本身有数据加载不能直接渲染
		return
	}
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	o.Read(&article)


	//高级查询	首先指定表	多对多查询二
	var users []models.User
	o.QueryTable("User").Filter("Articles__Article__Id",id).Distinct().All(&users)
	this.Data["users"] = users


	//给更新条件赋值
	article.ReadCount += 1
	o.Update(&article)

	//返回数据
	this.Data["article"] = article

	//插入多对多关系，根据用户名获取用户对象
	userName := this.GetSession("userName")
	var user models.User
	user.Name = userName.(string)
	o.Read(&user,"Name")

	//多对多的插入操作
	//获取ORM对象
	//获取被插入数据的对象 文章

	//获取多对多操作对象
	m2m := o.QueryM2M(&article,"Users")
	//用多对多操作对象
	m2m.Add(user)


	this.TplName = "content.html"
}

func (this *ArticleController)ShowUpdate(){
	//获取数据
	id,err := this.GetInt("id")
	if err != nil {
		beego.Error("获取文章id错误")
		this.Redirect("/article/index",302) 		//渲染 如果页面本身有数据加载不能直接渲染
		return
	}
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	o.Read(&article)

	//返回数据
	this.Data["article"] = article


	this.TplName = "update.html"
}


//封装上传文件处理函数
func UploadFile(this *ArticleController,fileName string,errHtml string) string {

	//文件二进制流 文件头 错误信息
	file,head,err := this.GetFile(fileName)
	if err != nil {
		beego.Error("获取文件错误")
		this.Data["errmsg"] = "图片上传失败"
		this.TplName = errHtml
		return""
	}
	defer file.Close()
	//校验大小
	if head.Size>5000000 {
		beego.Error("获取数据失败")
		this.Data["errmsg"] = "图片数据过大"
		this.TplName = errHtml
		return""
	}

	//格式
	ext := path.Ext(head.Filename)
	beego.Info(ext)
	if ext != ".jpg" && ext != ".png" && ext!="jpeg" {
		beego.Error("获取数据失败")
		this.Data["errmsg"] = "上传文件格式错误"
		this.TplName = errHtml
		return ""
	}

	//防止重名
	filename:=time.Now().Format("200601021504052222")

	//把上传的文件存储到某文件夹
	this.SaveToFile(fileName,"./static/img/"+filename+ext)
	o := orm.NewOrm()
	var article models.Article
	article.Img = "/static/img/"+filename+ext

	_,err = o.Insert(&article)
	if err != nil {
		beego.Error("插入失败",err)
		this.Data["errmsg"] = "数据插入数据库失败"
		this.TplName = errHtml
		return ""
	}
	return article.Img
}

func (this *ArticleController)HandleUpdate(){
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	savePath := UploadFile(this,"uploadname","update.html")
	id,_:= this.GetInt("id")
	if articleName == ""||content == "" || savePath == "" {
		beego.Error("获取数据失败")
		this.Data["errmsg"] = "获取数据错误"
		this.Redirect("/article/update?id="+strconv.Itoa(id),302)
		return
	}

	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	o.Read(&article)
	//更新 需要先赋值 beego中的ORM如果需要更新，Id必须有值
	article.Title = articleName
	article.Content = content
	article.Img = savePath
	o.Update(&article)

	//返回数据
	this.Redirect("/article/index",302)

}

func (this *ArticleController)HandleDelete(){
	id,err := this.GetInt("id")
	if err != nil {
		beego.Error("获取id错误")
		this.Redirect("/article/index",302)
		return
	}
	//处理数据
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	o.Delete(&article,"Id")
	//返回数据
	this.Redirect("/article/index",302)
}

//展示添加分类页面
func (this *ArticleController)ShowAddType(){
	//获取所有
	o := orm.NewOrm()
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)


	this.Data["articleTypes"] = articleTypes

	this.TplName = "addType.html"
}

func (this *ArticleController)HandleAddType(){

	typeName:=this.GetString("typeName")
	if typeName == ""{
		beego.Error("类型名称传输失败")
		this.Redirect("/article/addType",302)
		return
	}
	//处理数据
	//插入操作
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName = typeName
	o.Insert(&articleType)

	//
	this.Redirect("/article/addType",302)


}

func (this *ArticleController)HandleDeleteType(){
	id,err := this.GetInt("Id")
	if err != nil {
		beego.Error("获取ID失败")
		this.Redirect("/article/addType",302)
		return
	}
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.Id = id
	o.Delete(&articleType,"Id")
	this.Redirect("/article/addType",302)

}


