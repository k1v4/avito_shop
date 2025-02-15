package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/k1v4/avito_shop/internal/entity"
	"github.com/k1v4/avito_shop/pkg/jwtPkg"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrWrongPassword      = errors.New("wrong password")
	ErrInvalidAppId       = errors.New("invalid app id")
	ErrUserExist          = errors.New("user exist")
)

type ShopUseCase struct {
	repo IShopRepository
}

func NewShopUseCase(r IShopRepository) *ShopUseCase {
	return &ShopUseCase{repo: r}
}

func (uc *ShopUseCase) Login(ctx context.Context, username, password string) (string, error) {
	const op = "ShopUseCase.Login"

	user, err := uc.repo.FindUser(ctx, username)
	if err != nil {
		if errors.Is(err, ErrNoUser) {
			token, err := uc.Register(ctx, username, password)
			if err != nil {
				return "", fmt.Errorf("%s: %w", op, err)
			}
			return token, nil
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err = bcrypt.CompareHashAndPassword(user.Passhash, []byte(password)); err != nil {
		return "", ErrWrongPassword
	}

	token, err := jwtPkg.NewToken(user, time.Duration(1*time.Hour))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (uc *ShopUseCase) Register(ctx context.Context, username, password string) (string, error) {
	const op = "ShopUseCase.Register"

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	saveUserId, err := uc.repo.SaveUser(ctx, username, passHash)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := jwtPkg.NewToken(entity.User{
		Id:       saveUserId,
		Username: username,
		Passhash: passHash,
	}, time.Duration(1*time.Hour))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}
