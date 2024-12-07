package storage

import (
	"database/sql"
	"errors"
	models "social_network/internal/model"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type postsRepo struct {
	conn  *sqlx.DB
	cache *redis.Client
}

func NewPostsRepo(db *sqlx.DB) *postsRepo {
	return &postsRepo{conn: db}
}

func (db *postsRepo) GetPostFeed(offset int, limit int) (*[]models.Post, error) {
	posts := new([]models.Post)
	stm, err := db.conn.Preparex(`SELECT id,
								user_id,
								content
							FROM public.posts 
							OFFSET $1
							LIMIT $2;
							`)
	if err != nil {
		return nil, err
	}

	err = stm.Select(posts, offset, limit)
	if err == sql.ErrNoRows {
		err := errors.New("posts not found")
		return nil, err
	} else if err != nil {
		return nil, err
	}
	return posts, nil
}

func (db *postsRepo) GetPostByID(id string) (*models.Post, error) {
	post := new(models.Post)

	err := db.conn.Get(post, "SELECT id, user_id, content FROM public.posts WHERE id = $1", id)
	if err == sql.ErrNoRows {
		err := errors.New("post not found")
		return nil, err
	} else if err != nil {
		return nil, err
	}

	return post, nil
}

func (db *postsRepo) DeletePost(id string) (bool, error) {
	stm, err := db.conn.Preparex("DELETE FROM public.posts WHERE id = $1:")
	if err != nil {
		return false, err
	}

	_, err = stm.Exec(id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (db *postsRepo) CreatePost(post *models.Post) (*models.NewPostResponse, error) {
	result := new(models.NewPostResponse)
	rows, err := db.conn.NamedQuery(`INSERT INTO public.posts (user_id, content, author) 
				VALUES(:user_id, :content, :author)
				RETURNING id;`, post)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	if rows.Next() {
		err := rows.StructScan(result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (db *postsRepo) GetPosts(limit int, offset int) (*[]models.Post, error) {

	return nil, nil
}
func (db *postsRepo) UpdatePost(post *models.Post) error {
	rows, err := db.conn.NamedQuery(`UPDATE public.posts 
									SET content = :content
									WHERE id = :id;`, post)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}
