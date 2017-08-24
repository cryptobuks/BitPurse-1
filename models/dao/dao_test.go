package dao

import (
	"testing"
	"os"
	"git.coding.net/zhouhuangjing/BitPurse/models/common/models"
	"git.coding.net/zhouhuangjing/BitPurse/models/common/enums"
	"math/rand"
	"github.com/astaxie/beego"
)

func TestCreateTables(t *testing.T) {

	if err := GenerateTables(true); err != nil {
		t.Error("insert by orm failed ", err)
	}
}

func TestAlterTables(t *testing.T) {

	if err := GenerateTables(false); err != nil {
		t.Error("insert by orm failed ", err)
	}
}

func TestORM(t *testing.T) {

	u := &models.User{Id: 1}
	token := &models.Token{
		Id:          1,
		TokenSymbol: "btc",
		TokenName:   "bitcoin",
		TokenIntro:  "NO. 1",
	}

	ut := models.UserToken{
		TokenAddress: "111",
		PrivateKey:   "pk",
		TokenBalance: 34,
		TokenExtra:   "test",
		User:         u,
		Token:        token,
	}

	if _, err := orm_.Insert(&ut); err != nil {
		t.Error("insert by orm failed ", err)
		return
	}

}

func TestConnect(t *testing.T) {
	if _, err := Connect(); err != nil {
		t.Error("connect to mysql failed ", err)
	}
}
func TestNewUser(t *testing.T) {
	u := models.User{
		UserName:     string(rand.Int()),
		UserPassword: string(rand.Int()),
		MailAddress:  string(rand.Int()),
		MailCode:     string(rand.Int()),
		PhoneNo:      string(rand.Int()),
		CountryCode:  string(rand.Int()),
		UserPortrait: string(rand.Int()),
		UserIntro:    string(rand.Int()),
	}

	res, err := orm_.Insert(&u)
	if err != nil {
		t.Error(err)
	}
	beego.Debug(res)

}

func TestNewToken(t *testing.T) {
	res := NewToken(enums.TOKEN_BITCOIN, "BTC", "bitcoin", "NO. 1")
	if res < 0 {
		t.Error("new bitcoin failed")
	}
	res = NewToken(enums.TOKEN_ETHEREUM, "ETH", "ethereum", "NO. 2")
	if res < 0 {
		t.Error("new bitcoin failed")
	}
}
func TestNewTokenByUser(t *testing.T) {
	ut := NewTokenByUser(1, enums.TOKEN_BITCOIN, "bitcoin address", "private key")
	if ut == nil {
		t.Error("get token by user failed ")
		return
	}
}

func TestGetTokenByUser(t *testing.T) {
	ut := GetTokenByUser(1, 1)
	if ut == nil {
		t.Error("get token by user failed ")
		return
	}
}
func TestGetTokenByUserEmpty(t *testing.T) {
	ut := GetTokenByUser(0, 0)
	if ut == nil {
		t.Error("get token by user failed ", ut)
	}
}

func TestMain(m *testing.M) {

	ORM()
	retCode := m.Run()
	os.Exit(retCode)
}