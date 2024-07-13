package service

import (
	models "social_network/internal/model"
	"social_network/internal/pkg/storage"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func GetUser(db *sqlx.DB, id string) (*models.User, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	user, err := storage.GetUserByID(db, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func Register(db *sqlx.DB, user *models.User) (*models.UserRegisterResponse, error) {
	newUser, err := storage.CreateUser(db, user)
	if err != nil {
		return nil, err
	}
	salt := uuid.NewMD5(uuid.New(), []byte(user.Password))
	passwordHash := HashedPassword(user.Password, salt.String())
	err = storage.CreateAuthData(db, newUser.UserID, passwordHash, salt.String())
	if err != nil {
		return nil, err
	}

	return newUser, nil
}
