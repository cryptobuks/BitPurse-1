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

// todo: 多重签名, 用于冷热钱包的流转, 比如设定一个比例, 比如说: 90%(进配置)的
// 钱进冷钱包, 这个是需要多重签名才能搞定的,
// 其次热钱包应该定一个定时器, 定期将所有用户的代币存到热钱包, 如果热钱的比例超过设定的比例
// 则将剩余的部分打到冷钱包
// 出于安全的考虑, 多重签名就指接收签名后的结果, 不接受私钥请求
// 多重签名流程
// 1 多个签名生成一个多重签名的redeem script帐号
// 2 交易的时候, 输入则为这个redeem script , 输出参照普通交易
// 3 生成交易数据后, 依次用两个签名来签名,
// 4 两个签名结束后, 等于已经解锁, 可以发布到网络中去了

// @router /tokens/:_token/tx/sign [post]
func (tc *TokenController) SignTx(_token enums.TOKEN) {
	tx := tc.GetString("tx")
	signedTx, complete := service.SignTx(_token, tx)

	type Result struct {
		Tx       string `json:"tx"`
		Complete bool   `json:"complete"`
	}

	r := Result{Tx: signedTx, Complete: complete}
	tc.Data["json"] = &r
	tc.ServeJSON()
}

// @router /tokens/:_token/tx/send [post]
func (tc *TokenController) SendTx(_token enums.TOKEN) {
	tx := tc.GetString("tx")
	signedTx := service.SendTx(_token, tx)

	type Result struct {
		Tx string `json:"tx"`
	}

	r := Result{Tx: signedTx}
	tc.Data["json"] = &r
	tc.ServeJSON()
}

// 获取冷钱包的余额和热钱包的余额, 如果热钱包的余额小于配置的比例
// 则生成一笔交易, 把需要的数量转到热钱包去
// 冷热钱包的地址 由配置指定, 这个问题不大
// 首先手动生成一个 multiple signature 地址 存到配置
// 由于是多重签名地址, 所有权其实不归服务器账户所有
// 是由签名方共同所有, 故没法直接读取
// 得把所有已经确认的交易统计起来, 或者直接去链上查询

// @router /tokens/:_token/cold2hot/new [post]
func (tc *TokenController) NewCold2HotTx(_token enums.TOKEN) {
	tx := service.NewCold2HotTx(_token)

	type Result struct {
		Tx string `json:"tx"`
	}

	r := Result{Tx: tx}
	tc.Data["json"] = &r
	tc.ServeJSON()
}
