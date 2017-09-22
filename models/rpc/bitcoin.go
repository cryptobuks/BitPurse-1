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
	"github.com/astaxie/beego/httplib"
	"strconv"
	"math/rand"
	//"reflect"
)

var (
	client_ *btcrpcclient.Client
)

func Init() {
	var err error
	client_, err = btcrpcclient.New(&btcrpcclient.ConnConfig{
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default

		Host: beego.AppConfig.String("BITCOIN_HOST"),
		User: beego.AppConfig.String("BITCOIN_USER"),
		Pass: beego.AppConfig.String("BITCOIN_PASS"),
	}, nil)
	if err != nil {
		fmt.Println("client_ init failed")
	}
}

func init() {

}

type BitcoinRpc struct {
	IRpc
}

func (_self *BitcoinRpc) Address2BalanceMap() map[string]float64 {

	result := _self.Call("listaddressgroupings", map[string]interface{}{})

	if r1, ok := result.([]interface{}); ok {
		m := make(map[string]float64)
		for _, v1 := range r1 {
			if r2, ok := v1.([]interface{}); ok {
				for _, v2 := range r2 {
					if v, ok := v2.([]interface{}); ok {
						addr := fmt.Sprintf("%s", v[0])
						amountStr := fmt.Sprintf("%f", v[1])
						if amount, err := strconv.ParseFloat(amountStr, 64); err == nil && amount > 0 {
							m[addr] = amount
						}
					}
				}
			}
		}
		return m

	} else {
		beego.Error("ListAddressGroupings conversion failed", result)
	}

	return nil
}

func (_self *BitcoinRpc) Call(method string, params map[string]interface{}) interface{} {
	host := beego.AppConfig.String("BITCOIN_HOST")
	url := fmt.Sprintf("http://%s", host)
	user := beego.AppConfig.String("BITCOIN_USER")
	pass := beego.AppConfig.String("BITCOIN_PASS")

	req := httplib.Post(url)
	req.SetBasicAuth(user, pass)

	id := rand.Int()
	reqData := map[string]interface{}{
		"jsonrpc": "1.0",
		"method":  method,
		"params":  params,
		"id":      id,
	}
	if jsReq, err := req.JSONBody(reqData); err == nil {

		type Results struct {
			Result interface{}            `json:"result"`
			Error  map[string]interface{} `json:"error"`
			ID     int                    `json:"id"`
		}

		var r Results
		if err := jsReq.ToJSON(&r); err == nil && r.ID == id {
			return r.Result
		} else {
			beego.Error(r.Error)
		}
	} else {
		beego.Error(err)
	}
	return nil
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
func parseAddress(_address string) *btcutil.Address {
	if address, err1 := btcutil.DecodeAddress(_address, configs.GetNetParams()); err1 == nil {
		return &address
	}
	beego.Error("Invalid address", _address)

	return nil
}

func (_self *BitcoinRpc) ListUnspentByAddress(_address string) []btcjson.ListUnspentResult {
	if address, err1 := btcutil.DecodeAddress(_address, configs.GetNetParams()); err1 == nil {
		unspent, err2 := client_.ListUnspentMinMaxAddresses(1, 9999999, []btcutil.Address{address})
		if err2 != nil {
			beego.Error(err1)
			return nil
		}

		return unspent
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

//  cannot calculate hot address since hot address may be change address
func (_self *BitcoinRpc) ListUnspent() []btcjson.ListUnspentResult {
	if r, err := client_.ListUnspent(); err == nil {
		return r
	} else {
		beego.Error(err)
	}
	return nil
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
	if len(_from) == 0 || len(_to) == 0 {
		beego.Error("[NewTx]No from or to", _from, _to)
		return ""
	}

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
	} else {
		beego.Error(err)
		return ""
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

	change := fromAmount - toAmount
	if change > 0 {
		if len(_changeAddress) > 0 {
			if address, err1 := btcutil.DecodeAddress(_changeAddress, configs.GetNetParams()); err1 == nil {
				if amount, err2 := btcutil.NewAmount(change); err2 == nil {
					amounts[address] = amount
				}
			}
		}
	} else if change < 0 {
		beego.Error("negative change", change)
		return ""
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

func (_self *BitcoinRpc) ValidateAddress(_address string) bool {

	if addr := parseAddress(_address); addr != nil {
		if r, err := client_.ValidateAddress(*addr); err == nil && r.IsValid {
			return true
		}
	}

	beego.Error("Invalid address", _address)
	return false
}
