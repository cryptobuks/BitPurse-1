package routers

import (
	"git.coding.net/zhouhuangjing/BitPurse/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	beego.Include(&controllers.TokenController{})
}
