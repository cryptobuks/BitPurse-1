package routers

import (
	"github.com/forchain/BitPurse/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	beego.Include(&controllers.TokenController{})
}
