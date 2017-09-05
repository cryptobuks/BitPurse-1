package main

import (
	_ "git.coding.net/zhouhuangjing/BitPurse/models/cache"
	_ "git.coding.net/zhouhuangjing/BitPurse/models/dao"
	_ "git.coding.net/zhouhuangjing/BitPurse/models/service"
	_ "git.coding.net/zhouhuangjing/BitPurse/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
