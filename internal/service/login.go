package service

import (
	"crypto/sha256"
	"errors"
	"fmt"
	models "social_network/internal/model"
	"social_network/internal/pkg/storage"

	"github.com/jmoiron/sqlx"
)

func Login(db *sqlx.DB, loginData *models.LoginRequest) (*models.LoginResponse, error) {

	authData, err := storage.GetAuthData(db, loginData)
	if err != nil {
		return nil, err
	}

	if authData.PasswordHash == HashedPassword(loginData.Password, authData.Salt) {

		token, err := storage.CreateAccessToken(db, loginData.UserID)
		if err != nil || token == "" {
			return nil, err
		}
		return &models.LoginResponse{Token: token}, nil
	}

	return nil, errors.New("something goes wrong")
}

func HashedPassword(password string, salt string) string {

	hash := sha256.Sum256([]byte(password + salt))
	return fmt.Sprintf("%x", hash)
}
