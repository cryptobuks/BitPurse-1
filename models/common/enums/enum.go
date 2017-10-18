package enums

type TOKEN uint8

const (
	TOKEN_UNKNOWN  TOKEN = iota
	TOKEN_BITCOIN
	TOKEN_ETHEREUM
)

type TX uint64

const (
	TX_UNKNOWN TX = 1 << iota // just add empty tx to db, waiting for previous tx done
	TX_DONE                   // tx has at least 1 confirmation, fund in user account
	TX_SENT                   // tx sent to network with zero conformation
	TX_STORED                 // fund moved from user account to hot account
)

type OP uint8

const (
	OP_UNKNOWN = iota
	OP_SEND
	OP_RECEIVE
)
