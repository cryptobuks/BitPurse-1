package enums

type TOKEN uint8

const (
	TOKEN_UNKNOWN  TOKEN = iota
	TOKEN_BITCOIN
	TOKEN_ETHEREUM
)

type TX uint64

const (
	TX_UNKNOWN TX = 1 << iota
	TX_DONE
	TX_SPENT
	TX_STORED
)

type OP uint8

const (
	OP_UNKNOWN = iota
	OP_SEND
	OP_RECEIVE
)
