package tests

import (
	"fmt"
	"github.com/astaxie/beego/httplib"
	"testing"
	"io/ioutil"
)

func TestTokenController_Withdraw(t *testing.T) {
	url := fmt.Sprintf("%s/tokens/1/withdraw", HOST)
	req := httplib.Post(url)
	req.Param("amount", "9.9")
	req.Param("address", "2")
	req.Param("userID", "1")

	req.Debug(true)
	if resp, err := req.Response(); err != nil || resp == nil || resp.StatusCode > 300 {
		if resp == nil {
			t.Error(err)
			return
		}

		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		t.Error(resp.Status, string(body))
	}
}

func TestWithdraw_NewWithdrawal(t *testing.T) {
	url := fmt.Sprintf("%s/tokens/1/withdrawal/new", HOST)
	req := httplib.Post(url)
	req.Param("tag", "test tag")
	req.Param("address", "2NEhic4wTnBittzJru5r6SWP8LNjHjdE7nZ")
	req.Param("userID", "1")

	req.Debug(true)

	if resp, err := req.Response(); err != nil || resp.StatusCode > 300 {

		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		t.Error(resp.Status, string(body))
	}
}
