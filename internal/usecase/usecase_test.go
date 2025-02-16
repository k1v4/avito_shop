package usecase

import (
	"context"
	"errors"
	"github.com/k1v4/avito_shop/internal/entity"
	"github.com/k1v4/avito_shop/internal/usecase/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"testing"
)

func TestLogin_Register(t *testing.T) {
	cases := []struct {
		name      string
		username  string
		password  string
		mockUser  entity.User
		mockErr   error
		mockToken string
		wantToken string
		wantErr   bool
	}{
		{
			name:      "success",
			username:  "user1",
			password:  "hashed_password",
			mockUser:  entity.User{Id: 1, Username: "user1", Passhash: []byte("hashed_password")},
			mockErr:   nil,
			mockToken: mock.Anything,
			wantToken: mock.Anything,
			wantErr:   false,
		},
		{
			name:      "wrong_password",
			username:  "user1",
			password:  "wrongpass",
			mockUser:  entity.User{Id: 1, Username: "user1", Passhash: []byte("hashed_password")},
			mockErr:   ErrWrongPassword,
			mockToken: "",
			wantToken: "",
			wantErr:   true,
		},
		{
			name:      "user_not_found_register_success",
			username:  "newuser",
			password:  "newpass",
			mockUser:  entity.User{},
			mockErr:   ErrNoUser,
			mockToken: mock.Anything,
			wantToken: mock.Anything,
			wantErr:   false,
		},
		{
			name:      "wrong_password",
			username:  "user1",
			password:  "wrongpass",
			mockUser:  entity.User{Id: 1, Username: "user1", Passhash: []byte("hashed_password")},
			mockErr:   nil,
			mockToken: "",
			wantToken: "",
			wantErr:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			if !errors.Is(tc.mockErr, ErrWrongPassword) {
				passHash, _ := bcrypt.GenerateFromPassword(tc.mockUser.Passhash, bcrypt.DefaultCost)
				tc.mockUser.Passhash = passHash
			}

			mockRepo := new(mocks.IShopRepository)
			mockCache := new(redis.Client)
			uc := NewShopUseCase(mockRepo, mockCache)

			mockRepo.
				On("FindUser", mock.Anything, tc.username).
				Return(tc.mockUser, tc.mockErr)

			if errors.Is(tc.mockErr, ErrNoUser) {
				mockRepo.
					On("SaveUser", mock.Anything, tc.username, mock.Anything).
					Return(1, nil)
			}

			token, err := uc.Login(context.Background(), tc.username, tc.password)

			if (err != nil) != tc.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if len(tc.wantToken) != 0 && len(strings.TrimSpace(token)) == 0 {
				t.Errorf("Login() token = %v, want not nill", token)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestBuyItem(t *testing.T) {
	cases := []struct {
		name     string
		userId   int
		itemName string
		mockItem entity.Item
		mockUser entity.User
		mockErr  error
		wantErr  bool
	}{
		{
			name:     "success",
			userId:   1,
			itemName: "item1",
			mockItem: entity.Item{Id: 1, Name: "item1", Price: 100},
			mockUser: entity.User{Id: 1, Coins: 200},
			mockErr:  nil,
			wantErr:  false,
		},
		{
			name:     "not_enough_coins",
			userId:   1,
			itemName: "item1",
			mockItem: entity.Item{Id: 1, Name: "item1", Price: 100},
			mockUser: entity.User{Id: 1, Coins: 50},
			mockErr:  ErrNoCoins,
			wantErr:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.IShopRepository)
			mockCache := new(redis.Client)
			uc := NewShopUseCase(mockRepo, mockCache)

			mockRepo.
				On("GetItemByName", mock.Anything, tc.itemName).
				Return(tc.mockItem, nil)
			mockRepo.
				On("GetUserById", mock.Anything, tc.userId).
				Return(tc.mockUser, nil)
			if tc.mockErr == nil {
				mockRepo.
					On("BuyItem", mock.Anything, tc.userId, tc.mockItem.Id, 1).
					Return(nil)

				mockRepo.
					On("TakeGiveCoins", mock.Anything, tc.userId, -tc.mockItem.Price).
					Return(nil)
			}

			err := uc.BuyItem(context.Background(), tc.userId, tc.itemName)

			if (err != nil) != tc.wantErr {
				t.Errorf("BuyItem() error = %v, wantErr %v", err, tc.wantErr)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSendCoins(t *testing.T) {
	cases := []struct {
		name       string
		fromUserId int
		toUserName string
		amount     int
		mockFrom   entity.User
		mockTo     entity.User
		mockErr    error
		wantErr    bool
	}{
		{
			name:       "success",
			fromUserId: 1,
			toUserName: "user2",
			amount:     100,
			mockFrom:   entity.User{Id: 1, Coins: 200},
			mockTo:     entity.User{Id: 2, Username: "user2"},
			mockErr:    nil,
			wantErr:    false,
		},
		{
			name:       "not_enough_coins",
			fromUserId: 1,
			toUserName: "user2",
			amount:     300,
			mockFrom:   entity.User{Id: 1, Coins: 200},
			mockTo:     entity.User{Id: 2, Username: "user2"},
			mockErr:    ErrNoCoins,
			wantErr:    true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.IShopRepository)
			mockCache := new(redis.Client)
			uc := NewShopUseCase(mockRepo, mockCache)

			mockRepo.
				On("FindUser", mock.Anything, tc.toUserName).
				Return(tc.mockTo, nil)

			mockRepo.
				On("GetUserById", mock.Anything, tc.fromUserId).
				Return(tc.mockFrom, nil)

			if tc.mockErr == nil {
				mockRepo.
					On("TakeGiveCoins", mock.Anything, tc.mockTo.Id, tc.amount).
					Return(nil)

				mockRepo.
					On("TakeGiveCoins", mock.Anything, tc.fromUserId, -tc.amount).
					Return(nil)

				mockRepo.
					On("MakeRecord", mock.Anything, tc.fromUserId, tc.mockTo.Id, tc.amount).
					Return(nil)
			}

			err := uc.SendCoins(context.Background(), tc.toUserName, tc.fromUserId, tc.amount)

			if (err != nil) != tc.wantErr {
				t.Errorf("SendCoins() error = %v, wantErr %v", err, tc.wantErr)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

//func TestGetInfo(t *testing.T) {
//	cases := []struct {
//		name      string
//		userId    int
//		mockUser  entity.User
//		mockItems entity.Inventory
//		mockRec   []entity.BothDirection
//		mockErr   error
//		wantRes   entity.ResponseInfo
//		wantErr   bool
//	}{
//		{
//			name:   "success",
//			userId: 1,
//			mockUser: entity.User{
//				Id:       1,
//				Username: "user1",
//				Coins:    100,
//			},
//			mockItems: entity.Inventory{
//				Items: []entity.InventoryItem{
//					{ItemId: 1, Type: "item1", Quantity: 2},
//				},
//			},
//			mockRec: []entity.BothDirection{
//				{FromUser: 1, ToUser: 2, Amount: 50},
//				{FromUser: 3, ToUser: 1, Amount: 30},
//			},
//			mockErr: nil,
//			wantRes: entity.ResponseInfo{
//				Coins: 100,
//				Inventory: entity.Inventory{
//					Items: []entity.InventoryItem{
//						{ItemId: 1, Type: "item1", Quantity: 2},
//					},
//				},
//				CoinHistory: entity.CoinHistory{
//					Received: entity.Received{
//						ReceivedItems: []entity.ReceivedItem{
//							{FromUser: "user3", Amount: 30},
//						},
//					},
//					Sent: entity.Sent{
//						SentItems: []entity.SentItem{
//							{ToUser: "user2", Amount: 50},
//						},
//					},
//				},
//			},
//			wantErr: false,
//		},
//		{
//			name:      "user_not_found",
//			userId:    1,
//			mockUser:  entity.User{},
//			mockItems: entity.Inventory{},
//			mockRec:   []entity.BothDirection{},
//			mockErr:   errors.New("user not found"),
//			wantRes:   entity.ResponseInfo{},
//			wantErr:   true,
//		},
//		{
//			name:   "cache_hit",
//			userId: 1,
//			mockUser: entity.User{
//				Id:       1,
//				Username: "user1",
//				Coins:    100,
//			},
//			mockItems: entity.Inventory{},
//			mockRec:   []entity.BothDirection{},
//			mockErr:   nil,
//			wantRes: entity.ResponseInfo{
//				Coins: 100,
//				Inventory: entity.Inventory{
//					Items: []entity.InventoryItem{},
//				},
//				CoinHistory: entity.CoinHistory{
//					Received: entity.Received{
//						ReceivedItems: []entity.ReceivedItem{},
//					},
//					Sent: entity.Sent{
//						SentItems: []entity.SentItem{},
//					},
//				},
//			},
//			wantErr: false,
//		},
//	}
//
//	for _, tc := range cases {
//		t.Run(tc.name, func(t *testing.T) {
//			mockRepo := new(mocks.IShopRepository)
//			redisClient, redisMock := redismock.NewClientMock()
//
//			uc := NewShopUseCase(mockRepo, redisClient)
//
//			// Mock Redis cache
//			if tc.name == "cache_hit" {
//				cachedData, _ := json.Marshal(tc.wantRes)
//				redisMock.ExpectGet(fmt.Sprintf("%d", tc.userId)).SetVal(string(cachedData))
//			} else {
//				redisMock.ExpectGet(fmt.Sprintf("%d", tc.userId)).SetErr(redis.Nil)
//				redisMock.ExpectSet(fmt.Sprintf("%d", tc.userId), mock.Anything, 1*time.Minute).SetVal("OK")
//			}
//
//			// Mock repository calls
//			mockRepo.On("GetUserById", mock.Anything, tc.userId).Return(tc.mockUser, tc.mockErr)
//			if tc.mockErr == nil {
//				mockRepo.On("GetItemUser", mock.Anything, tc.userId).Return(tc.mockItems, nil)
//				mockRepo.On("TakeRecords", mock.Anything, tc.userId).Return(tc.mockRec, nil)
//				mockRepo.On("GetItemById", mock.Anything, mock.Anything).Return("item1", nil)
//				mockRepo.On("GetUserById", mock.Anything, mock.Anything).Return(entity.User{Username: "user2"}, nil)
//			}
//
//			res, err := uc.GetInfo(context.Background(), tc.userId)
//
//			if (err != nil) != tc.wantErr {
//				t.Errorf("GetInfo() error = %v, wantErr %v", err, tc.wantErr)
//			}
//
//			if !tc.wantErr {
//				assert.Equal(t, tc.wantRes, res)
//			}
//
//			// Проверяем, что все ожидания Redis выполнены
//			if err := redisMock.ExpectationsWereMet(); err != nil {
//				t.Errorf("Redis expectations were not met: %v", err)
//			}
//
//			mockRepo.AssertExpectations(t)
//		})
//	}
//}
