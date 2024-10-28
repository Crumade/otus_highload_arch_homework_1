package service

import (
	models "social_network/internal/model"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func (srv *service) GetPostFeed(cache *redis.Client, offset int, limit int) (*[]models.Post, error) {

	posts, err := srv.postsRepo.GetPostFeed(offset, limit) //storage.GetPostFeed(db, offset, limit)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (srv *service) DeletePost(id string) (bool, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return false, err
	}

	result, err := srv.postsRepo.DeletePost(id) //storage.DeletePost(db, id)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (srv *service) GetPost(id string) (*models.Post, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	post, err := srv.postsRepo.GetPostByID(id) //storage.GetPostByID(db, id)
	if err != nil {
		return nil, err
	}
	return post, nil
}
