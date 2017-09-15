package configs

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/astaxie/beego"
)

func GetNetParams() *chaincfg.Params {
	switch net := beego.AppConfig.String("NET_PARAMS"); net {
	case "REGRESSION":
		return &chaincfg.RegressionNetParams
	case "MAIN":
		return &chaincfg.MainNetParams
	case "TEST":
		return &chaincfg.TestNet3Params
	}

	return &chaincfg.MainNetParams
}
