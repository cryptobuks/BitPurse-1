package service

import (
	"git.coding.net/zhouhuangjing/BitPurse/models/cache"
	"git.coding.net/zhouhuangjing/BitPurse/models/common/enums"
	"git.coding.net/zhouhuangjing/BitPurse/models/common/models"
	"git.coding.net/zhouhuangjing/BitPurse/models/common/types"
	"git.coding.net/zhouhuangjing/BitPurse/models/dao"
	"github.com/astaxie/beego"
)

type IService interface {
	TokenType() enums.TOKEN
	SetTokenType(token enums.TOKEN)
	Deposit(userId types.ID) *models.UserToken
	Withdraw(_address string, _amount float64)
	Watch(_address string)
	NewAddress() (string, string)
	WalletNotify(_txId string) *models.TokenRecord
}

type TokenService struct {
	IService

	tokenType_ enums.TOKEN
}

func (ts *TokenService) SetTokenType(token enums.TOKEN) {
	ts.tokenType_ = token
}

func (ts *TokenService) TokenType() enums.TOKEN {
	return ts.tokenType_
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

	ts.Withdraw(w.Address, _amount)

	return w
}

func WalletNotify(_token enums.TOKEN, _txId string) *models.TokenRecord {
	ts := Get(_token)
	return ts.WalletNotify(_txId)
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
}
