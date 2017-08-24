package models

import (
	"time"
	"git.coding.net/zhouhuangjing/BitPurse/models/common/types"
	"git.coding.net/zhouhuangjing/BitPurse/models/common/enums"
)

type UserToken struct {
	Id   types.ID
	User *User `orm:"rel(fk)"`

	Token *Token `orm:"rel(fk)"`

	TokenAddress string
	PrivateKey   string
	TokenBalance float64
	TokenExtra   string
}

type Withdrawal struct {
	Id      types.ID
	User    *User `orm:"rel(fk)"`
	Address uint8
	Tag     string
	Token   *Token `orm:"rel(fk)"`
}

type User struct {
	Id           types.ID
	UserName     string `orm:"unique"`
	UserPassword string
	MailAddress  string
	MailCode     string
	PhoneNo      string `orm:"size(11);unique"`
	CountryCode  string
	CreateTime   time.Time `orm:"auto_now_add;type(datetime)"`
	UserPortrait string
	UserIntro    string
	Tokens       []*UserToken `orm:"reverse(many)"`
	Withdrawals  []*Withdrawal  `orm:"reverse(many)"`
	Records      []*TokenRecord  `orm:"reverse(many)"`
}

type Token struct {
	Id          types.ID
	TokenType   enums.TOKEN `orm:"unique"`
	TokenSymbol string `orm:"unique"`
	TokenName   string
	TokenIntro  string
}

type TokenRecord struct {
	Id            types.ID
	RecordTime    uint64
	RecordType    uint8
	Token         *Token `orm:"rel(fk)"`
	TransactionId string
	User          *User `orm:"rel(fk)"`
}
