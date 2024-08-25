package service

import (
	models "social_network/internal/model"
	"social_network/internal/pkg/storage"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func GetPostFeed(db *sqlx.DB, cache *redis.Client, offset int, limit int) (*[]models.Post, error) {

	posts, err := storage.GetPostFeed(db, offset, limit)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func DeletePost(db *sqlx.DB, id string) (bool, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return false, err
	}

	result, err := storage.DeletePost(db, id)
	if err != nil {
		return false, err
	}
	return result, nil
}

func GetPost(db *sqlx.DB, id string) (*models.Post, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	post, err := storage.GetPostByID(db, id)
	if err != nil {
		return nil, err
	}
	return post, nil
}
