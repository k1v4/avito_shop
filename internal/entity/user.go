package entity

type User struct {
	Id       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Passhash []byte `json:"password" db:"password"`
}
