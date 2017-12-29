package main

import (
	_ "./routers"

	"github.com/astaxie/beego"
	"./models/dao"
	"./models/service"
	"./models/cache"
	"./models/rpc"
)

func main() {

	dao.Init()
	service.Init()
	cache.Init()
	rpc.Init()

	beego.Run()
}
