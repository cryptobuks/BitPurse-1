package service

import (
	"git.coding.net/zhouhuangjing/BitPurse/models/common/enums"
	"git.coding.net/zhouhuangjing/BitPurse/models/common/models"
	"git.coding.net/zhouhuangjing/BitPurse/models/common/types"
	"git.coding.net/zhouhuangjing/BitPurse/models/dao"
	"git.coding.net/zhouhuangjing/BitPurse/models/rpc"
	"github.com/astaxie/beego"
)

type BitcoinService struct {
	TokenService
	rpc.BitcoinRpc
}

func initWatch() {
	tokens := dao.GetTokensByType(enums.TOKEN_BITCOIN)
	s := Get(enums.TOKEN_BITCOIN)
	for _, t := range tokens {
		s.Watch(t.TokenAddress)
	}
}

func InitBitcoin() IService {
	bs := Get(enums.TOKEN_BITCOIN)
	if bs == nil {
		bs = new(BitcoinService)
		bs.SetTokenType(enums.TOKEN_BITCOIN)
		Reg(enums.TOKEN_BITCOIN, bs)
	}

	initWatch()

	beego.Info("init bitcoin service")

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

func (bs *BitcoinService) Withdraw(_address string, _amount float64) {
	hash := bs.BitcoinRpc.Withdraw(_address, _amount)
	beego.Debug(hash)
}

func (bs *BitcoinService) WalletNotify(_txId string) *models.TokenRecord {
	tx := bs.BitcoinRpc.GetTransaction(_txId)
	if tx != nil {
		for _, v := range tx.Details {
			ut := dao.GetTokenByAddress(v.Address)
			u := ut.User
			if u == nil {
				beego.Error("user is nil")
				return nil
			}

			var category uint8
			if v.Category == "receive" {
				category = 1
			} else if v.Category == "send" {
				category = 0
			}

			tr := dao.NewTokenRecord(u.Id, ut.Token.TokenType, category, _txId)
			return tr
		}
	}

	return nil
}

// 检查
func (bs *BitcoinService) Watch(_address string) {
	bs.BitcoinRpc.Watch(_address)
	beego.Debug("I am watching", bs.TokenType(), _address)
}
