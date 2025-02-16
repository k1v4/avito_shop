package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/k1v4/avito_shop/internal/entity"
	"github.com/k1v4/avito_shop/pkg/jwtPkg"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrWrongPassword      = errors.New("wrong password")
	ErrUserExist          = errors.New("user exist")
)

type ShopUseCase struct {
	repo  IShopRepository
	cache *redis.Client
}

func NewShopUseCase(r IShopRepository, red *redis.Client) *ShopUseCase {
	return &ShopUseCase{
		repo:  r,
		cache: red,
	}
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

func (uc *ShopUseCase) BuyItem(ctx context.Context, userId int, itemName string) error {
	const op = "ShopUseCase.BuyItem"

	item, err := uc.repo.GetItemByName(ctx, itemName)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	user, err := uc.repo.GetUserById(ctx, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if user.Coins < item.Price {
		return ErrNoCoins
	}

	err = uc.repo.BuyItem(ctx, userId, item.Id, 1)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = uc.repo.TakeGiveCoins(ctx, userId, -item.Price)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (uc *ShopUseCase) SendCoins(ctx context.Context, toUserName string, fromUserId, amount int) error {
	const op = "ShopUseCase.SendCoins"

	toUser, err := uc.repo.FindUser(ctx, toUserName)
	if err != nil {
		return ErrNoUser
	}

	fromUser, err := uc.repo.GetUserById(ctx, fromUserId)
	if err != nil {
		return ErrNoUser
	}

	if fromUser.Coins < amount {
		return ErrNoCoins
	}

	toUserId := toUser.Id

	err = uc.repo.TakeGiveCoins(ctx, toUserId, amount)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = uc.repo.TakeGiveCoins(ctx, fromUserId, -amount)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = uc.repo.MakeRecord(ctx, fromUserId, toUserId, amount)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (uc *ShopUseCase) GetInfo(ctx context.Context, userId int) (entity.ResponseInfo, error) {
	const op = "ShopUseCase.GetInfo"
	var res entity.ResponseInfo

	err := uc.cache.Get(ctx, fmt.Sprintf("%d", userId)).Scan(&res)
	if err != nil && !errors.Is(err, redis.Nil) {
		return res, fmt.Errorf("%s: %w", op, err)
	}

	// Coins
	user, err := uc.repo.GetUserById(ctx, userId)
	if err != nil {
		return entity.ResponseInfo{}, fmt.Errorf("%s: %w", op, err)
	}

	res.Coins = user.Coins

	itemsUser, err := uc.repo.GetItemUser(ctx, userId)
	if err != nil {
		return entity.ResponseInfo{}, fmt.Errorf("%s: %w", op, err)
	}

	for i, _ := range itemsUser.Items {
		name, err := uc.repo.GetItemById(ctx, itemsUser.Items[i].ItemId)
		if err != nil {
			return entity.ResponseInfo{}, fmt.Errorf("%s: %w", op, err)
		}

		itemsUser.Items[i].Type = name
	}

	res.Inventory = itemsUser

	bothDir, err := uc.repo.TakeRecords(ctx, userId)
	if err != nil {
		return entity.ResponseInfo{}, fmt.Errorf("%s: %w", op, err)
	}

	var sentItems []entity.SentItem
	var receivedItems []entity.ReceivedItem

	for _, item := range bothDir {
		var s entity.SentItem
		var r entity.ReceivedItem

		if item.FromUser == userId {
			u, err := uc.repo.GetUserById(ctx, item.ToUser)
			if err != nil {
				return entity.ResponseInfo{}, fmt.Errorf("%s: %w", op, err)
			}

			r.Amount = item.Amount
			r.FromUser = u.Username

			receivedItems = append(receivedItems, r)
		} else {
			u, err := uc.repo.GetUserById(ctx, item.FromUser)
			if err != nil {
				return entity.ResponseInfo{}, fmt.Errorf("%s: %w", op, err)
			}

			s.Amount = item.Amount
			s.ToUser = u.Username

			sentItems = append(sentItems, s)
		}
	}

	coinHistory := entity.CoinHistory{
		Received: entity.Received{ReceivedItems: receivedItems},
		Sent:     entity.Sent{SentItems: sentItems},
	}

	res.CoinHistory = coinHistory

	uc.cache.Set(context.Background(), fmt.Sprintf("%d", userId), res, 1*time.Minute)

	return res, nil
}
