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
	ToUser string `json:"to_user"`
	Amount int    `json:"amount"`
}
