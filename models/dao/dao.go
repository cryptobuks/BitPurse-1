package dao

import "database/sql"
import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego"
	"git.coding.net/zhouhuangjing/BitPurse/models/common"
)

var db *sql.DB

func Init() error {
	var err error
	db, err = Connect()
	if err != nil {
		beego.Error("init db failed", err)
	}
	return err
}

func Connect() (*sql.DB, error) {
	if db != nil {
		return db, nil
	}

	var err error
	db, err = sql.Open("mysql", "root@/BitPurse")
	if err == nil {
		if err = db.Ping(); err != nil {
			defer db.Close()
			return db, err
		}
	}

	return db, err
}

func GetTokenByUser(_userId common.ID, _tokenType common.TOKEN) *common.UserToken {
	row := db.QueryRow("SELECT *  FROM UserToken WHERE UserId=? AND TokenType=?", _userId, _tokenType)
	ut := new(common.UserToken)
	err := row.Scan(&ut.UserTokenID, &ut.UserId, &ut.TokenType, &ut.TokenAddress, &ut.TokenBalance, &ut.TokenExtra)

	if err == sql.ErrNoRows {
		return nil
	}
	return ut
}

func NewTokenByUser(_userId common.ID, _tokenType common.TOKEN, _address common.TokenAddress) *common.UserToken {
	result, err := db.Exec(
		"INSERT INTO UserToken (UserId, TokenType, TokenAddress, TokenBalance, TokenExtra) VALUES (?, ?, ?, ?, ?)",
		_userId, _tokenType, string(_address), 0, "")

	if err != nil {
		beego.Error(err)
		return nil
	}
	utId, err := result.LastInsertId()
	ut := common.UserToken{
		UserTokenID:  common.ID(utId),
		UserId:       _userId,
		TokenType:    uint8(_tokenType),
		TokenAddress: string(_address),
		TokenExtra:   "",
	}

	return &ut
}
