package tasks

import "github.com/ethereum/go-ethereum/core/types"

type Task interface {
	GetLambPriceFromMarket(secret []byte) error
	SetLambPriceToOracle(secret []byte) (*types.Transaction, error)
}
