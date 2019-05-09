package routers

import (
	"newsWeb/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {

	//路由过滤器	第一个参数是过滤匹配支持正则，过滤位置， 过滤操作（函数）参数是context
	beego.InsertFilter("/article/*",beego.BeforeExec,fileterFunc)

    beego.Router("/", &controllers.MainController{})
    beego.Router("/register",&controllers.UserController{},"get:ShowRegister;post:HandleRegister")
    beego.Router("/login",&controllers.UserController{},"get:ShowLogin;post:HandleLogin")
    //首页展示
    beego.Router("/article/index",&controllers.ArticleController{},"get,post:ShowIndex")
    //添加文章业务
    beego.Router("/article/addArticle",&controllers.ArticleController{},"get:ShowAddArticle;post:HandleAddArticle")
    //查看文章详情
    beego.Router("/article/content",&controllers.ArticleController{},"get:ShowContent")
    //编辑文章
	beego.Router("/article/update",&controllers.ArticleController{},"get:ShowUpdate;post:HandleUpdate")
    //
	beego.Router("/article/delete",&controllers.ArticleController{},"get:HandleDelete")

	beego.Router("/article/addType",&controllers.ArticleController{},"get:ShowAddType;post:HandleAddType")

    beego.Router("/article/deleteType",&controllers.ArticleController{},"get:HandleDeleteType")

    //退出登陆业务
	beego.Router("/article/logout",&controllers.UserController{},"get:Logout")

}

func fileterFunc(ctx *context.Context)  {
	//登陆校验
	userName := ctx.Input.Session("userName")
	if userName == nil {
		ctx.Redirect(302,"/login")
	}
}
