package rpc

import "git.coding.net/zhouhuangjing/BitPurse/models/common"

type IRpc interface {
	NewAddress() common.TokenAddress
	Deposit()
	Withdraw()
}
