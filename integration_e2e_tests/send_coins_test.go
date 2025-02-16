package integration_tests

import (
	"bytes"
	"encoding/json"
	"github.com/k1v4/avito_shop/internal/entity"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_SendCoins(t *testing.T) {
	// Аутентификация пользователя 1
	authReqUser1 := map[string]string{
		"username": "user_1",
		"password": "123",
	}

	authBody, _ := json.Marshal(authReqUser1)
	respAuth1, err := http.Post("http://localhost:8080/api/auth", "application/json", bytes.NewBuffer(authBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respAuth1.StatusCode)

	var authRespUser1 map[string]string
	err = json.NewDecoder(respAuth1.Body).Decode(&authRespUser1)
	assert.NoError(t, err)
	tokenUser1 := authRespUser1["token"]

	// Аутентификация пользователя 2
	authReqUser2 := map[string]string{
		"username": "user_2",
		"password": "321",
	}
	authBody, _ = json.Marshal(authReqUser2)
	respAuth2, err := http.Post("http://localhost:8080/api/auth", "application/json", bytes.NewBuffer(authBody))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respAuth2.StatusCode)

	var authResp map[string]string
	err = json.NewDecoder(respAuth2.Body).Decode(&authResp)

	assert.NoError(t, err)
	tokenUser2 := authResp["token"]

	// user1 отправляет 100 монет user2
	user1ToUser2 := entity.SendCoinRequest{
		ToUserName: "user_2",
		Amount:     100,
	}
	client := &http.Client{}

	sendCoinsBody, _ := json.Marshal(user1ToUser2)
	req, err := http.NewRequest("POST", "http://localhost:8080/api/sendCoin", bytes.NewBuffer(sendCoinsBody))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenUser1)
	respSentCoins, err := client.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respSentCoins.StatusCode)

	// Проверяем баланс обоих. Должно стать user1:900 user2: 1100
	// User_1
	req, err = http.NewRequest("GET", "http://localhost:8080/api/info", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+tokenUser1)
	respGetInfo, err := client.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respGetInfo.StatusCode)

	var infoResp1 entity.ResponseInfo
	err = json.NewDecoder(respGetInfo.Body).Decode(&infoResp1)
	assert.NoError(t, err)

	assert.Equal(t, 900, infoResp1.Coins)

	//User_2
	req, err = http.NewRequest("GET", "http://localhost:8080/api/info", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+tokenUser2)
	respGetInfo2, err := client.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respGetInfo2.StatusCode)

	var infoResp2 entity.ResponseInfo
	err = json.NewDecoder(respGetInfo2.Body).Decode(&infoResp2)
	assert.NoError(t, err)

	assert.Equal(t, 1100, infoResp2.Coins)

	// Проверяем таблицу операций
	sendInfoUser1 := infoResp1.CoinHistory.Sent.SentItems[0]
	assert.Equal(t, sendInfoUser1, entity.SentItem{
		ToUser: "user_2",
		Amount: 100,
	})

	sendInfoUser2 := infoResp2.CoinHistory.Received.ReceivedItems[0]
	assert.Equal(t, sendInfoUser2, entity.ReceivedItem{
		FromUser: "user_1",
		Amount:   100,
	})

}
