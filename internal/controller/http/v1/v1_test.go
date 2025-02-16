package v1

import (
	"encoding/json"
	"errors"
	"github.com/k1v4/avito_shop/internal/entity"
	"github.com/k1v4/avito_shop/internal/usecase"
	"github.com/k1v4/avito_shop/internal/usecase/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInfo(t *testing.T) {
	cases := []struct {
		name       string
		token      string
		mockInfo   entity.ResponseInfo
		mockErr    error
		statusCode int
		respBody   string
		wantErr    bool
		isMock     bool
	}{
		{
			name:  "success",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDIyOTQxNDIsImlkIjoxMjIxMiwidXNlcm5hbWUiOiJUcmV2b3I2OCJ9.ChQ2sD0avEBglTA7CHBDKRzjP4FIRNeJy68hm4todkQ",
			mockInfo: entity.ResponseInfo{
				Coins: 100,
				Inventory: entity.Inventory{
					Items: []entity.InventoryItem{
						{ItemId: 1, Type: "type1", Quantity: 5},
						{ItemId: 2, Type: "type2", Quantity: 10},
					},
				},
				CoinHistory: entity.CoinHistory{
					Received: entity.Received{
						ReceivedItems: []entity.ReceivedItem{
							{FromUser: "user2", Amount: 50},
							{FromUser: "user3", Amount: 30},
						},
					},
					Sent: entity.Sent{
						SentItems: []entity.SentItem{
							{ToUser: "user4", Amount: 20},
							{ToUser: "user5", Amount: 10},
						},
					},
				},
			},
			mockErr:    nil,
			statusCode: http.StatusOK,
			respBody: `{
				"coins": 100,
				"inventory": {
					"items": [
						{"type": "type1", "quantity": 5},
						{"type": "type2", "quantity": 10}
					]
				},
				"coinHistory": {
					"received": {
						"items": [
							{"fromUser": "user2", "amount": 50},
							{"fromUser": "user3", "amount": 30}
						]
					},
					"sent": {
						"items": [
							{"toUser": "user4", "amount": 20},
							{"toUser": "user5", "amount": 10}
						]
					}
				}
			}`,
			wantErr: false,
			isMock:  true,
		},
		{
			name:       "no_token",
			token:      "",
			mockInfo:   entity.ResponseInfo{},
			mockErr:    nil,
			statusCode: http.StatusBadRequest,
			respBody:   `{"error":"bad request"}`,
			wantErr:    true,
			isMock:     false,
		},
		{
			name:       "internal_error",
			token:      "valid_token",
			mockInfo:   entity.ResponseInfo{},
			mockErr:    errors.New("internal error"),
			statusCode: http.StatusUnauthorized,
			respBody:   `{"error":"bad request"}`,
			wantErr:    true,
			isMock:     false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/info", nil)
			if tc.token != "" {
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockService := new(mocks.IShopService)

			if tc.isMock {
				mockService.
					On("GetInfo", c.Request().Context(), mock.Anything).
					Return(tc.mockInfo, tc.mockErr)
			}

			handler := &conatainerRoutes{t: mockService}
			err := handler.Info(c)

			if (err != nil) != tc.wantErr {
				t.Errorf("Info() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			assert.Equal(t, tc.statusCode, rec.Code)
			assert.JSONEq(t, tc.respBody, rec.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}

func TestSendCoins(t *testing.T) {
	cases := []struct {
		name       string
		reqBody    string
		token      string
		mockErr    error
		statusCode int
		respBody   string
		wantErr    bool
		isMock     bool
	}{
		{
			name:       "success",
			reqBody:    `{"toUserName":"user2","amount":100}`,
			token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDIyOTQxNDIsImlkIjoxMjIxMiwidXNlcm5hbWUiOiJUcmV2b3I2OCJ9.ChQ2sD0avEBglTA7CHBDKRzjP4FIRNeJy68hm4todkQ",
			mockErr:    nil,
			statusCode: http.StatusOK,
			respBody:   `{}`,
			wantErr:    false,
			isMock:     true,
		},
		{
			name:       "negative_amount",
			reqBody:    `{"toUserName":"user2","amount":-100}`,
			token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDIyOTQxNDIsImlkIjoxMjIxMiwidXNlcm5hbWUiOiJUcmV2b3I2OCJ9.ChQ2sD0avEBglTA7CHBDKRzjP4FIRNeJy68hm4todkQ",
			mockErr:    nil,
			statusCode: http.StatusBadRequest,
			respBody:   `{"error":"bad request"}`,
			wantErr:    true,
			isMock:     false,
		},
		{
			name:       "bad_body",
			reqBody:    `{"tttt":"user2","amount":-100}`,
			token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDIyOTQxNDIsImlkIjoxMjIxMiwidXNlcm5hbWUiOiJUcmV2b3I2OCJ9.ChQ2sD0avEBglTA7CHBDKRzjP4FIRNeJy68hm4todkQ",
			mockErr:    nil,
			statusCode: http.StatusBadRequest,
			respBody:   `{"error":"bad request"}`,
			wantErr:    true,
			isMock:     false,
		},
		{
			name:       "no_token",
			reqBody:    `{"toUserName":"user2","amount":100}`,
			token:      "",
			mockErr:    nil,
			statusCode: http.StatusUnauthorized,
			respBody:   `{"error":"bad request"}`,
			wantErr:    true,
			isMock:     false,
		},
		{
			name:       "internal_error",
			reqBody:    `{"toUserName":"user2","amount":100}`,
			token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDIyOTQxNDIsImlkIjoxMjIxMiwidXNlcm5hbWUiOiJUcmV2b3I2OCJ9.ChQ2sD0avEBglTA7CHBDKRzjP4FIRNeJy68hm4todkQ",
			mockErr:    errors.New("internal error"),
			statusCode: http.StatusInternalServerError,
			respBody:   `{"error":"internal error"}`,
			wantErr:    true,
			isMock:     true,
		},
		{
			name:       "internal_error",
			reqBody:    `{"toUserName":"user2","amount":100}`,
			token:      "a.a.a",
			mockErr:    errors.New("internal error"),
			statusCode: http.StatusUnauthorized,
			respBody:   `{"error":"bad request"}`,
			wantErr:    true,
			isMock:     false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/sendCoin", strings.NewReader(tc.reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			if tc.token != "" {
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockService := new(mocks.IShopService)

			if tc.isMock {
				mockService.
					On("SendCoins", c.Request().Context(), mock.Anything, mock.Anything, mock.Anything).
					Return(tc.mockErr)
			}

			handler := &conatainerRoutes{t: mockService}
			err := handler.SendCoins(c)

			if (err != nil) != tc.wantErr {
				t.Errorf("SendCoins() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			assert.Equal(t, tc.statusCode, rec.Code)
			assert.JSONEq(t, tc.respBody, rec.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}

func TestBuy(t *testing.T) {
	cases := []struct {
		name       string
		item       string
		token      string
		mockErr    error
		statusCode int
		respBody   string
		wantErr    bool
		isMock     bool
	}{
		{
			name:       "success",
			item:       "item1",
			token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDIyOTQxNDIsImlkIjoxMjIxMiwidXNlcm5hbWUiOiJUcmV2b3I2OCJ9.ChQ2sD0avEBglTA7CHBDKRzjP4FIRNeJy68hm4todkQ",
			mockErr:    nil,
			statusCode: http.StatusOK,
			respBody:   `{}`,
			wantErr:    false,
			isMock:     true,
		},
		{
			name:       "no_token",
			item:       "item1",
			token:      "",
			mockErr:    nil,
			statusCode: http.StatusUnauthorized,
			respBody:   `{"error":"bad request"}`,
			wantErr:    true,
			isMock:     false,
		},
		{
			name:       "no_coins",
			item:       "item1",
			token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDIyOTQxNDIsImlkIjoxMjIxMiwidXNlcm5hbWUiOiJUcmV2b3I2OCJ9.ChQ2sD0avEBglTA7CHBDKRzjP4FIRNeJy68hm4todkQ",
			mockErr:    usecase.ErrNoCoins,
			statusCode: http.StatusBadRequest,
			respBody:   `{"error":"not enough coins"}`,
			wantErr:    true,
			isMock:     true,
		},
		{
			name:       "fail_validate_token",
			item:       "item1",
			token:      "valid_token",
			mockErr:    errors.New("internal error"),
			statusCode: http.StatusUnauthorized,
			respBody:   `{"error":"bad request"}`,
			wantErr:    true,
			isMock:     false,
		},
		{
			name:       "no_item",
			item:       "",
			mockErr:    nil,
			statusCode: http.StatusBadRequest,
			respBody:   `{"error":"bad request"}`,
			wantErr:    true,
			isMock:     false,
		},
		{
			name:       "internal_error",
			item:       "wallet",
			token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDIyOTQxNDIsImlkIjoxMjIxMiwidXNlcm5hbWUiOiJUcmV2b3I2OCJ9.ChQ2sD0avEBglTA7CHBDKRzjP4FIRNeJy68hm4todkQ",
			mockErr:    errors.New("internal error"),
			statusCode: http.StatusInternalServerError,
			respBody:   `{"error":"internal error"}`,
			wantErr:    true,
			isMock:     true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/buy/"+tc.item, nil)
			if tc.token != "" {
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("item")
			c.SetParamValues(tc.item)

			mockService := new(mocks.IShopService)

			if tc.isMock {
				mockService.
					On("BuyItem", c.Request().Context(), mock.Anything, tc.item).
					Return(tc.mockErr)
			}

			handler := &conatainerRoutes{t: mockService}
			err := handler.Buy(c)

			if (err != nil) != tc.wantErr {
				t.Errorf("Buy() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			assert.Equal(t, tc.statusCode, rec.Code)
			assert.JSONEq(t, tc.respBody, rec.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}

func TestAuth(t *testing.T) {
	cases := []struct {
		name       string
		reqBody    entity.AuthRequest
		mockToken  string
		mockErr    error
		statusCode int
		respBody   entity.AuthResponse
		wantErr    bool
		isMock     bool
	}{
		{
			name: "success",
			reqBody: entity.AuthRequest{
				Username: "user1",
				Password: "pass1",
			},
			mockToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDIyOTQxNDIsImlkIjoxMjIxMiwidXNlcm5hbWUiOiJUcmV2b3I2OCJ9.ChQ2sD0avEBglTA7CHBDKRzjP4FIRNeJy68hm4todkQ",
			mockErr:    nil,
			statusCode: http.StatusOK,
			respBody: entity.AuthResponse{
				Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDIyOTQxNDIsImlkIjoxMjIxMiwidXNlcm5hbWUiOiJUcmV2b3I2OCJ9.ChQ2sD0avEBglTA7CHBDKRzjP4FIRNeJy68hm4todkQ",
			},
			wantErr: false,
			isMock:  true,
		},
		{
			name: "empty_username",
			reqBody: entity.AuthRequest{
				Username: "",
				Password: "pass1",
			},
			mockToken:  "",
			mockErr:    errors.New("invalid credentials"),
			statusCode: http.StatusBadRequest,
			respBody:   entity.AuthResponse{},
			wantErr:    true,
			isMock:     false,
		},
		{
			name: "empty_password",
			reqBody: entity.AuthRequest{
				Username: "user1",
				Password: "",
			},
			mockToken:  "",
			mockErr:    errors.New("invalid credentials"),
			statusCode: http.StatusBadRequest,
			respBody:   entity.AuthResponse{},
			wantErr:    true,
			isMock:     false,
		},
		{
			name: "wrong_password",
			reqBody: entity.AuthRequest{
				Username: "user1",
				Password: "wrongpass",
			},
			mockToken:  "",
			mockErr:    usecase.ErrWrongPassword,
			statusCode: http.StatusUnauthorized,
			respBody:   entity.AuthResponse{},
			wantErr:    true,
			isMock:     true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()

			// Сериализуем reqBody в JSON
			reqBodyJSON, err := json.Marshal(tc.reqBody)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/auth", strings.NewReader(string(reqBodyJSON)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockService := new(mocks.IShopService)

			if tc.isMock {
				mockService.
					On("Login", c.Request().Context(), tc.reqBody.Username, tc.reqBody.Password).
					Return(tc.mockToken, tc.mockErr)
			}

			handler := &conatainerRoutes{t: mockService}
			err = handler.Auth(c)

			if (err != nil) != tc.wantErr {
				t.Errorf("Auth() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			assert.Equal(t, tc.statusCode, rec.Code)

			if !tc.wantErr {
				var respBody entity.AuthResponse
				err = json.Unmarshal(rec.Body.Bytes(), &respBody)
				assert.NoError(t, err)
				assert.Equal(t, tc.respBody, respBody)
			} else {
				assert.Contains(t, rec.Body.String(), "bad request")
			}

			mockService.AssertExpectations(t)
		})
	}
}
