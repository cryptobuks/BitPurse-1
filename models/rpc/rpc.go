package rpc

type IRpc interface {
	NewAddress() (string, string)
	Deposit()
	Withdraw(_address string, _amount float64) string
	Watch(_address string)
}
