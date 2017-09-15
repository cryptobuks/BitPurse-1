package service

import (
	"git.coding.net/zhouhuangjing/BitPurse/models/cache"
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
	TokenType() enums.TOKEN
	Withdraw(_address string, _amount float64) string
	Watch(_address string)
	NewAddress() (string, string)
	WalletNotify(_txId string) *models.TokenRecord
	User2Hot() string
	IsUserAddress(_address string) bool
	HotAddress() string
	ColdAddress() string
	FreeRate() float64
	NewCold2HotTx() string
	GetBalanceByAddress(_address string) float64
	HotRate() float64
	NewTx(_from []string, _to map[string]float64, _change string) string
	SignTx(_tx string) (string, bool)
	SendTx(_tx string) string
}

// category/type/btc/
// 监控的时候将自己的地址放到集合里面, 这个操作是o1,
// 执行的时候一次性取出来遍历,  列表和集合都可以, 优先列表吧, 操作比较多
func Watch(_ts IService, _address string) {
	c := cache.New("WATCH")
	if c == nil {
		return
	}
	c.Put(_address, _ts.TokenType(), 0)
	_ts.Watch(_address)

	beego.Info("watch", _ts.TokenType(), _address)
}
func Deposit(_token enums.TOKEN, _userId types.ID) *models.UserToken {

	ts := Get(_token)
	// 存款的意思就是先检查用户是否已经生成了账户, 如果没有则先生成
	// 并将用户加入监控列表, 一旦发现账户有变动, 则刷新数据库
	// 监控用比特币的回调来做, 效率较高

	ut := dao.GetTokenByUser(_userId, ts.TokenType())
	if ut == nil {
		ta, pk := ts.NewAddress()
		ut = dao.NewTokenByUser(_userId, ts.TokenType(), ta, pk)
		if len(ta) > 0 {
			Watch(ts, ta)
		}
	}

	return ut
}

func Withdraw(_token enums.TOKEN, _userId types.ID, _address types.ID, _amount float64) *models.Withdrawal {
	ts := Get(_token)

	w := dao.GetWithdrawal(_address)
	if w.User.Id != _userId {
		beego.Error("no right to withdraw", _userId, _address)
		return nil
	}

	hash := ts.Withdraw(w.Address, _amount)
	if hash != "" {
		dao.NewTokenRecord(_userId, ts.TokenType(), enums.OP_SEND, hash)
	}

	return w
}

func WalletNotify(_token enums.TOKEN, _txId string) *models.TokenRecord {
	ts := Get(_token)
	return ts.WalletNotify(_txId)
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
	ts := Get(_token)
	return ts.SignTx(_tx)
}

func SendTx(_token enums.TOKEN, _tx string) string {
	ts := Get(_token)
	return ts.SendTx(_tx)
}

// 得到冷钱包余额, 得到热钱包余额, 然后得到应该的转账金额, 生成一笔冷到热的交易

// 先来个BalanceByAddress
func NewCold2HotTx(_token enums.TOKEN) string {
	ts := Get(_token)

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

	return ""
}

var (
	map_ = make(map[enums.TOKEN]IService)
)

func Get(_id enums.TOKEN) IService {
	return map_[_id]
}

func Reg(_id enums.TOKEN, _service IService) {
	map_[_id] = _service
}

func init() {
	InitBitcoin()
	User2Hot()
}
