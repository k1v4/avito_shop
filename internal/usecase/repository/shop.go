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
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	var id int
	err = s.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
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
	err = s.Pool.QueryRow(ctx, sq, args...).Scan(&user.Id, &user.Username, &user.Passhash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, usecase.ErrNoUser
		}

		return entity.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
