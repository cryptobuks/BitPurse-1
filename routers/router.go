package routers

import (
	"git.coding.net/zhouhuangjing/BitPurse/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})

	beego.Router("/addresses/bitcoin", &controllers.BitcoinController{})
	beego.Router("/addresses/ethereum", &controllers.EthereumController{})
}
