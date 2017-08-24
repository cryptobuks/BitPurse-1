package rpc

import (
	"github.com/btcsuite/btcrpcclient"
	"fmt"
	"github.com/astaxie/beego"
)

var (
	client *btcrpcclient.Client
)

func init() {

	var err error
	client, err = btcrpcclient.New(&btcrpcclient.ConnConfig{
		HTTPPostMode: true,
		DisableTLS:   true,
		Host:         "127.0.0.1:19011",
		User:         "admin2",
		Pass:         "123",
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

func (b *BitcoinRpc) Withdraw() {

}
