package models

type Post struct {
	ID     string `json:"id" db:"id"`
	Text   string `json:"text" db:"content"`
	Author string `json:"author_user_id" db:"user_id"`
}
