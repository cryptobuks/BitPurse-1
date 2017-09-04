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

// TODO: 异步操作, 多重签名,  监控代币充值状态并更新记录, 以太坊, 矿工费比例(数据库),
// TODO: 提币
// @router /users/:_userId/tokens/:_token/deposit [post]
func (tc *TokenController) Deposit(_userId types.ID, _token enums.TOKEN) {

	ut := service.Deposit(_token, _userId)

	type Result struct {
		Address string `json:"address"`
	}

	r := Result{Address: ut.TokenAddress}
	tc.Data["json"] = &r
	tc.ServeJSON()
}

// post withdraw address id, amount
// @router /users/:_userId/tokens/:_token/withdraw [post]
func (tc *TokenController) Withdraw(_userId types.ID, _token enums.TOKEN) {

	addr, err1 := tc.GetInt("address")
	if err1 != nil {
		beego.Error(err1)
		tc.Finish()
	}

	amount, err2 := tc.GetFloat("amount")
	if err2 != nil {
		beego.Error(err2)
		tc.Finish()
	}

	w := service.Withdraw(_token, _userId, types.ID(addr), amount)

	type Result struct {
		Address string `json:"address"`
	}

	r := Result{Address: w.Address}
	tc.Data["json"] = &r
	tc.ServeJSON()
}

// @router /users/:_userId/tokens/:_token/withdrawal [post]
func (tc *TokenController) NewWithdrawal() {

}

// post withdraw address id, amount
// @router /tokens/:_token/tx/:_txId/notify [get]
func (tc *TokenController) WatchNotify(_token enums.TOKEN, _txId string) {
	tr := service.WalletNotify(_token, _txId)
	if tr == nil {
		beego.Error(_token, _txId)
		tc.Abort("405")
	} else {
		type Result struct {
			Address types.ID `json:"address"`
		}

		r := Result{Address: tr.Id}
		tc.Data["json"] = &r
		tc.ServeJSON()
	}
}

func (tc *TokenController) GenMultiSig(_token enums.TOKEN) {

}
