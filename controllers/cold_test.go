package controllers

import (
	"testing"
	"os"
	"github.com/astaxie/beego/httplib"
	"io/ioutil"
	"fmt"
)

func TestFoo(t *testing.T) {
	t.Error("I am foo")

}

func TestTokenController_NewCold2HotTx(t *testing.T) {
	url := fmt.Sprintf("%s/tokens/1/cold2hot/new", HOST)
	req := httplib.Post(url)

	req.Debug(true)

	if resp, err := req.Response(); err == nil {

		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		t.Log(string(body))
		if resp.StatusCode > 300 {
			t.Error(resp.Status, string(body))
		}
	} else {
		t.Error(resp.Status)
	}
}

func TestMain(m *testing.M) {

	code := m.Run()
	os.Exit(code)

}
