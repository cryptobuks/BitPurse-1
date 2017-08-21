package dao

import (
	"testing"
	"os"
)

func TestConnect(t *testing.T) {
	if _, err := Connect(); err != nil {
		t.Error("connect to mysql failed ", err)
	}
}

func TestGetTokenByUser(t *testing.T) {
	ut := GetTokenByUser(1, 1)
	if ut != nil {
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
	Connect()

	retCode := m.Run()
	os.Exit(retCode)
}
