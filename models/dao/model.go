package dao

import (
	"github.com/astaxie/beego/orm"
	"github.com/forchain/BitPurse/models/common/models"
)

func registerModels() {
	orm.RegisterModel(new(models.UserToken))
	orm.RegisterModel(new(models.Withdrawal))
	orm.RegisterModel(new(models.User))
	orm.RegisterModel(new(models.TokenRecord))
	orm.RegisterModel(new(models.Token))
}
