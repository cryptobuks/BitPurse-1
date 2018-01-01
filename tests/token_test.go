package tests

import (
	"fmt"
	mycache "github.com/forchain/BitPurse/models/cache"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/httplib"
	"testing"
)

const HOST = "http://localhost:18080"



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
