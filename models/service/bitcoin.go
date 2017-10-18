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

func (_self *BitcoinService) CheckBalance() bool {
	checked := true
	if unspent := _self.BitcoinRpc.ListUnspent(); unspent != nil {
		groupings := _self.BitcoinRpc.Address2BalanceMap()

		unspentMap := make(map[string]float64)
		unspentSum := 0.0
		for k, u := range unspent {
			if u.Address == "2NEhic4wTnBittzJru5r6SWP8LNjHjdE7nZ" {
				beego.Debug("cold", k, u.Amount)
			}
			unspentMap[u.Address] = u.Amount
			unspentSum += u.Amount

			if amount, ok := groupings[u.Address]; !ok {
				checked = false
				beego.Error("Not found ", u.Address)
			} else {
				if amount != u.Amount {
					checked = false
					//beego.Error("Not equal ", u.Amount, amount)
				}
			}
		}
		beego.Debug("unspent ", unspentSum, len(unspentMap))

		groupSum := 0.0
		beego.Debug(len(unspent), len(groupings))
		for address, amount := range groupings {
			groupSum += amount
			if ua, ok := unspentMap[address]; !ok {
				checked = false
				beego.Error("Not found ", address)
			} else {
				if amount != ua {
					checked = false
					//beego.Error("Not equal ", ua, amount)
				}
			}
		}
		beego.Debug("groupings ", groupSum, len(groupings))
		if groupSum != unspentSum {
			checked = false
		}

		beego.Debug("balance", _self.BitcoinRpc.Balance())
	}

	return checked
}

func (_self *BitcoinService) UnspentBalance() float64 {

	if result := _self.BitcoinRpc.ListUnspent(); result != nil {
		sum := 0.0
		for _, v := range result {
			sum += v.Amount
		}
		return sum
	}

	return -1
}

func (_self *BitcoinService) TokenID() enums.TOKEN {
	return enums.TOKEN_BITCOIN
}

