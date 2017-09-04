package main

import (
	_ "git.coding.net/zhouhuangjing/BitPurse/models/cache"
	"git.coding.net/zhouhuangjing/BitPurse/models/dao"
	"git.coding.net/zhouhuangjing/BitPurse/models/service"
	_ "git.coding.net/zhouhuangjing/BitPurse/routers"
	"github.com/astaxie/beego"
)

func main() {
	if err := dao.Init(); err == nil {
		service.InitBitcoin()
		//service.CheckUpdate()

		beego.Run()
	}
}
