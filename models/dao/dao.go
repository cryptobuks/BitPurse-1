package dao

import "database/sql"
import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego"
	"git.coding.net/zhouhuangjing/BitPurse/models/common/enums"
	"git.coding.net/zhouhuangjing/BitPurse/models/common/types"
	"git.coding.net/zhouhuangjing/BitPurse/models/common/models"
	"github.com/astaxie/beego/orm"
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

func GetTokenByUser(_userId types.ID, _tokenID enums.TOKEN) *models.UserToken {

	qs := ORM().QueryTable(new(models.UserToken))
	ut := new(models.UserToken)
	err := qs.Filter("Token", _tokenID).Filter("User", _userId).One(ut)

	if err == orm.ErrNoRows {
		beego.Error(err)
		return nil
	}
	return ut
}

func NewTokenByUser(_userId types.ID, _tokenID enums.TOKEN, _address string, _privateKey string) *models.UserToken {

	o := ORM()
	// need confirm if needs query
	u := &models.User{Id: _userId}
	t := &models.Token{Id: types.ID(_tokenID)}

	ut := &models.UserToken{
		User:  u,
		Token: t,

		TokenAddress: _address,
		PrivateKey:   _privateKey,
		TokenExtra:   "",
	}
	o.Insert(ut)

	return ut
}

func NewToken(_type enums.TOKEN, _symbol string, _name string, _intro string) int64 {
	t := models.Token{
		TokenType:   _type,
		TokenName:   _name,
		TokenSymbol: _symbol,
		TokenIntro:  _intro,
	}
	res, err := ORM().InsertOrUpdate(&t, "TokenType")
	if err != nil {
		beego.Error(err)
		return -1
	}
	return res
}
