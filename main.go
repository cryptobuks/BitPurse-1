package main

import (
	_ "git.coding.net/zhouhuangjing/BitPurse/routers"

	"github.com/astaxie/beego"
	"git.coding.net/zhouhuangjing/BitPurse/models/dao"
	"git.coding.net/zhouhuangjing/BitPurse/models/service"
	"git.coding.net/zhouhuangjing/BitPurse/models/cache"
	"git.coding.net/zhouhuangjing/BitPurse/models/rpc"
)

func main() {

	dao.Init()
	service.Init()
	cache.Init()
	rpc.Init()

	beego.Run()
}
