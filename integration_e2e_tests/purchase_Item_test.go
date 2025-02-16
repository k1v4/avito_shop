package integration_tests

import (
	"bytes"
	"encoding/json"
	"github.com/k1v4/avito_shop/internal/entity"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestPurchaseItem_E2E(t *testing.T) {
	// Шаг 1: Аутентификация пользователя
	authReq := map[string]string{
		"username": "testuser",
		"password": "password123",
	}
	authBody, _ := json.Marshal(authReq)
	resp, err := http.Post("http://localhost:8080/api/auth", "application/json", bytes.NewBuffer(authBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var authResp map[string]string
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	assert.NoError(t, err)
	token := authResp["token"]

	// Шаг 2: Покупка предмета
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://localhost:8080/api/buy/hoody", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Шаг 3: Проверка баланса и инвентаря
	req, err = http.NewRequest("GET", "http://localhost:8080/api/info", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var infoResp entity.ResponseInfo
	err = json.NewDecoder(resp.Body).Decode(&infoResp)
	assert.NoError(t, err)

	// проверка баланса
	assert.Equal(t, 700, infoResp.Coins)

	// проверка инвентаря
	inventory := infoResp.Inventory
	found := false
	for _, item := range inventory.Items {
		if item.Type == "hoody" && item.Quantity == 1 {
			found = true
			break
		}
	}
	assert.True(t, found, "Предмет hoody должен быть в инвентаре")
}
