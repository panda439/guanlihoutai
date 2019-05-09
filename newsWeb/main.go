package main

import (
	_ "newsWeb/controllers"
	_ "newsWeb/routers"
	"github.com/astaxie/beego"
	_"newsWeb/models"
)

func main() {
	//给视图函数建立映射
	beego.AddFuncMap("prePage",PrePage)
	beego.AddFuncMap("nextPage",NextPage)
	beego.Run()
}

func PrePage(pageNum int) int {
	if pageNum <= 1 {
		beego.Info("当前页数为",pageNum)
		return 1
	}
	beego.Info("当前页数为",pageNum)
	return pageNum - 1
}

func NextPage(pageNum int,pageCount float64) int {
	if pageNum>=int(pageCount) {
		beego.Info("当前页数为",pageNum)
		return int(pageCount)
	}
	beego.Info("当前页数为",pageNum)
	return pageNum + 1
}

