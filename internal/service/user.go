package service

import (
	models "social_network/internal/model"

	"github.com/google/uuid"
)

func (srv *service) GetUser(id string) (*models.User, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	user, err := srv.userRepo.GetUserByID(id) //storage.GetUserByID(db, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (srv *service) Register(user *models.User) (*models.UserRegisterResponse, error) {
	newUser, err := srv.userRepo.CreateUser(user) //storage.CreateUser(db, user)
	if err != nil {
		return nil, err
	}
	salt := uuid.NewMD5(uuid.New(), []byte(user.Password))
	passwordHash := HashedPassword(user.Password, salt.String())
	err = srv.userRepo.CreateAuthData(newUser.UserID, passwordHash, salt.String()) //storage.CreateAuthData(db, newUser.UserID, passwordHash, salt.String())
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (srv *service) SearchUser(firstName string, lastName string) (*[]models.User, error) {

	users, err := srv.userRepo.SearchUser(firstName, lastName) //storage.SearchUser(db, firstName, lastName)
	if err != nil {
		return nil, err
	}
	return users, nil
}
