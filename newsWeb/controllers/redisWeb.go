package controllers

import (
	"github.com/gomodule/redigo/redis"
	"github.com/astaxie/beego"
)

func init() {
	//连接函数
	conn ,err := redis.Dial("tcp4",":6379")
	if err != nil {
		beego.Error("redis连接失败")
		return
	}
	defer conn.Close()


	//操作函数
	//conn.Send("set","redis","test")
	//conn.Flush()
	//conn.Receive()
	resp,err := conn.Do("mget","t1","t2","t3")


	//回复助手函数
	result,_ := redis.Values(resp,err)
	//把对应的函数扫描到相应变量里面
	var v1,v2 int
	var v3 string
	redis.Scan(result,&v1,&v3,&v2)
	beego.Info("获取的数据为",v1,v2,v3)

}

