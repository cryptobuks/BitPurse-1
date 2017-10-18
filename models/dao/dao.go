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
		And(fmt.Sprintf("record_status & %d =  %d", enums.TX_DONE|enums.TX_STORED, enums.TX_DONE|enums.TX_STORED))

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

func MarkSendListSent(_txID string, _address string) bool {

	qb, _ := orm.NewQueryBuilder("mysql")

	set := fmt.Sprintf("record_status = record_status | %d", enums.TX_SENT)

	qb = qb.Update("token_record").Set(set)
	// 这里是条件！！！ 不是行为 !!
	qb = qb.Where("transaction_id = ?").And("record_address = ?").And("record_status & 1= ?").And("record_type = ?")

	result, err1 := ORM().Raw(qb.String(), _txID, _address, enums.TX_UNKNOWN, enums.OP_SEND).Exec()
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

func MarkRecordStatusDone(_txID string, _address string, _type enums.OP) bool {

	qb, _ := orm.NewQueryBuilder("mysql")

	set := fmt.Sprintf("record_status = record_status | %d", enums.TX_DONE)

	qb = qb.Update("token_record").Set(set)
	//  这里是条件 不是行为！！！
	qb = qb.Where("transaction_id = ?").And("record_address = ?").And("record_status & 4 =  ?").And("record_type = ?")

	result, err1 := ORM().Raw(qb.String(), _txID, _address, enums.TX_SENT, _type).Exec()
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
		return nil
	} else if err == orm.ErrMissPK {
		fmt.Println("No primary key found.")
		return nil
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
func GetTokenRecordByTxAddress(_txID string, _address string, _type enums.OP) *models.TokenRecord {
	o := ORM()

	tr := models.TokenRecord{TransactionId: _txID, RecordAddress: _address, RecordType: _type}
	err := o.Read(&tr, "TransactionId", "RecordAddress", "RecordType")
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

func GetSendList(_tokenID enums.TOKEN) []*models.TokenRecord {
	ut := new(models.TokenRecord)
	qs := ORM().QueryTable(ut).Filter("Token", _tokenID)
	qs = qs.Filter("RecordType", enums.OP_SEND).Filter("TransactionId", "").Filter("RecordStatus", enums.TX_UNKNOWN)

	records := make([]*models.TokenRecord, 0)
	if num, err := qs.All(&records); num > 0 && err == nil {
		return records
	}
	return nil
}

func UpdateSendListTx(_tokenID enums.TOKEN, _hash string) int64 {
	qs := ORM().QueryTable(new(models.TokenRecord)).Filter("Token", _tokenID)
	qs = qs.Filter("RecordType", enums.OP_SEND).Filter("TransactionId", "").Filter("RecordStatus", enums.TX_UNKNOWN)

	if num, err := qs.Update(orm.Params{"TransactionId": _hash}); num > 0 && err == nil {
		return num
	}
	return 0
}

//  1 deposit 2 withdraw
func NewTokenRecord(_userId types.ID, _tokenID enums.TOKEN, _recordType enums.OP, _txId string, _address string, _amount float64) *models.TokenRecord {

	o := ORM()
	// need confirm if needs query
	u := &models.User{Id: _userId}
	t := &models.Token{Id: _tokenID}

	status := enums.TX_UNKNOWN
	if _recordType == enums.OP_RECEIVE {
		status = status | enums.TX_SENT
	}
	tr := &models.TokenRecord{
		User:          u,
		Token:         t,
		TransactionId: _txId,

		RecordType:   _recordType,
		RecordStatus: status,

		RecordAddress: _address,
		RecordAmount:  _amount,
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
