package models

import "git.coding.net/zhouhuangjing/BitPurse/models/common/types"

type UserToken struct {
	Id   types.ID
	User *User `orm:"rel(fk)"`

	Token *Token `orm:"rel(fk)"`

	TokenAddress string
	PrivateKey   string
	TokenBalance float64
	LockBalance  float64
}

func (_self *UserToken) Balance() float64 {
	return _self.TokenBalance - _self.LockBalance
}

func (_self *UserToken) Lock(_amount float64) float64 {
	return _self.LockBalance + _amount
}

func (_self *UserToken) Unlock(_amount float64) float64 {
	if _self.TokenBalance <= _amount {
		return 0
	}
	return _self.TokenBalance - _amount
}
