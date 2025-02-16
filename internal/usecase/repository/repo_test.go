package repository

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/k1v4/avito_shop/internal/config"
	"github.com/k1v4/avito_shop/internal/entity"
	"github.com/k1v4/avito_shop/internal/usecase"
	"github.com/k1v4/avito_shop/pkg/DB/postgres"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"testing"
)

func TestLinksRepository(t *testing.T) {
	ctx := context.Background()

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)
	err = godotenv.Load("../../../.env") // Укажите полный путь к файлу .env
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	os.Setenv("POSTGRES_HOST", "localhost")

	cfg := config.MustLoadConfig()
	if cfg == nil {
		panic("load config fail")
	}

	urlDB := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBConfig.UserName,
		cfg.DBConfig.Password,
		cfg.DBConfig.Host,
		cfg.DBConfig.Port,
		cfg.DBConfig.DbName,
	)

	pg, err := postgres.New(urlDB, postgres.MaxPoolSize(cfg.DBConfig.PoolMax))
	if err != nil {
		panic(err)
	}

	username := "testUser"
	password := "testPassword"

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	linksRepository := NewShopRepository(pg)

	// SaveUser
	userSave, err := linksRepository.SaveUser(ctx, username, passHash)
	assert.NoError(t, err)
	assert.NotEmpty(t, userSave)
	assert.NotEqual(t, 0, userSave)

	// FindUser
	userFind, err := linksRepository.FindUser(ctx, username)
	assert.NoError(t, err)
	assert.Equal(t, entity.User{
		Id:       userSave,
		Username: username,
		Passhash: passHash,
		Coins:    1000,
	}, userFind)

	userFindErr, err := linksRepository.FindUser(ctx, "1")
	assert.Error(t, err)
	assert.Equal(t, entity.User{}, userFindErr)

	// BuyItem
	err = linksRepository.BuyItem(ctx, userSave, 1, 1)
	assert.NoError(t, err)

	err = linksRepository.BuyItem(ctx, -1, 1, 1)
	assert.Error(t, err)

	// GetItemUser
	userItems, err := linksRepository.GetItemUser(ctx, userSave)
	assert.NoError(t, err)
	assert.Equal(t, entity.Inventory{Items: []entity.InventoryItem{
		{
			ItemId:   1,
			Type:     "",
			Quantity: 1,
		},
	}}, userItems)

	userItemsErr, err := linksRepository.GetItemUser(ctx, -1)
	assert.Equal(t, entity.Inventory{}, userItemsErr)

	// GetItemByName
	item, err := linksRepository.GetItemByName(ctx, "cup")
	assert.NoError(t, err)
	assert.Equal(t, item.Id, 2)

	itemEmpty, err := linksRepository.GetItemByName(ctx, "hahahahha")
	assert.Error(t, err)
	assert.Equal(t, itemEmpty, entity.Item{})

	// GetItemById
	id, err := linksRepository.GetItemById(ctx, item.Id)
	assert.NoError(t, err)
	assert.Equal(t, id, item.Name)

	idEmpty, err := linksRepository.GetItemById(ctx, -1)
	assert.Error(t, err)
	assert.ErrorIs(t, err, usecase.ErrNoItem)
	assert.Equal(t, idEmpty, "")

	// GetUserById
	byId, err := linksRepository.GetUserById(ctx, userSave)
	assert.NoError(t, err)
	assert.Equal(t, byId.Id, userSave)

	// TakeGiveCoins
	err = linksRepository.TakeGiveCoins(ctx, userSave, 100)
	assert.NoError(t, err)

	err = linksRepository.TakeGiveCoins(ctx, -1, 1)

	// MakeRecord
	err = linksRepository.MakeRecord(ctx, userSave, 1, 100)
	assert.NoError(t, err)

	// TakeRecords
	records, err := linksRepository.TakeRecords(ctx, userSave)
	assert.NoError(t, err)
	assert.NotEmpty(t, records)
}
