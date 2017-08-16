package controllers

import (
	"github.com/astaxie/beego"
	"git.coding.net/zhouhuangjing/BitPurse/models/tokens"
)

var (
	helper = tokens.BitcoinHelper{}
)

type BitcoinController struct {
	beego.Controller
}

func init() {

}

// new bitcoin address
func (b *BitcoinController) Post() {

	address := helper.NewAddress()

	type Result struct {
		Address string `json:"address"`
	}

	r := Result{Address: address}
	b.Data["json"] = &r
	b.ServeJSON()
}
