package tokens

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
		Host:         "127.0.0.1:18332",
		User:         "INVENT_A_UNIQUE_USERNAME",
		Pass:         "INVENT_A_UNIQUE_PASSWORD",
	}, nil)
	if err != nil {
		fmt.Println("client init failed")
	}
}

type BitcoinHelper struct {
	Helper
}

func (b *BitcoinHelper) NewAddress() string {
	// make GetNewAddress later
	if address, err := client.GetNewAddress(""); err != nil {
		beego.Error(err)
		return ""
	} else {
		return address.String()
	}
}

func (b *BitcoinHelper) Deposit() {

}

func (b *BitcoinHelper) Withdraw() {

}
