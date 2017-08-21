package service

import (
	"git.coding.net/zhouhuangjing/BitPurse/models/dao"
	"git.coding.net/zhouhuangjing/BitPurse/models/common"
)

type IService interface {
	TokenType() common.TOKEN
	SetTokenType(token common.TOKEN)
	Deposit(userId common.ID) *common.UserToken
	Withdraw(userId common.ID)
	Watch(userId common.ID)
	NewAddress() common.TokenAddress
}

type TokenService struct {
	IService

	tokenType_ common.TOKEN
}

func (ts *TokenService) SetTokenType(token common.TOKEN) {
	ts.tokenType_ = token
}

func (ts *TokenService) TokenType() common.TOKEN {
	return ts.tokenType_
}

func (ts *TokenService) Deposit(userId common.ID) *common.UserToken {

	// 存款的意思就是先检查用户是否已经生成了账户, 如果没有则先生成
	// 并将用户加入监控列表, 一旦发现账户有变动, 则刷新数据库
	// 监控用比特币的回调来做, 效率较高

	ut := dao.GetTokenByUser(userId, ts.TokenType())
	if ut == nil {
		ta := ts.NewAddress()
		ut = dao.NewTokenByUser(userId, ts.TokenType(), ta)
		if len(ta) > 0 {
			ts.Watch(userId)
		}
	}

	return ut
}

var (
	map_ = make(map[common.TOKEN]IService)
)

func Get(_id common.TOKEN) IService {
	return map_[_id]
}

func Reg(_id common.TOKEN, _service IService) {
	map_[_id] = _service
}
