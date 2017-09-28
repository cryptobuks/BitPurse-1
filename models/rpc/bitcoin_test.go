package rpc

import (
	"testing"
	"os"
	"path/filepath"
	"github.com/astaxie/beego"
)

func TestBitcoinRpc_Call(t *testing.T) {

	r := new(BitcoinRpc)

	type GetMemPoolInfoResult struct {
		Size          int     `json:"size"`
		Bytes         int     `json:"bytes"`
		Usage         int     `json:"usage"`
		MaxMemPool    int     `json:"maxmempool"`
		MemPoolMinFee float64 `json:"mempoolminfee"`
	}
	r1 := new(GetMemPoolInfoResult)
	if ok := r.Call("getmempoolinfo", []interface{}{}, r1); !ok {

		t.Error(ok)
	} else {
		t.Log(r1)
	}
}

func TestBitcoinRpc_Address2BalanceMap(t *testing.T) {
	r := new(BitcoinRpc)
	if m := r.Address2BalanceMap(); m == nil {
		t.Error()
	} else {
		t.Log(len(m))
	}
}

func TestBitcoinRpc_NewTx(t *testing.T) {
	r := new(BitcoinRpc)

	hot := beego.AppConfig.String("HOT_ADDRESS")
	cold := beego.AppConfig.String("COLD_ADDRESS")
	from := []string{hot}
	to := map[string]float64{cold: 10.1}
	to["n3ZMMByDb478XAgkZwvgbWwt25AqeGHBj1"] = 2

	if hex := r.NewTx(from, to, hot); hex == "" {
		t.Error()
	} else {
		t.Log(hex)
	}

}

func TestMain(m *testing.M) {

	if path, err := filepath.Abs("../.."); err == nil {
		beego.TestBeegoInit(path + "")

		Init()
	}
	code := m.Run()
	os.Exit(code)
}
