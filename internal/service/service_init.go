package service

import (
	models "social_network/internal/model"

	"go.uber.org/zap"
)

type Service struct {
	userRepo  userRepo
	postsRepo postsRepo
	logger    *zap.Logger
}

func NewService(userRepo userRepo, postsRepo postsRepo, logger *zap.Logger) *Service {
	return &Service{
		userRepo:  userRepo,
		postsRepo: postsRepo,
		logger:    logger}
}

type userRepo interface {
	GetUserByID(string) (*models.User, error)
	SearchUser(string, string) (*[]models.User, error)
	CreateUser(*models.User) (*models.UserRegisterResponse, error)
	GetAuthData(*models.LoginRequest) (*models.AuthData, error)
	CreateAccessToken(string) (string, error)
	CreateAuthData(string, string, string) error
}

type postsRepo interface {
	GetPostFeed(int, int) (*[]models.Post, error)
	GetPostByID(string) (*models.Post, error)
	DeletePost(string) (bool, error)
}
