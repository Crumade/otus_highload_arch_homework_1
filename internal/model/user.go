package models

type User struct {
	ID         string `json:"id,omitempty" db:"id"`
	FirstName  string `json:"first_name" db:"first_name"`
	SecondName string `json:"second_name" db:"second_name"`
	BirthDate  string `json:"birthdate" db:"birthdate"`
	Gender     string `json:"gender" db:"gender,omitempty"`
	Biography  string `json:"biography" db:"biography,omitempty"`
	City       string `json:"city" db:"city"`
	Password   string `json:"password,omitempty"`
}

type UserRegisterResponse struct {
	UserID string `json:"user_id" db:"id"`
}
