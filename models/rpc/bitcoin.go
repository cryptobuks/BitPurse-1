package rpc

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcrpcclient"
	"github.com/btcsuite/btcutil"
	"git.coding.net/zhouhuangjing/BitPurse/models/common/configs"
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcd/wire"
)

var (
	client_ *btcrpcclient.Client
)

func init() {

	var err error
	client_, err = btcrpcclient.New(&btcrpcclient.ConnConfig{
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default

		Host: "127.0.0.1:19011",
		User: "admin2",
		Pass: "123",
	}, nil)
	if err != nil {
		fmt.Println("client_ init failed")
	}
}

type BitcoinRpc struct {
	IRpc
}

func (_self *BitcoinRpc) parseTx(_tx string) *wire.MsgTx {
	if serializedTx, err := hex.DecodeString(_tx); err == nil {
		var msgTx wire.MsgTx
		if err := msgTx.Deserialize(bytes.NewReader(serializedTx)); err == nil {
			return &msgTx
		}
	}
	return nil
}

func (_self *BitcoinRpc) serializeTx(_tx *wire.MsgTx) string {
	if _tx != nil {
		// Serialize the transaction and convert to hex string.
		buf := bytes.NewBuffer(make([]byte, 0, _tx.SerializeSize()))
		if err := _tx.Serialize(buf); err == nil {
			txHex := hex.EncodeToString(buf.Bytes())
			return txHex
		}
	}
	return ""
}

func (_self *BitcoinRpc) ListUnspentByAddress(_address string) *[]btcjson.ListUnspentResult {
	if address, err1 := btcutil.DecodeAddress(_address, configs.GetNetParams()); err1 == nil {
		unspent, err2 := client_.ListUnspentMinMaxAddresses(1, 9999999, []btcutil.Address{address})
		if err2 != nil {
			beego.Error(err1)
			return nil
		}

		return &unspent
	}

	return nil
}

func (_self *BitcoinRpc) Balance() float64 {
	r, err1 := client_.GetBalance("")
	if err1 != nil {
		beego.Error(err1)
		return -1
	}
	return r.ToBTC()
}

func (_self *BitcoinRpc) ListUnspent() *[]btcjson.ListUnspentResult {
	r, err1 := client_.ListUnspent()
	if err1 != nil {
		beego.Error(err1)
		return nil
	}
	return &r
}

func (_self *BitcoinRpc) Transfer(_from []string, _to map[string]float64, _changeAddress string) string {
	if tx := _self.NewTx(_from, _to, _changeAddress); tx != "" {
		if tx, complete := _self.SignTx(tx); tx != "" && complete {
			if tx := _self.SendTx(tx); tx != "" {
				return tx
			}
		}
	}
	beego.Error("Transfer failed", _from, _to, _changeAddress)

	return ""
}

func (_self *BitcoinRpc) SendTx(_tx string) string {
	tx := _self.parseTx(_tx)
	if signedTx, err := client_.SendRawTransaction(tx, true); err == nil {
		return signedTx.String()
	}

	return ""
}
func (_self *BitcoinRpc) SignTx(_tx string) (string, bool) {
	tx := _self.parseTx(_tx)
	if signedTx, complete, err := client_.SignRawTransaction(tx); err == nil {
		return _self.serializeTx(signedTx), complete
	}

	return "", false
}

func (_self *BitcoinRpc) NewTx(_from []string, _to map[string]float64, _changeAddress string) string {

	addresses := make([]btcutil.Address, 0)
	for _, v := range _from {
		if address, err1 := btcutil.DecodeAddress(v, configs.GetNetParams()); err1 == nil {
			addresses = append(addresses, address)
		}
	}

	inputs := make([]btcjson.TransactionInput, 0)
	fromAmount := 0.0
	if unspent, err := client_.ListUnspentMinMaxAddresses(1, 99999, addresses); err == nil {
		for _, v := range unspent {
			fromAmount += v.Amount
			inputs = append(inputs, btcjson.TransactionInput{Txid: v.TxID, Vout: v.Vout})
		}
	}

	toAmount := 0.0
	amounts := make(map[btcutil.Address]btcutil.Amount)
	for to, a := range _to {
		toAmount += a
		if address, err1 := btcutil.DecodeAddress(to, configs.GetNetParams()); err1 == nil {
			if amount, err2 := btcutil.NewAmount(a); err2 == nil {
				amounts[address] = amount
			}
		}
	}

	change := toAmount - toAmount
	if change > 0 && len(_changeAddress) > 0 {
		if address, err1 := btcutil.DecodeAddress(_changeAddress, configs.GetNetParams()); err1 == nil {
			if amount, err2 := btcutil.NewAmount(change); err2 == nil {
				amounts[address] = amount
			}
		}
	}

	var lockTime int64 = 0
	if tx, err3 := client_.CreateRawTransaction(inputs, amounts, &lockTime); err3 == nil {
		return _self.serializeTx(tx)
	}

	return ""
}

// 要记录公钥和私钥
func (_self *BitcoinRpc) NewAddress() (string, string) {
	if address, err := client_.GetNewAddress(""); err != nil {
		beego.Error(err)
		return "", ""
	} else {
		pk, err := client_.DumpPrivKey(address)
		if err != nil {
			beego.Error(err)
			return "", ""
		}

		return address.String(), pk.String()
	}
}

func (_self *BitcoinRpc) GetTransaction(_txId string) *btcjson.GetTransactionResult {
	h, err1 := chainhash.NewHashFromStr(_txId)
	if err1 == nil {
		tx, err2 := client_.GetTransaction(h)
		if err2 == nil {
			return tx
		}
	}

	return nil
}

func (_self *BitcoinRpc) Watch(_address string) {
	//address, err := btcutil.DecodeAddress(_address, &chaincfg.RegressionNetParams)
	//if err != nil {
	//	beego.Error(err)
	//	return
	//}
	// need web socket, try walletnotify instead
	//err1 := client_.NotifyReceived([]btcutil.Address{address})
	//if err1 != nil {
	//	beego.Error(err1)
	//}
}
