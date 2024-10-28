package storage

import (
	"database/sql"
	"errors"
	models "social_network/internal/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	conn *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *userRepo {
	return &userRepo{conn: db}
}

func (db *userRepo) GetUserByID(id string) (*models.User, error) {
	user := new(models.User)

	err := db.conn.Get(user, "SELECT first_name, second_name, birthdate, gender, biography, city FROM public.users WHERE id = $1", id)
	if err == sql.ErrNoRows {
		err := errors.New("user not found")
		return nil, err
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *userRepo) SearchUser(firstName string, lastName string) (*[]models.User, error) {
	users := new([]models.User)
	stm, err := db.conn.Preparex(`SELECT id,
						first_name, 
						second_name, 
						birthdate, 
						coalesce(gender, '') as gender, 
						coalesce(biography, '') as biography, 
						city 
					FROM public.users 
					WHERE second_name like  $1 
					AND first_name like $2
					ORDER BY id`)
	if err != nil {
		return nil, err
	}

	err = stm.Select(users, lastName+"%", firstName+"%")
	if err == sql.ErrNoRows {
		err := errors.New("user not found")
		return nil, err
	} else if err != nil {
		return nil, err
	}
	return users, nil
}

func (db *userRepo) CreateUser(user *models.User) (*models.UserRegisterResponse, error) {
	result := new(models.UserRegisterResponse)
	rows, err := db.conn.NamedQuery(`INSERT INTO users (first_name, second_name, birthdate, gender, biography, city) 
				VALUES(:first_name, :second_name, :birthdate, :gender, :biography, :city)
				RETURNING id;`, user)
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

func (db *userRepo) GetAuthData(loginData *models.LoginRequest) (*models.AuthData, error) {
	authData := new(models.AuthData)
	err := db.conn.Get(authData, "SELECT password_hash, salt FROM public.user_data WHERE user_id = $1", loginData.UserID)
	if err == sql.ErrNoRows {
		err := errors.New("user not found")
		return nil, err
	} else if err != nil {
		return nil, err
	}

	return authData, nil
}

func (db *userRepo) CreateAccessToken(userID string) (string, error) {
	token := uuid.New().String()

	_, err := db.conn.Exec("INSERT INTO tokens(access_token, user_id) VALUES($1, $2)", token, userID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (db *userRepo) CreateAuthData(userID string, passwordHash string, salt string) error {
	_, err := db.conn.Exec("INSERT INTO user_data(user_id, password_hash, salt) VALUES($1, $2, $3)", userID, passwordHash, salt)
	if err != nil {
		return err
	}

	return nil
}
