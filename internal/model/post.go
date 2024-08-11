package models

type Post struct {
	ID     string `json:"id"`
	Text   string `json:"text"`
	Author string `json:"author_user_id"`
}
