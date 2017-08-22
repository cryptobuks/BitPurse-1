package service

import (
	"git.coding.net/zhouhuangjing/BitPurse/models/common"
	"git.coding.net/zhouhuangjing/BitPurse/models/rpc"
	"github.com/astaxie/beego"
)

type BitcoinService struct {
	TokenService
	rpc.BitcoinRpc
}

func InitBitcoin() IService {
	bs := Get(common.TOKEN_BITCOIN)
	if bs == nil {
		bs = new(BitcoinService)
		bs.SetTokenType(common.TOKEN_BITCOIN)
		Reg(common.TOKEN_BITCOIN, bs)
	}

	return bs
}

func (bs *BitcoinService) NewAddress() common.TokenAddress {

	beego.Debug("Do something before new bitcoin address")
	return bs.BitcoinRpc.NewAddress()
}

func (bs *BitcoinService) Deposit(userId common.ID) *common.UserToken {
	beego.Debug("Deposit", bs.TokenType())
	return bs.TokenService.Deposit(userId)
}

func (bs *BitcoinService) Withdraw(userId common.ID) {

}

func (bs *BitcoinService) Watch(userId common.ID) {
	beego.Debug("I am watching", bs.TokenType(), userId)
}

