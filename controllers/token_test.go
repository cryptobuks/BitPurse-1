package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"testing"
)

const HOST = "http://localhost:18080"

func TestTokenController_Withdraw(t *testing.T) {
	url := fmt.Sprintf("%s/users/1/tokens/1/withdraw", HOST)
	req := httplib.Post(url)
	req.Param("amount", "9.9")
	req.Param("address", "1")

	req.Debug(true)
	if res, err := req.String(); err != nil {
		t.Error(err)
	} else {
		beego.Debug(res)
	}
}
