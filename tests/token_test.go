package test

import (
	"fmt"
	mycache "git.coding.net/zhouhuangjing/BitPurse/models/cache"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
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

func TestTokenController_Deposit(t *testing.T) {
	url := fmt.Sprintf("%s/users/1/tokens/1/deposit", HOST)
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

func TestNewCache(t *testing.T) {
	config := fmt.Sprintf(`{"conn": "127.0.0.1:6379", "key": "%s"}`, "Test")
	c, err := cache.NewCache("my", config)
	x := c.(mycache.IMyCache)
	beego.Debug(x, err)

}

func TestTokenController_WatchNotify(t *testing.T) {
	url := fmt.Sprintf("%s/tokens/1/tx/0c86354dda907b3301dea5c2c8c32749879483f8a7add114974c20666559397f/notify", HOST)
	req := httplib.Get(url)
	if res, err := req.Response(); err != nil || res.StatusCode >= 300 {
		t.Error(err, res.Status)
	} else {
		beego.Debug(res)
	}
}
