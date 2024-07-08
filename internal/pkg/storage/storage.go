package storage

import (
	"fmt"
	"log"
	models "social_network/internal/model"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	HOST     = "localhost"
	PORT     = 5432
	USER     = "postgres"
	PASSWORD = "731596"
	DBNAME   = "social_network"
)

func NewConnection() (*sqlx.DB, error) {
	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, USER, PASSWORD, DBNAME,
	)

	db, err := sqlx.Connect("pgx", connString)
	if err != nil {
		log.Println("Connection error: " + err.Error())
		return nil, err
	}

	db.SetConnMaxIdleTime(time.Second * 30)
	db.SetConnMaxLifetime(time.Second * 30)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	if err = db.Ping(); err != nil {
		log.Println("Ping error: " + err.Error())
		return nil, err
	}

	return db, nil
}

func GetUserByID(db *sqlx.DB, id string) (*models.User, error) {
	user := new(models.User)

	err := db.QueryRowx("SELECT first_name, second_name, birthdate, gender, biography, city FROM public.users WHERE id = $1", id).StructScan(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func CreateUser(db *sqlx.DB, user *models.User) (*models.UserRegisterResponse, error) {
	result := new(models.UserRegisterResponse)
	rows, err := db.NamedQuery(`INSERT INTO users (first_name, second_name, birthdate, gender, biography, city) 
				VALUES(:first_name, :second_name, :birthdate, :gender, :biography, :city)
				RETURNING id;`, user)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err := rows.StructScan(result)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	return result, nil
}

func GetAuthData(db *sqlx.DB, loginData *models.LoginRequest) (*models.AuthData, error) {
	authData := new(models.AuthData)
	err := db.QueryRowx("SELECT password_hash, salt FROM public.user_data WHERE user_id = $1", loginData.UserID).StructScan(authData)
	if err != nil {
		return nil, err
	}

	return authData, nil
}

func CreateAccessToken(db *sqlx.DB, userID string) (string, error) {
	token := uuid.New().String()

	_, err := db.Exec("INSERT INTO tokens(access_token, user_id) VALUES($1, $2)", token, userID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func CreateAuthData(db *sqlx.DB, userID string, passwordHash string, salt string) error {
	_, err := db.Exec("INSERT INTO user_data(user_id, password_hash, salt) VALUES($1, $2, $3)", userID, passwordHash, salt)
	if err != nil {
		return err
	}

	return nil
}
