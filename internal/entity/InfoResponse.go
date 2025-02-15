package entity

type ResponseInfo struct {
	Coins       int         `json:"coins"`
	Inventory   Inventory   `json:"inventory"`
	CoinHistory CoinHistory `json:"coinHistory"`
}

type Inventory struct {
	Items []InventoryItem `json:"items"`
}

type InventoryItem struct {
	ItemId   int
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type CoinHistory struct {
	Received Received `json:"received"`
	Sent     Sent     `json:"sent"`
}

type Received struct {
	ReceivedItems []ReceivedItem `json:"items"`
}

type ReceivedItem struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type Sent struct {
	SentItems []SentItem `json:"items"`
}

type SentItem struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}