func (_self *BitcoinService) ColdAddress() string {
	return beego.AppConfig.String("COLD_ADDRESS")
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

func (_self *BitcoinService) Transfer(_from []string, _to map[string]float64, _changeAddress string) string {
	return _self.BitcoinRpc.Transfer(_from, _to, _changeAddress)
}

func (_self *BitcoinService) ValidateAddress(_address string) bool {
	return _self.BitcoinRpc.ValidateAddress(_address)
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
	if rate, err := beego.AppConfig.Float("HOT_RATE"); err == nil && rate > 0 && rate < 1 {
		return rate
	} else {
		beego.Error("Invalid hot rate", rate)
		return -1
	}
}

func (_self *BitcoinService) GetBalanceByAddress(_address string) float64 {
	// balance will be locked if one transaction is sent but packed not into block yet
	if tx := _self.BitcoinRpc.ListUnspentByAddress(_address); tx != nil {
		amount := 0.0
		for _, v := range tx {
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

func (_self *BitcoinService) User2General() string {
	if balanceMap := _self.BitcoinRpc.Address2BalanceMap(); balanceMap != nil {
		from := make([]string, 0)
		amount := 0.0
		hotAddr := _self.HotAddress()
		coldAddr := _self.ColdAddress()
		for addr, balance := range balanceMap {
			if balance > 0 && coldAddr != addr && hotAddr != addr {
				from = append(from, addr)
				amount += balance
			}
		}

		hotAmount, ok1 := balanceMap[hotAddr]
		coldAmount, ok2 := balanceMap[coldAddr]

		if !ok1 || !ok2 {
			beego.Error("no general amount", coldAmount, hotAmount)
			return ""
		}

		to := make(map[string]float64)

		sum := amount + coldAmount + hotAmount
		if hotRate := _self.HotRate(); hotRate > 0 {
			if (hotAmount+amount)/sum <= hotRate {
				to[hotAddr] = amount
			} else {
				toCold := (hotAmount + amount) - sum*hotRate
				toHot := amount - toCold
				to[hotAddr] = toHot
				to[coldAddr] = toCold
			}
		}

		if tx := _self.Transfer(from, to, hotAddr); tx != "" {
			return tx
		}
	}

	beego.Error("User2General failed")
	return ""
}

func (_self *BitcoinService) User2GeneralOld() string {
	if list := _self.BitcoinRpc.ListUnspent(); list != nil {
		var amount float64
		var coldAmount float64

		from := make([]string, 0)
		fromMap := make(map[string]interface{})
		for _, v := range list {
			if _self.ColdAddress() == v.Address {
				coldAmount += v.Amount
			} else if _self.HotAddress() != v.Address {
				amount += v.Amount

				if _, ok := fromMap[v.Address]; !ok {
					from = append(from, v.Address)
					var empty interface{}
					fromMap[v.Address] = empty
				}
			}
		}

		balanceMap := _self.BitcoinRpc.Address2BalanceMap()
		hotAmount, ok := balanceMap[_self.HotAddress()]

		if !ok || coldAmount == 0 || hotAmount == 0 {
			beego.Error("no general amount", coldAmount, hotAmount)
			return ""
		}

		to := make(map[string]float64)

		sum := amount + coldAmount + hotAmount
		if hotRate := _self.HotRate(); hotRate > 0 {
			if (hotAmount+amount)/sum <= hotRate {
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

	beego.Error("User2General failed")

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
func (_self *BitcoinService) WalletNotify(_txId string) []*models.TokenRecord {
	if tx := _self.BitcoinRpc.GetTransaction(_txId); tx != nil {
		records := make([]*models.TokenRecord, 0)
		needSend := false
		for _, v := range tx.Details {
			var tr *models.TokenRecord
			if v.Category == "receive" {
				if ut := dao.GetTokenByAddress(v.Address); ut != nil {
					u := ut.User
					if u == nil {
						beego.Error("user is nil")
						return nil
					}
					if tx.Confirmations == 0 {
						tr = dao.NewTokenRecord(u.Id, ut.Token.Id, enums.OP_RECEIVE, _txId, v.Address, v.Amount)
					} else {
						tr = dao.GetTokenRecordByTxAddress(_txId, v.Address, enums.OP_RECEIVE)
						dao.UpdateTokenBalance(ut.Id, ut.TokenBalance+v.Amount)
						dao.MarkRecordStatusDone(_txId, v.Address, enums.OP_RECEIVE)
					}
				} else {
					beego.Error("[WalletNotify]No user token", v.Address)
				}
			} else if v.Category == "send" {
				if tr = dao.GetTokenRecordByTxAddress(_txId, v.Address, enums.OP_SEND); tr != nil {
					if ut := dao.GetTokenByUser(tr.User.Id, tr.Token.Id); ut != nil {
						if tx.Confirmations == 0 {
							dao.MarkSendListSent(_txId, v.Address)
						} else {
							needSend = true
							if amount := -v.Amount; amount > 0 {
								dao.UpdateLockBalance(ut.Id, ut.Unlock(amount+_self.WithdrawFee()))
								dao.UpdateTokenBalance(ut.Id, ut.TokenBalance-amount)
								dao.MarkRecordStatusDone(_txId, v.Address, enums.OP_SEND)
							}
						}
					} else {
						beego.Error("[WalletNotify]No user token", tr.User.Id, tr.Token.Id)
					}
				} else {
					beego.Error("[WalletNotify]No token record", _txId)
				}
			}

			records = append(records, tr)
		}

		if needSend {
			if list := dao.GetSendList(enums.TOKEN_BITCOIN); list != nil && len(list) > 0 {
				to := make(map[string]float64)
				sum := 0.0
				for _, v := range list {
					sum += v.RecordAmount
					to[v.RecordAddress] = v.RecordAmount
				}
				hotBalance := _self.GetBalanceByAddress(_self.HotAddress())
				if hotBalance > sum {
					from := []string{_self.HotAddress()}
					if hash := _self.Transfer(from, to, _self.HotAddress()); hash != "" {
						if num := dao.UpdateSendListTx(enums.TOKEN_BITCOIN, hash); int(num) == len(to) {
							return records
						}
					}
				}
			}
		}
		return records
	}

	beego.Error("[WalletNotify]No tx", _txId)
	return nil
}
