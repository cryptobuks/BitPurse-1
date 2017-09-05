package dao

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" // import your required driver
)

var orm_ orm.Ormer

func GenerateTables(force bool) error {
	// Database alias.
	name := "default"

	// Error.
	err := orm.RunSyncdb(name, force, true)
	if err != nil {
		beego.Error(err)
	}
	return err
}

func ORM() orm.Ormer {
	return orm_
}

func InitORM() {
	orm.Debug = true
	// set default database
	ds := beego.AppConfig.String("dataSource")
	if len(ds) == 0 {
		beego.Warn("no data source")
	}
	err := orm.RegisterDataBase("default", "mysql", ds, 30)
	if err != nil {
		panic(err)
	}

	// register model
	registerModels()

	orm_ = orm.NewOrm()

	beego.Info("init orm ")
}
