package usecase

import "errors"

var (
	ErrNoUser  = errors.New("user not found")
	ErrNoItem  = errors.New("item not found")
	ErrNoCoins = errors.New("not enough coins")
)
