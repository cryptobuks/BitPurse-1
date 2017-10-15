package dao

import (
	"fmt"
	"git.coding.net/zhouhuangjing/BitPurse/models/common/enums"
	"git.coding.net/zhouhuangjing/BitPurse/models/common/models"
	"git.coding.net/zhouhuangjing/BitPurse/models/common/types"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func MarkRecordStatusStored(_tokenID enums.TOKEN) bool {

	qb, _ := orm.NewQueryBuilder("mysql")

	set := fmt.Sprintf("record_status = record_status | %d", enums.TX_STORED)

	qb.Update("token_record").Set(set).Where("token_id = ?").
		And(fmt.Sprintf("record_status & %d =  %d", enums.TX_DONE|enums.TX_STORED, enums.TX_DONE|enums.TX_STORED)).
		And(fmt.Sprintf("record_status & %d = 0", enums.TX_SPENT))

	result, err1 := ORM().Raw(qb.String(), _tokenID).Exec()
	if err1 != nil {
		beego.Error(err1)
		return false
	}

	id, err2 := result.LastInsertId()
	rows, err3 := result.RowsAffected()
	if err2 != nil || err3 != nil {
		beego.Error(err2, id, err3, rows)
		return false
	}

	return true
}

func MarkRecordStatusDone(_txID string) bool {

	qb, _ := orm.NewQueryBuilder("mysql")

	set := fmt.Sprintf("record_status = record_status | %d", enums.TX_DONE)

	qb.Update("token_record").Set(set).Where("transaction_id = ?")

	result, err1 := ORM().Raw(qb.String(), _txID).Exec()
	if err1 != nil {
		beego.Error(err1)
		return false
	}

	id, err2 := result.LastInsertId()
	rows, err3 := result.RowsAffected()
	if err2 != nil || err3 != nil {
		beego.Error(err2, id, err3, rows)
		return false
	}

	return true
}

func GetUser(_userID types.ID) *models.User {
	o := ORM()
	u := models.User{
		Id: _userID,
	}
	if err := o.Read(&u); err == nil {
		return &u
	}

	beego.Error("No user", _userID)
	return nil
}

func NewWithdrawal(_userID types.ID, _token enums.TOKEN, _address string, _tag string) types.ID {
	o := ORM()
	w := &models.Withdrawal{
		User:    &models.User{Id: _userID},
		Address: _address,
		Tag:     _tag,
		Token:   &models.Token{Id: _token},
	}

	if r, err := o.Insert(w); err == nil && r > 0 {
		return types.ID(r)
	}

	beego.Error("NewWithdrawal failed", _userID, _token, _address, _tag)
	return -1
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

func GetWithdrawal(_id types.ID) *models.Withdrawal {
	o := ORM()
	w := models.Withdrawal{Id: _id}

	err := o.Read(&w)

	if err == orm.ErrNoRows {
		fmt.Println("No result found.")
	} else if err == orm.ErrMissPK {
		fmt.Println("No primary key found.")
	} else {
		fmt.Println(w.Id, w.Address)
	}

	return &w

}

func GetTokenByAddress(_address string) *models.UserToken {
	o := ORM()

	qs := o.QueryTable(new(models.UserToken))

	ut := new(models.UserToken)
	err := qs.Filter("TokenAddress", _address).RelatedSel("user").One(ut)
	if err != nil {
		beego.Error(err)
		return nil
	}
	return ut

}
func GetTokenRecordByTx(_txID string) *models.TokenRecord {
	o := ORM()

	tr := models.TokenRecord{TransactionId: _txID}
	err := o.Read(&tr, "TransactionID")
	if err != nil {
		beego.Error(err)
		return nil
	}
	return &tr

}

func NewTokenByUser(_userId types.ID, _tokenID enums.TOKEN, _address string, _privateKey string) *models.UserToken {

	o := ORM()
	// need confirm if needs query
	u := &models.User{Id: _userId}
	t := &models.Token{Id: _tokenID}

	ut := &models.UserToken{
		User:  u,
		Token: t,

		TokenAddress: _address,
		PrivateKey:   _privateKey,
	}
	o.Insert(ut)

	return ut
}

func UpdateTokenBalance(_utID types.ID, _amount float64) bool {
	o := ORM()

	ut := &models.UserToken{Id: _utID, TokenBalance: _amount}

	if num, err := o.Update(ut, "TokenBalance"); err != nil || num != 1 {
		beego.Error("update balance failed", num, err)
		return false
	}

	return true
}

func UpdateLockBalance(_utID types.ID, _balance float64) bool {
	o := ORM()

	ut := &models.UserToken{Id: _utID, LockBalance: _balance}

	if num, err := o.Update(ut, "LockBalance"); err != nil || num != 1 {
		beego.Error("Lock balance failed", num, err)
		return false
	}

	return true
}

//  1 deposit 2 withdraw
func NewTokenRecord(_userId types.ID, _tokenID enums.TOKEN, _recordType enums.OP, _txId string) *models.TokenRecord {

	o := ORM()
	// need confirm if needs query
	u := &models.User{Id: _userId}
	t := &models.Token{Id: _tokenID}

	tr := &models.TokenRecord{
		User:          u,
		Token:         t,
		TransactionId: _txId,

		RecordType:   _recordType,
		RecordStatus: enums.TX_UNKNOWN,
	}
	index, err := o.Insert(tr)
	if err != nil {
		beego.Error(index, err)
		return nil
	}

	return tr
}

func NewToken(_ID enums.TOKEN, _symbol string, _name string, _intro string) int64 {
	t := models.Token{
		Id:          _ID,
		TokenName:   _name,
		TokenSymbol: _symbol,
		TokenIntro:  _intro,
	}
	res, err := ORM().InsertOrUpdate(&t, "TokenID")
	if err != nil {
		beego.Error(err)
		return -1
	}
	return res
}

func Init() {
	InitORM()
}

func init() {
}
