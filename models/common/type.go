package common

type ID int
type TokenAddress string

type UserToken struct {
	UserTokenID  ID
	UserId       ID
	TokenType    uint8
	// TokenAddress is not work on read from database
	TokenAddress string
	TokenBalance float64
	TokenExtra   string
}

