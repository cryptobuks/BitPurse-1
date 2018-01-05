package routers

import (
	"github.com/Newtrong/BitPurse/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	beego.Include(&controllers.TokenController{})
}
