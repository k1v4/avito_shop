package usecase

import (
	"context"
	"github.com/k1v4/avito_shop/internal/entity"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=IShopRepository
type IShopRepository interface {
	SaveUser(ctx context.Context, username string, passhash []byte) (int, error)
	FindUser(ctx context.Context, username string) (entity.User, error)
	BuyItem(ctx context.Context, userId, itemId, quantity int) error
	GetItemUser(ctx context.Context, userId int) (entity.Inventory, error)
	GetItemByName(ctx context.Context, itemId string) (entity.Item, error)
	GetItemById(ctx context.Context, itemId int) (string, error)
	GetUserById(ctx context.Context, userId int) (entity.User, error)
	TakeGiveCoins(ctx context.Context, userId, amount int) error
	MakeRecord(ctx context.Context, fromUserId, toUserId, amount int) error
	TakeRecords(ctx context.Context, userId int) ([]entity.BothDirection, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=IShopService
type IShopService interface {
	Login(ctx context.Context, username, password string) (string, error)
	Register(ctx context.Context, username, password string) (string, error)
	BuyItem(ctx context.Context, userId int, itemName string) error
	SendCoins(ctx context.Context, toUserName string, fromUserId, amount int) error
	GetInfo(ctx context.Context, userId int) (entity.ResponseInfo, error)
}
