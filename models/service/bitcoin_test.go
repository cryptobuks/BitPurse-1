package service

import (
	"testing"
	"github.com/Newtrong/BitPurse/models/common/enums"
	"path/filepath"
	"github.com/astaxie/beego"
	"os"
	"github.com/Newtrong/BitPurse/models/rpc"
	"time"
	"fmt"
)

func TestWatchNotify(t *testing.T) {

}

func TestBitcoinService_User2General(t *testing.T) {

	beego.Debug("TestBitcoinService_User2General")
	oldTime := beego.AppConfig.String("USER_2_GENERAL_TIME")
	t.Log(oldTime)

	now := time.Now()
	next := now.Add(2 * time.Second)
	newTime := fmt.Sprintf("%d:%d:%d", next.Hour(), next.Minute(), next.Second())
	beego.AppConfig.Set("USER_2_GENERAL_TIME", newTime)
	t.Log(beego.AppConfig.String("USER_2_GENERAL_TIME"))

	beego.AppConfig.Set("USER_2_GENERAL_INTERVAL", "3s")

	if oldTime == newTime {
		t.Error(newTime)
	}

	firstTimer := make(chan *time.Timer)
	firstTicker := make(chan *time.Ticker)
	delay := User2General(firstTimer, firstTicker)

	beego.Debug("waiting for timer")
	<-firstTimer
	beego.Debug("timer is on")

	beego.Debug("waiting for ticker")
	<-firstTicker
	beego.Debug("ticker is on")

	if delay > time.Minute {
		t.Error()
	}
	t.Log(delay)
}

func TestBitcoinService_CheckBalance(t *testing.T) {
	s := Get(enums.TOKEN_BITCOIN)
	if r := s.CheckBalance(); !r {
		t.Error()
	}
}

func TestBitcoinService_UnspentBalance(t *testing.T) {
	s := Get(enums.TOKEN_BITCOIN)
	if sum := s.UnspentBalance(); sum < 0 {
		t.Error()
	} else {
		t.Log(sum)
	}
}

func TestMain(m *testing.M) {
	if path, err := filepath.Abs("../.."); err == nil {
		beego.TestBeegoInit(path + "")

		rpc.Init()
		InitBitcoin()
	}

	retCode := m.Run()
	os.Exit(retCode)
}
