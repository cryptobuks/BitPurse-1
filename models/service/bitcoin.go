package service

import (
	"git.coding.net/zhouhuangjing/BitPurse/models/common/enums"
	"git.coding.net/zhouhuangjing/BitPurse/models/common/models"
	"git.coding.net/zhouhuangjing/BitPurse/models/common/types"
	"git.coding.net/zhouhuangjing/BitPurse/models/rpc"
	"github.com/astaxie/beego"
)

type BitcoinService struct {
	TokenService
	rpc.BitcoinRpc
}

func InitBitcoin() IService {
	bs := Get(enums.TOKEN_BITCOIN)
	if bs == nil {
		bs = new(BitcoinService)
		bs.SetTokenType(enums.TOKEN_BITCOIN)
		Reg(enums.TOKEN_BITCOIN, bs)
	}

	return bs
}

func (bs *BitcoinService) NewAddress() (string, string) {

	beego.Debug("Do something before new bitcoin address")
	return bs.BitcoinRpc.NewAddress()
}

func (bs *BitcoinService) Deposit(userId types.ID) *models.UserToken {
	beego.Debug("Deposit", bs.TokenType())
	return bs.TokenService.Deposit(userId)
}

func (bs *BitcoinService) Withdraw(userId types.ID) {

}

func (bs *BitcoinService) Watch(userId types.ID) {
	beego.Debug("I am watching", bs.TokenType(), userId)
}
