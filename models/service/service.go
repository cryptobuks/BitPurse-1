package service

import (
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
	Withdraw(userId types.ID)
	Watch(userId types.ID)
	NewAddress() (string, string)
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

func Deposit(ts IService, userId types.ID) *models.UserToken {

	// 存款的意思就是先检查用户是否已经生成了账户, 如果没有则先生成
	// 并将用户加入监控列表, 一旦发现账户有变动, 则刷新数据库
	// 监控用比特币的回调来做, 效率较高

	ut := dao.GetTokenByUser(userId, ts.TokenType())
	if ut == nil {
		ta, pk := ts.NewAddress()
		ut = dao.NewTokenByUser(userId, ts.TokenType(), ta, pk)
		if len(ta) > 0 {
			ts.Watch(userId)
		}
	}

	return ut
}

func Withdraw(ts IService, userId types.ID, _address types.ID, _amount float64) *models.UserToken {
	//  提款的意思就是, 调用比特币的提款服务, 记得要添加记录

	beego.Debug("needs implement", _address, _amount)
	ut := dao.GetTokenByUser(userId, ts.TokenType())
	if ut == nil {
		beego.Error("no deposit address")
		return nil
	}

	return ut
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
