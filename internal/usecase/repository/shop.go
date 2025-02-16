package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/k1v4/avito_shop/internal/entity"
	"github.com/k1v4/avito_shop/internal/usecase"
	"github.com/k1v4/avito_shop/pkg/DB/postgres"
)

const defaultEntityCap = 64

type ShopRepository struct {
	*postgres.Postgres
}

func NewShopRepository(pg *postgres.Postgres) *ShopRepository {
	return &ShopRepository{
		Postgres: pg,
	}
}

func (s *ShopRepository) SaveUser(ctx context.Context, username string, passhash []byte) (int, error) {
	const op = "ShopRepository.SaveUser"

	sql, args, err := s.Builder.Insert("users").
		Columns("username", "password").
		Values(username, passhash).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var id int
	err = s.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *ShopRepository) FindUser(ctx context.Context, username string) (entity.User, error) {
	const op = "ShopRepository.FindUser"

	sq, args, err := s.Builder.Select("*").
		From("users").
		Where(squirrel.Eq{"username": username}).
		ToSql()
	if err != nil {
		return entity.User{}, fmt.Errorf("%s: %w", op, err)
	}

	var user entity.User
	err = s.Pool.QueryRow(ctx, sq, args...).Scan(&user.Id, &user.Username, &user.Passhash, &user.Coins)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, usecase.ErrNoUser
		}

		return entity.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *ShopRepository) BuyItem(ctx context.Context, userId, itemId, quantity int) error {
	const op = "ShopRepository.BuyItem"

	sq, args, err := s.Builder.Insert("inventory").
		Columns("user_id", "item_id", "quantity").
		Values(userId, itemId, quantity).
		Suffix("ON CONFLICT (user_id, item_id) DO UPDATE SET quantity = inventory.quantity + EXCLUDED.quantity").
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.Pool.Exec(ctx, sq, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *ShopRepository) GetItemUser(ctx context.Context, userId int) (entity.Inventory, error) {
	const op = "ShopRepository.GetItemUser"

	sq, args, err := s.Builder.Select("item_id", "quantity").
		From("inventory").
		Where(squirrel.Eq{"user_id": userId}).
		ToSql()
	if err != nil {
		return entity.Inventory{}, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := s.Pool.Query(ctx, sq, args...)
	if err != nil {
		return entity.Inventory{}, fmt.Errorf("%s: %w", op, err)
	}

	var inventory []entity.InventoryItem
	for rows.Next() {
		var item entity.InventoryItem
		err = rows.Scan(&item.ItemId, &item.Quantity)
		if err != nil {
			return entity.Inventory{}, fmt.Errorf("%s: %w", op, err)
		}

		inventory = append(inventory, item)
	}

	return entity.Inventory{
		Items: inventory,
	}, nil
}

func (s *ShopRepository) GetItemByName(ctx context.Context, itemName string) (entity.Item, error) {
	const op = "ShopRepository.GetItemByName"

	sq, args, err := s.Builder.Select("*").
		From("items").
		Where(squirrel.Eq{"name": itemName}).
		ToSql()
	if err != nil {
		return entity.Item{}, fmt.Errorf("%s: %w", op, err)
	}

	var item entity.Item
	err = s.Pool.QueryRow(ctx, sq, args...).Scan(&item.Id, &item.Name, &item.Price)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Item{}, usecase.ErrNoItem
		}

		return entity.Item{}, fmt.Errorf("%s: %w", op, err)
	}

	return item, nil
}

func (s *ShopRepository) GetItemById(ctx context.Context, itemId int) (string, error) {
	const op = "ShopRepository.GetItemById"

	sq, args, err := s.Builder.Select("name").
		From("items").
		Where(squirrel.Eq{"id": itemId}).
		ToSql()
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var res string
	err = s.Pool.QueryRow(ctx, sq, args...).Scan(&res)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", usecase.ErrNoItem
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}

func (s *ShopRepository) GetUserById(ctx context.Context, userId int) (entity.User, error) {
	const op = "ShopRepository.GetCoins"

	sq, args, err := s.Builder.
		Select("*").
		From("users").
		Where(squirrel.Eq{"id": userId}).
		ToSql()
	if err != nil {
		return entity.User{}, fmt.Errorf("%s: %w", op, err)
	}

	var user entity.User
	err = s.Pool.QueryRow(ctx, sq, args...).Scan(&user.Id, &user.Username, &user.Passhash, &user.Coins)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, usecase.ErrNoUser
		}

		return entity.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// TakeGiveCoins надо передавать значение amount со знаком согласно операции (добавить: +, убрать: - )
func (s *ShopRepository) TakeGiveCoins(ctx context.Context, userId, amount int) error {
	const op = "ShopRepository.TakeGiveCoins"

	sq, args, err := s.Builder.
		Update("users").
		Set("amount", squirrel.Expr("amount + ?", amount)).
		Where(squirrel.Eq{"id": userId}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.Pool.Exec(ctx, sq, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *ShopRepository) MakeRecord(ctx context.Context, fromUserId, toUserId, amount int) error {
	const op = "ShopRepository.MakeRecord"

	sq, args, err := s.Builder.Insert("coin_history").
		Columns("from_user", "to_user", "amount").
		Values(fromUserId, toUserId, amount).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.Pool.Exec(ctx, sq, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *ShopRepository) TakeRecords(ctx context.Context, userId int) ([]entity.BothDirection, error) {
	const op = "ShopRepository.TakeRecords"

	sq, args, err := s.Builder.Select("from_user", "to_user", "amount").
		From("coin_history").
		Where(squirrel.Or{
			squirrel.Eq{"from_user": userId},
			squirrel.Eq{"to_user": userId},
		}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := s.Pool.Query(ctx, sq, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var items []entity.BothDirection
	for rows.Next() {
		var both entity.BothDirection

		err = rows.Scan(&both.ToUser, &both.FromUser, &both.Amount)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		items = append(items, both)
	}

	return items, nil
}
