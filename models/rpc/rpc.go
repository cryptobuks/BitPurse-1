package rpc


type IRpc interface {
	NewAddress() (string, string)
	Deposit()
	Withdraw()
}
