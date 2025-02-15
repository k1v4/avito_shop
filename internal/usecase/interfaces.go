package usecase

import (
	"context"
	"github.com/k1v4/avito_shop/internal/entity"
)

type (
	IShopRepository interface {
		SaveUser(ctx context.Context, username string, passhash []byte) (int, error)
		FindUser(ctx context.Context, username string) (entity.User, error)
	}

	IShopService interface {
		Login(ctx context.Context, username, password string) (string, error)
	}
)
