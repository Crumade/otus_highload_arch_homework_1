package models

type LoginRequest struct {
	UserID   string `json:"id"`
	Password string `json:"password"`
}

type AuthData struct {
	PasswordHash string `db:"password_hash"`
	Salt         string `db:"salt"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
