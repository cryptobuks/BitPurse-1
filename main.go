package main

import (
	_ "git.coding.net/zhouhuangjing/BitPurse/routers"
	"github.com/astaxie/beego"
	"git.coding.net/zhouhuangjing/BitPurse/models/dao"
	"git.coding.net/zhouhuangjing/BitPurse/models/service"
)

func main() {
	if err := dao.Init(); err == nil {
		service.InitBitcoin()

		beego.Run()
	}
}
