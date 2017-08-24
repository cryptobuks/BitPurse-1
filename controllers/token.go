package controllers

import (
	"git.coding.net/zhouhuangjing/BitPurse/models/common/enums"
	"git.coding.net/zhouhuangjing/BitPurse/models/common/types"
	"git.coding.net/zhouhuangjing/BitPurse/models/service"
	"github.com/astaxie/beego"
)

type TokenController struct {
	beego.Controller
}

// @router /users/:userId/tokens/:token/deposit [get]
func (bc *TokenController) Deposit(userId types.ID, token enums.TOKEN) {
	s := service.Get(token)

	ut := service.Deposit(s, userId)

	type Result struct {
		Address string `json:"address"`
	}

	r := Result{Address: ut.TokenAddress}
	bc.Data["json"] = &r
	bc.ServeJSON()
}

// post withdraw address id, amount
// @router /users/:_userId/tokens/:_token/withdraw [post]
func (bc *TokenController) Withdraw(_userId types.ID, _token enums.TOKEN) {
	s := service.Get(_token)

	addr, err1 := bc.GetInt("address")
	if err1 != nil {
		beego.Error(err1)
		bc.Finish()
	}

	amount, err2 := bc.GetFloat("amount")
	if err2 != nil {
		beego.Error(err2)
		bc.Finish()
	}

	ut := service.Withdraw(s, _userId, types.ID(addr), amount)

	type Result struct {
		Address string `json:"address"`
	}

	r := Result{Address: ut.TokenAddress}
	bc.Data["json"] = &r
	bc.ServeJSON()
}
