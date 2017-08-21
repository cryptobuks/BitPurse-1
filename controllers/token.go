package controllers

import (
	"github.com/astaxie/beego"
	"git.coding.net/zhouhuangjing/BitPurse/models/rpc"
	"git.coding.net/zhouhuangjing/BitPurse/models/service"
	"git.coding.net/zhouhuangjing/BitPurse/models/common"
)

var (
	rpc_     rpc.IRpc
)

type TokenController struct {
	beego.Controller
}

// new bitcoin address
func (bc *TokenController) NewBitcoinAddress() {

	address := rpc_.NewAddress()

	type Result struct {
		Address string `json:"address"`
	}

	r := Result{Address: string(address)}
	bc.Data["json"] = &r
	bc.ServeJSON()
}

// @router /users/:userId/tokens/:token/deposit [get]
func (bc *TokenController) Deposit(userId common.ID, token common.TOKEN) {
	s := service.Get(token)

	uk := s.Deposit(userId)

	type Result struct {
		Address string `json:"address"`
	}

	r := Result{Address: string(uk.TokenAddress)}
	bc.Data["json"] = &r
	bc.ServeJSON()
}


