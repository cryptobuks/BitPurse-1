package models

import (
	"../enums"
	"../types"
	"time"
)

type Withdrawal struct {
	Id      types.ID
	User    *User  `orm:"rel(fk)"`
	Address string
	Tag     string
	Token   *Token `orm:"rel(fk)"`
}

type User struct {
	Id           types.ID
	UserName     string         `orm:"unique"`
	UserPassword string
	MailAddress  string
	MailCode     string
	PhoneNo      string         `orm:"size(11);unique"`
	CountryCode  string
	CreateTime   time.Time      `orm:"auto_now_add;type(datetime)"`
	UserPortrait string
	UserIntro    string
	Tokens       []*UserToken   `orm:"reverse(many)"`
	Withdrawals  []*Withdrawal  `orm:"reverse(many)"`
	Records      []*TokenRecord `orm:"reverse(many)"`
}

type Token struct {
	Id          enums.TOKEN `orm:"pk"`
	TokenSymbol string      `orm:"unique"`
	TokenName   string
	TokenIntro  string
	TokenFee    float64
}

type TokenRecord struct {
	Id            types.ID
	RecordTime    time.Time `orm:"auto_now_add;type(datetime)"`
	RecordType    enums.OP
	Token         *Token    `orm:"rel(fk)"`
	TransactionId string
	User          *User     `orm:"rel(fk)"`
	RecordStatus  enums.TX
	RecordAddress string
	RecordAmount  float64
}
