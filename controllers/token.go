package controllers

import (
	"github.com/forchain/BitPurse/models/common/enums"
	"github.com/forchain/BitPurse/models/common/types"
	"github.com/forchain/BitPurse/models/service"
	"github.com/astaxie/beego"
	"fmt"
)

type TokenController struct {
	beego.Controller
}

// @router /tokens/:_token/deposit [post]
func (tc *TokenController) Deposit(_token enums.TOKEN) {
	if userID, err := tc.GetInt("userID"); err == nil {
		if ut := service.Deposit(_token, types.ID(userID)); ut != nil {
			type Result struct {
				Address string `json:"address"`
			}

			r := Result{Address: ut.TokenAddress}
			tc.Data["json"] = &r
			tc.ServeJSON()
			return
		}
	}

	tc.CustomAbort(405, fmt.Sprintf("Deposit failed %d", _token))
}

// post withdraw address id, amount
// @router /tokens/:_token/withdraw [post]
func (tc *TokenController) Withdraw(_token enums.TOKEN) {

	addr, err1 := tc.GetInt("address")

	amount, err2 := tc.GetFloat("amount")

	userID, err3 := tc.GetInt("userID")
	if err1 != nil || err2 != nil || err3 != nil {
		beego.Error(err1, err2, err3)
		tc.CustomAbort(405, fmt.Sprintf("Withdraw Invalid %s %s %s", err1, err2, err3))
		return
	}

	if w := service.Withdraw(_token, types.ID(userID), types.ID(addr), amount); w != nil {
		type Result struct {
			Address string `json:"address"`
		}

		r := Result{Address: w.Address}
		tc.Data["json"] = &r
		tc.ServeJSON()
		return
	}

	tc.CustomAbort(405, fmt.Sprintf("Withdraw failed %d %d %f %d", _token, addr, amount, userID))
}

// post withdraw address id, amount
// @router /tokens/:_token/tx/:_txId/notify [get]
func (tc *TokenController) WatchNotify(_token enums.TOKEN, _txId string) {
	beego.Debug("notify watching..", _txId)
	records := service.WalletNotify(_token, _txId)
	if records == nil || len(records) == 0 {
		beego.Error(_token, _txId)
		tc.Abort("405")
	} else {
		type Result struct {
			Num int `json:"num"`
		}

		r := Result{Num: len(records)}
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
	if tx := service.NewCold2HotTx(_token); tx != "" {
		type Result struct {
			Tx string `json:"tx"`
		}

		r := Result{Tx: tx}
		tc.Data["json"] = &r
		tc.ServeJSON()
	} else {
		tc.CustomAbort(405, fmt.Sprintf("NewCold2HotTx failed %d", _token))
	}
}

// @router /tokens/:_token/withdrawal/new [post]
func (tc *TokenController) NewWithdrawal(_token enums.TOKEN) () {
	userID, err := tc.GetInt("userID")
	address := tc.GetString("address")
	tag := tc.GetString("tag")

	if err != nil || len(address) == 0 || len(tag) == 0 {
		tc.CustomAbort(405, fmt.Sprintf("New withdrawal failed %s %s %s", err, address, tag))
		return
	}
	if w := service.NewWithdrawal(types.ID(userID), _token, address, tag); w > 0 {
		type Result struct {
			ID types.ID `json:"id"`
		}

		r := Result{ID: w}
		tc.Data["json"] = &r
		tc.ServeJSON()
		return
	}

	tc.CustomAbort(405, fmt.Sprintf("New withdrawal failed %d %d %s %s", _token, userID, address, tag))
}
