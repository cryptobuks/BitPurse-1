package service

import (
	"git.coding.net/zhouhuangjing/BitPurse/models/common/enums"
	"git.coding.net/zhouhuangjing/BitPurse/models/common/models"
	"git.coding.net/zhouhuangjing/BitPurse/models/common/types"
	"git.coding.net/zhouhuangjing/BitPurse/models/dao"
	"github.com/astaxie/beego"
	"time"
	"strings"
	"strconv"
)

type IService interface {
	TokenID() enums.TOKEN
	Withdraw(_address string, _amount float64) string
	NewAddress() (string, string)
	WalletNotify(_txId string) *models.TokenRecord
	User2Hot() string
	IsUserAddress(_address string) bool
	HotAddress() string
	ColdAddress() string
	FeeRate() float64
	WithdrawFee() float64
	NewCold2HotTx() string
	GetBalanceByAddress(_address string) float64
	HotRate() float64
	NewTx(_from []string, _to map[string]float64, _change string) string
	SignTx(_tx string) (string, bool)
	SendTx(_tx string) string
	ValidateAddress(_address string) bool
}

func Deposit(_token enums.TOKEN, _userId types.ID) *models.UserToken {
	if ts := Get(_token); ts != nil {
		if u := dao.GetUser(_userId); u != nil {
			ut := dao.GetTokenByUser(_userId, ts.TokenID())
			if ut == nil {
				ta, pk := ts.NewAddress()
				ut = dao.NewTokenByUser(_userId, ts.TokenID(), ta, pk)
			}

			return ut
		}
	}

	return nil
}

func Withdraw(_tokenID enums.TOKEN, _userID types.ID, _withdrawalID types.ID, _amount float64) *models.Withdrawal {
	if ts := Get(_tokenID); ts != nil {
		if _amount < ts.WithdrawFee() {
			beego.Error("no enough withdraw fee", _amount, ts.WithdrawFee())
			return nil
		}
		w := dao.GetWithdrawal(_withdrawalID)
		if w == nil || w.User.Id != _userID {
			beego.Error("no right to withdraw", _userID, _withdrawalID)
			return nil
		}

		ut := dao.GetTokenByUser(_userID, _tokenID)
		balance := ut.Balance()
		if balance < _amount {
			beego.Error("no enough balance", _amount, ts.WithdrawFee())
			return nil
		}

		amount := _amount - ts.WithdrawFee()
		hotBalance := ts.GetBalanceByAddress(ts.HotAddress())
		if hotBalance < amount {
			beego.Error("no enough hot balance", _amount, hotBalance)
			return nil
		}

		if hash := ts.Withdraw(w.Address, amount); hash != "" {
			dao.NewTokenRecord(_userID, ts.TokenID(), enums.OP_SEND, hash)
			dao.UpdateLockBalance(_userID, ut.Lock(_amount))
		}
		return w
	}

	return nil
}

func WalletNotify(_token enums.TOKEN, _txId string) *models.TokenRecord {
	if ts := Get(_token); ts != nil {
		return ts.WalletNotify(_txId)
	}
	return nil
}

func User2Hot() {

	hotTime := beego.AppConfig.String("USER_2_HOT_TIME")

	timeParts := strings.Split(hotTime, ":")
	hour := 0
	if h, err1 := strconv.Atoi(timeParts[0]); err1 == nil {
		hour = h % 24
	}

	min := 0
	if len(timeParts) > 1 {
		if m, err2 := strconv.Atoi(timeParts[1]); err2 == nil {
			min = m % 60
		}
	}
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, time.Local)

	if now.Hour()*60+now.Minute() < hour*60+min {
		next.Add(24 * time.Hour)
	}

	delay := time.Until(next)

	timer := time.NewTimer(delay)
	go func() {
		<-timer.C
		ticker := time.NewTicker(24 * time.Hour)
		for t := range ticker.C {
			beego.Info("User2Hot", t)
			for token, s := range map_ {
				if tx := s.User2Hot(); tx != "" {
					dao.MarkRecordStatusStored(token)
				}
			}
		}
	}()
}

func SignTx(_token enums.TOKEN, _tx string) (string, bool) {
	if ts := Get(_token); ts != nil {
		return ts.SignTx(_tx)
	}

	return "", false
}

func SendTx(_token enums.TOKEN, _tx string) string {
	if ts := Get(_token); ts != nil {
		return ts.SendTx(_tx)
	}

	return ""
}

// 得到冷钱包余额, 得到热钱包余额, 然后得到应该的转账金额, 生成一笔冷到热的交易

// 先来个BalanceByAddress
func NewCold2HotTx(_token enums.TOKEN) string {
	if ts := Get(_token); ts != nil {
		hotBalance := ts.GetBalanceByAddress(ts.HotAddress())
		coldBalance := ts.GetBalanceByAddress(ts.ColdAddress())

		sum := coldBalance + hotBalance
		hotRate := ts.HotRate()
		if hotBalance < sum*hotRate {
			amount := sum*hotRate - hotBalance
			to := make(map[string]float64)
			to[ts.HotAddress()] = amount
			return ts.NewTx([]string{ts.ColdAddress()}, to, ts.HotAddress())
		}

		beego.Info("No need transfer from cold to hot since hot is greater than 10% ")
	}

	return ""
}

func NewWithdrawal(_userID types.ID, _token enums.TOKEN, _address string, _tag string) types.ID {
	if len(_tag) == 0 || len(_tag) > 99 {
		beego.Error("Invalid tag", _tag)
		return -1
	}

	if ts := Get(_token); ts != nil {
		if ts.ValidateAddress(_address) {
			if ut := dao.GetTokenByUser(_userID, _token); ut != nil {
				if r := dao.NewWithdrawal(_userID, _token, _address, _tag); r > 0 {
					return r
				}
			}
		}
	}

	return -2
}

var (
	map_ = make(map[enums.TOKEN]IService)
)

func Get(_id enums.TOKEN) IService {
	if r, ok := map_[_id]; ok {
		return r
	}
	beego.Error("No service", _id)
	return nil
}

func Reg(_id enums.TOKEN, _service IService) {
	map_[_id] = _service
}

func init() {
	InitBitcoin()
	User2Hot()
}
