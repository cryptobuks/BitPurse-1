package tests

import (
	"path/filepath"
	"github.com/astaxie/beego"
)

func Init() error {

	if path, err := filepath.Abs("../.."); err == nil {
		beego.TestBeegoInit(path + "")
		return nil
	} else {
		return err
	}
}
