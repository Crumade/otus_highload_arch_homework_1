package service

import (
	"crypto/sha256"
	"errors"
	"fmt"
	models "social_network/internal/model"
)

func (srv *service) Login(loginData *models.LoginRequest) (*models.LoginResponse, error) {

	authData, err := srv.userRepo.GetAuthData(loginData) //storage.GetAuthData(db, loginData)
	if err != nil {
		return nil, err
	}

	if authData.PasswordHash == HashedPassword(loginData.Password, authData.Salt) {

		token, err := srv.userRepo.CreateAccessToken(loginData.UserID) //storage.CreateAccessToken(db, loginData.UserID)
		if err != nil || token == "" {
			return nil, err
		}
		return &models.LoginResponse{Token: token}, nil
	}

	return nil, errors.New("user id or password invalid")
}

func HashedPassword(password string, salt string) string {

	hash := sha256.Sum256([]byte(password + salt))
	return fmt.Sprintf("%x", hash)
}
