package main

import (
	_ "github.com/Newtrong/BitPurse/routers"

	"github.com/astaxie/beego"
	"github.com/Newtrong/BitPurse/models/dao"
	"github.com/Newtrong/BitPurse/models/service"
	"github.com/Newtrong/BitPurse/models/cache"
	"github.com/Newtrong/BitPurse/models/rpc"
)

func main() {

	dao.Init()
	service.Init()
	cache.Init()
	rpc.Init()

	beego.Run()
}
