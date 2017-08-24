package controllers

import (
	"github.com/astaxie/beego"
	"git.coding.net/zhouhuangjing/BitPurse/models/service"
	"git.coding.net/zhouhuangjing/BitPurse/models/common/types"
	"git.coding.net/zhouhuangjing/BitPurse/models/common/enums"
)

type TokenController struct {
	beego.Controller
}

// @router /users/:userId/tokens/:token/deposit [get]
func (bc *TokenController) Deposit(userId types.ID, token enums.TOKEN) {
	s := service.Get(token)

	ut :=service.Deposit(s, userId)

	type Result struct {
		Address string `json:"address"`
	}

	r := Result{Address: ut.TokenAddress}
	bc.Data["json"] = &r
	bc.ServeJSON()
}

// @router /users/:userId/tokens/:token/withdraw [get]
func (bc *TokenController) Withdraw(userId types.ID, token enums.TOKEN) {
	s := service.Get(token)



}
