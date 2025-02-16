package entity

import (
	"testing"
)

func TestMarshalBinary(t *testing.T) {
	response := ResponseInfo{
		Coins: 100,
		Inventory: Inventory{
			Items: []InventoryItem{
				{ItemId: 1, Type: "gold", Quantity: 10},
				{ItemId: 2, Type: "silver", Quantity: 20},
			},
		},
		CoinHistory: CoinHistory{
			Received: Received{
				ReceivedItems: []ReceivedItem{
					{FromUser: "user1", Amount: 50},
				},
			},
			Sent: Sent{
				SentItems: []SentItem{
					{ToUser: "user2", Amount: 30},
				},
			},
		},
	}

	data, err := response.MarshalBinary()
	if err != nil {
		t.Errorf("Failed to marshal ResponseInfo: %v", err)
	}

	// Проверяем, что данные сериализуются в правильный JSON
	expectedJSON := `{"coins":100,"inventory":{"items":[{"type":"gold","quantity":10},{"type":"silver","quantity":20}]},"coinHistory":{"received":{"items":[{"fromUser":"user1","amount":50}]},"sent":{"items":[{"toUser":"user2","amount":30}]}}}`
	if string(data) != expectedJSON {
		t.Errorf("Expected %s but got %s", expectedJSON, string(data))
	}
}

func TestUnmarshalBinary(t *testing.T) {
	jsonInput := `{"coins":100,"inventory":{"items":[{"type":"gold","quantity":10},{"type":"silver","quantity":20}]},"coinHistory":{"received":{"items":[{"fromUser":"user1","amount":50}]},"sent":{"items":[{"toUser":"user2","amount":30}]}}}`

	var response ResponseInfo
	err := response.UnmarshalBinary([]byte(jsonInput))
	if err != nil {
		t.Errorf("Failed to unmarshal ResponseInfo: %v", err)
	}

	// Проверяем, что данные десериализуются правильно
	if response.Coins != 100 {
		t.Errorf("Expected coins to be 100, got %d", response.Coins)
	}
	if len(response.Inventory.Items) != 2 {
		t.Errorf("Expected 2 inventory items, got %d", len(response.Inventory.Items))
	}
	if response.CoinHistory.Received.ReceivedItems[0].FromUser != "user1" {
		t.Errorf("Expected first received item from user to be 'user1', got %s", response.CoinHistory.Received.ReceivedItems[0].FromUser)
	}
	if response.CoinHistory.Sent.SentItems[0].Amount != 30 {
		t.Errorf("Expected first sent item amount to be 30, got %d", response.CoinHistory.Sent.SentItems[0].Amount)
	}
}
