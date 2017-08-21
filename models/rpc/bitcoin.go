package rpc

import (
	"github.com/btcsuite/btcrpcclient"
	"fmt"
	"github.com/astaxie/beego"
	"git.coding.net/zhouhuangjing/BitPurse/models/common"
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

func (b *BitcoinRpc) NewAddress() common.TokenAddress {
	if address, err := client.GetNewAddress(""); err != nil {
		beego.Error(err)
		return common.TokenAddress("")
	} else {
		return common.TokenAddress(address.String())
	}
}

// 充币, 只是生成用户的代币地址
func (b *BitcoinRpc) Deposit() {

}

func (b *BitcoinRpc) Withdraw() {

}
