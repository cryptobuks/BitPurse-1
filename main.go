package main

import (
	_ "github.com/forchain/BitPurse/routers"

	"github.com/astaxie/beego"
	"github.com/forchain/BitPurse/models/dao"
	"github.com/forchain/BitPurse/models/service"
	"github.com/forchain/BitPurse/models/cache"
	"github.com/forchain/BitPurse/models/rpc"
)

func main() {

	dao.Init()
	service.Init()
	cache.Init()
	rpc.Init()

	beego.Run()
}
