package dao

import (
	_ "github.com/go-sql-driver/mysql" // import your required driver
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

var orm_ orm.Ormer

func GenerateTables(force bool) error {
	// Database alias.
	name := "default"

	// Print log.
	verbose := true

	// Error.
	err := orm.RunSyncdb(name, force, verbose)
	if err != nil {
		beego.Error(err)
	}
	return err
}

func initORM() {
	orm.Debug = true
	// set default database
	err := orm.RegisterDataBase("default", "mysql", "root@/bit_purse?charset=utf8", 30)
	if err != nil {
		panic(err)
	}

	// register model
	registerModels()
}

func ORM() orm.Ormer {
	if orm_ == nil {
		initORM()
		orm_ = orm.NewOrm()
	}

	return orm_
}
