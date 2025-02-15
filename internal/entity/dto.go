package entity

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SendCoinRequest struct {
	ToUserName string `json:"toUser"`
	Amount     int    `json:"amount"`
}

type BothDirection struct {
	ToUser   int `json:"toUser"`
	FromUser int `json:"fromUser"`
	Amount   int `json:"amount"`
}
