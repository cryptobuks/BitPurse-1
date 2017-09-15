package service

import (
	"git.coding.net/zhouhuangjing/BitPurse/models/common/enums"
	"git.coding.net/zhouhuangjing/BitPurse/models/common/models"
	"git.coding.net/zhouhuangjing/BitPurse/models/dao"
	"git.coding.net/zhouhuangjing/BitPurse/models/rpc"
	"github.com/astaxie/beego"
)

type BitcoinService struct {
	rpc.BitcoinRpc
	IService
}

func InitBitcoin() IService {
	bs := new(BitcoinService)
	Reg(enums.TOKEN_BITCOIN, bs)

	beego.Info("init bitcoin service")

	return bs
}

func (_self *BitcoinService) TokenID() enums.TOKEN {
	return enums.TOKEN_BITCOIN
}

func (_self *BitcoinService) ColdAddress() string {
	return beego.AppConfig.String("HOT_ADDRESS")
}

func (_self *BitcoinService) HotAddress() string {
	return beego.AppConfig.String("HOT_ADDRESS")
}
func (_self *BitcoinService) SignTx(_tx string) (string, bool) {
	return _self.BitcoinRpc.SignTx(_tx)
}

func (_self *BitcoinService) SendTx(_tx string) string {
	return _self.BitcoinRpc.SendTx(_tx)
}

func (_self *BitcoinService) NewTx(_from []string, _to map[string]float64, _change string) string {
	return _self.BitcoinRpc.NewTx(_from, _to, _change)
}

func (_self *BitcoinService) WithdrawFee() float64 {
	if fee, err := beego.AppConfig.Float("WITHDRAW_FEE"); err == nil {
		return fee
	}
	return 0
}

func (_self *BitcoinService) FeeRate() float64 {
	if rate, err := beego.AppConfig.Float("FEE_RATE"); err == nil {
		return rate
	}
	return 0
}

func (_self *BitcoinService) IsUserAddress(_address string) bool {
	return _address != _self.HotAddress() && _self.ColdAddress() != _address
}

// ToDo: Load from database
func (_self *BitcoinService) HotRate() float64 {
	if rate, err := beego.AppConfig.Float("HOT_RATE"); err == nil {
		return rate
	}
	return -1
}

func (_self *BitcoinService) GetBalanceByAddress(_address string) float64 {
	if tx := _self.BitcoinRpc.ListUnspentByAddress(_address); tx != nil {
		amount := 0.0
		for _, v := range *tx {
			amount += v.Amount
		}
		return amount
	}

	return -1
}

func (_self *BitcoinService) NewCold2HotTx() string {
	coldBalance := _self.GetBalanceByAddress(_self.ColdAddress())
	hotBalance := _self.GetBalanceByAddress(_self.HotAddress())
	sum := coldBalance + hotBalance
	if r := _self.HotRate(); r > 0 {
		if hotBalance < sum*r {
			amount := sum*r - hotBalance

			from := []string{_self.ColdAddress()}
			to := map[string]float64{
				_self.HotAddress(): amount,
			}
			tx := _self.NewTx(from, to, _self.ColdAddress())
			return tx
		}
	}

	beego.Error("NewCold2HotTx failed", coldBalance, hotBalance, sum)

	return ""
}

func (_self *BitcoinService) User2Hot() string {
	if list := _self.BitcoinRpc.ListUnspent(); list == nil {
		var amount float64
		var coldAmount float64
		var hotAmount float64

		from := make([]string, 0)
		for _, v := range *list {
			if _self.HotAddress() == v.Address {
				hotAmount += v.Amount
			} else if _self.ColdAddress() == v.Address {
				coldAmount += v.Amount
			} else {
				amount += v.Amount
				from = append(from, v.Address)
			}
		}

		to := make(map[string]float64)

		sum := amount + coldAmount + hotAmount
		if hotRate, err := beego.AppConfig.Float("HOT_RATE"); err == nil {
			if hotAmount+amount/sum <= hotRate {
				to[_self.HotAddress()] = amount
			} else {
				toCold := (hotAmount + amount) - sum*hotRate
				toHot := amount - toCold
				to[_self.HotAddress()] = toHot
				to[_self.ColdAddress()] = toCold
			}
		}

		if tx := _self.Transfer(from, to, _self.HotAddress()); tx != "" {
			return tx
		}
	}

	beego.Error("User2Hot failed")

	return ""
}

func (_self *BitcoinService) NewAddress() (string, string) {

	beego.Debug("Do something before new bitcoin address")
	return _self.BitcoinRpc.NewAddress()
}

// 提款的时候还是应该先插入一条记录, 并设成未完成状态, 等收到通知后再设成已完成
// 提款是从热钱包提!!! 而且应该是从收款列表中选
func (_self *BitcoinService) Withdraw(_address string, _amount float64) string {

	hotAddr := []string{_self.HotAddress()}
	to := make(map[string]float64)
	to[_address] = _amount
	hash := _self.BitcoinRpc.Transfer(hotAddr, to, _self.HotAddress())

	beego.Debug(hash)

	return hash
}

// UTXO 比较特别, 所以接收和发送要分别处理
func (_self *BitcoinService) WalletNotify(_txId string) *models.TokenRecord {
	if tx := _self.BitcoinRpc.GetTransaction(_txId); tx != nil {
		for _, v := range tx.Details {
			ut := dao.GetTokenByAddress(v.Address)
			u := ut.User
			if u == nil {
				beego.Error("user is nil")
				return nil
			}

			var tr *models.TokenRecord
			if v.Category == "receive" {
				tr = dao.NewTokenRecord(u.Id, ut.Token.Id, enums.OP_RECEIVE, _txId)
			} else if v.Category == "send" {
				tr = dao.GetTokenRecordByTx(_txId)
			}

			dao.MarkRecordStatusDone(_txId)
			dao.UpdateLockBalance(u.Id, ut.Unlock(v.Amount+_self.WithdrawFee()))

			return tr
		}
	}

	return nil
}
