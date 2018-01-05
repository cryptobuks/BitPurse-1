package rpc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Newtrong/BitPurse/models/common/configs"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"math/rand"
)

var (
	client_ *rpcclient.Client
)

func Init() {
	var err error
	client_, err = rpcclient.New(&rpcclient.ConnConfig{
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

	r := new([][][]interface{})
	if ok := _self.Call("listaddressgroupings", []interface{}{}, r); ok {

		m := make(map[string]float64)
		for _, v1 := range *r {
			for _, v2 := range v1 {
				if addr, ok := v2[0].(string); ok {
					if amount, ok := v2[1].(float64); ok {
						m[addr] = amount
					}
				}
			}
		}
		return m
	}

	beego.Error("[Address2BalanceMap]Failed")

	return nil
}

func (_self *BitcoinRpc) Call(method string, params []interface{}, v interface{}) bool {
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

		// ReadAll will disable ToJSON since the buffer is empty !!!!!
		//body, _ := ioutil.ReadAll(jsReq.GetRequest().Body)
		//beego.Debug(string(body))

		var r Results
		if err := jsReq.ToJSON(&r); err == nil && r.ID == id {
			buf := new(bytes.Buffer)

			encoder := json.NewEncoder(buf)
			if err := encoder.Encode(r.Result); err == nil {
				if err := json.Unmarshal(buf.Bytes(), v); err == nil {
					return true
				} else {
					beego.Error(err)
				}
			} else {
				beego.Error(err)
			}
		} else {
			beego.Error(r.Error)
		}
	} else {
		beego.Error(err)
	}
	return false
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
		if address := parseAddress(v); address != nil {
			addresses = append(addresses, *address)
		}
	}

	toAmount := 0.0
	amounts := make(map[btcutil.Address]btcutil.Amount)
	for to, a := range _to {
		toAmount += a
		if address := parseAddress(to); address != nil {
			if amount, err2 := btcutil.NewAmount(a); err2 == nil {
				amounts[*address] = amount
			}

		}
	}

	inputs := make([]btcjson.TransactionInput, 0)
	fromAmount := 0.0
	if unspent, err := client_.ListUnspentMinMaxAddresses(1, 99999, addresses); err == nil {
		for _, v := range unspent {
			fromAmount += v.Amount
			inputs = append(inputs, btcjson.TransactionInput{Txid: v.TxID, Vout: v.Vout})
			if fromAmount >= toAmount {
				break
			}
		}
	} else {
		beego.Error(err)
		return ""
	}
	if len(inputs) == 0 {
		beego.Error("[NewTx]No inputs")
		return ""
	}

	var lockTime int64 = 0
	if tx, err3 := client_.CreateRawTransaction(inputs, amounts, &lockTime); err3 == nil {
		rawTx := _self.serializeTx(tx)

		type FundRawTransactionOptions struct {
			ChangeAddress          string `json:"changeAddress"`
			ChangePosition         int    `json:"changePosition"`
			SubtractFeeFromOutputs []int  `json:"subtractFeeFromOutputs"`
		}

		fee := []int{}
		for i := 0; i < len(amounts); i++ {
			fee = append(fee, i)
		}

		o := FundRawTransactionOptions{
			ChangeAddress:          _changeAddress,
			ChangePosition:         len(amounts),
			SubtractFeeFromOutputs: fee,
		}

		type FundRawTransactionResult struct {
			Hex       string  `json:"hex"`
			Fee       float64 `json:"fee"`
			ChangePos int     `json:"changepos"`
		}

		r := new(FundRawTransactionResult)
		if ok := _self.Call("fundrawtransaction", []interface{}{rawTx, o}, r); ok {
			return r.Hex
		}
	} else {
		beego.Error("[NewTx]CreateRawTransaction", err3)
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
