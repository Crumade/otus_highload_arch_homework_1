package service

import (
	models "social_network/internal/model"
	"social_network/internal/pkg/storage"

	"github.com/jmoiron/sqlx"
)

func GetPostFeed(db *sqlx.DB, offset int, limit int) (*[]models.Post, error) {

	posts, err := storage.GetPostFeed(db, offset, limit)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
