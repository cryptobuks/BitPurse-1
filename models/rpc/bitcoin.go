package rpc

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcrpcclient"
	"github.com/btcsuite/btcutil"
)

var (
	client *btcrpcclient.Client
)

func init() {

	var err error
	client, err = btcrpcclient.New(&btcrpcclient.ConnConfig{
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default

		Host: "127.0.0.1:19011",
		User: "admin2",
		Pass: "123",
	}, nil)
	if err != nil {
		fmt.Println("client init failed")
	}
}

type BitcoinRpc struct {
	IRpc
}

// 要记录公钥和私钥
func (b *BitcoinRpc) NewAddress() (string, string) {
	if address, err := client.GetNewAddress(""); err != nil {
		beego.Error(err)
		return "", ""
	} else {
		pk, err := client.DumpPrivKey(address)
		if err != nil {
			beego.Error(err)
			return "", ""
		}

		return address.String(), pk.String()
	}
}

// 充币, 只是生成用户的代币地址
func (b *BitcoinRpc) Deposit() {

}

func (b *BitcoinRpc) Withdraw(_address string, _amount float64) string {
	amount, err1 := btcutil.NewAmount(_amount)
	address, err2 := btcutil.DecodeAddress(_address, &chaincfg.RegressionNetParams)
	if err1 != nil || err2 != nil {
		return ""
	}
	if hash, err3 := client.SendToAddress(address, amount); err3 == nil {
		beego.Debug(hash.String())
		return hash.String()
	}
	return ""
}

func (b *BitcoinRpc) GetTransaction(_txId string) *btcjson.GetTransactionResult {
	h, err1 := chainhash.NewHashFromStr(_txId)
	if err1 == nil {
		tx, err2 := client.GetTransaction(h)
		if err2 == nil {
			return tx
		}
	}

	return nil
}

func (b *BitcoinRpc) Watch(_address string) {
	//address, err := btcutil.DecodeAddress(_address, &chaincfg.RegressionNetParams)
	//if err != nil {
	//	beego.Error(err)
	//	return
	//}
	// need web socket, try walletnotify instead
	//err1 := client.NotifyReceived([]btcutil.Address{address})
	//if err1 != nil {
	//	beego.Error(err1)
	//}
}
