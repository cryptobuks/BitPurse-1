package tests

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"testing"
	"io/ioutil"
)

func TestDeposit_OK(t *testing.T) {
	url := fmt.Sprintf("%s/tokens/1/deposit", HOST)
	req := httplib.Post(url)

	// test valid user and valid token
	req.Param("userID", "1")

	req.Debug(true)
	if res, err := req.String(); err != nil {
		t.Error(err)
	} else {
		beego.Debug(res)
	}
}

func TestDeposit_InvalidUser(t *testing.T) {
	url := fmt.Sprintf("%s/tokens/1/deposit", HOST)
	req := httplib.Post(url)

	// test valid user and invalid token
	req.Param("userID", "-1")

	req.Debug(true)

	if resp, err := req.Response(); err != nil || resp.StatusCode != 405 {

		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		t.Error(resp.Status, string(body))
	}
}

func TestDeposit_InvalidToken(t *testing.T) {
	url := fmt.Sprintf("%s/tokens/0/deposit", HOST)
	req := httplib.Post(url)

	// test valid user and valid token
	req.Param("userID", "1")

	req.Debug(true)
	if resp, err := req.Response(); err != nil || resp.StatusCode != 405 {

		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		t.Error(resp.Status, string(body))
	}
}
