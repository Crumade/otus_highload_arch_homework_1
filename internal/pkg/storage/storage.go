package storage

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	models "social_network/internal/model"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	HOST     = "host.docker.internal"
	PORT     = 5431
	USER     = "postgres"
	PASSWORD = "postgres"
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

func MigrateSchema(connDB *sqlx.DB) {
	driver, err := postgres.WithInstance(connDB.DB, &postgres.Config{})
	if err != nil {
		log.Fatal("Instance error: " + err.Error())
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:///app/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal("New DB Instance error: " + err.Error())
	}
	if err := m.Up(); err != nil {
		log.Fatal("Up migrations error: " + err.Error())
	}
}

func createIndexes(db *sqlx.DB) error {

	ext, err := db.Preparex("CREATE EXTENSION pg_trgm;")
	if err != nil {
		return err
	}
	_, err = ext.Exec()
	if err != nil {
		return err
	}

	index, err := db.Preparex("	CREATE INDEX users_names_idx ON users USING gist(second_name gist_trgm_ops, first_name gist_trgm_ops);")
	if err != nil {
		return err
	}
	_, err = index.Exec()
	if err != nil {
		return err
	}
	return nil
}

func MigrateUsers(connDB *sqlx.DB) error {
	file, err := os.Open("people.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	r := csv.NewReader(file)
	var placeholders []string
	var users []any
	index := 0
	start := time.Now()
	for {
		record, err := r.Read()
		if err == io.EOF {
			insertStatement := fmt.Sprintf("INSERT INTO users(first_name, second_name, birthdate, city) VALUES %s", strings.Join(placeholders, ","))
			//log.Printf("\n%+v", users...)
			_, err = connDB.Exec(insertStatement, users...)
			if err != nil {
				return err
			}
			users = nil
			placeholders = nil
			index = 0
			log.Printf("Для добавления в БД прошло времени %.2f c\n", time.Since(start).Seconds())
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fio := strings.Split(record[0], " ")

		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d,$%d,$%d)",
			index*4+1,
			index*4+2,
			index*4+3,
			index*4+4,
		))
		users = append(users, fio[1], fio[0], record[1], record[2])
		index++
		if len(users) == 65000 {

			insertStatement := fmt.Sprintf("INSERT INTO users(first_name, second_name, birthdate, city) VALUES %s", strings.Join(placeholders, ","))
			_, err = connDB.Exec(insertStatement, users...)
			if err != nil {
				return err
			}
			users = nil
			placeholders = nil
			index = 0
		}
	}

	err = createIndexes(connDB)
	if err != nil {
		return err
	}
	return nil
}

func GetUserByID(db *sqlx.DB, id string) (*models.User, error) {
	user := new(models.User)

	err := db.Get(user, "SELECT first_name, second_name, birthdate, gender, biography, city FROM public.users WHERE id = $1", id)
	if err == sql.ErrNoRows {
		err := errors.New("user not found")
		return nil, err
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func SearchUser(db *sqlx.DB, firstName string, lastName string) (*[]models.User, error) {
	users := new([]models.User)
	stm, err := db.Preparex(`SELECT id,
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

	err = stm.Select(users, firstName+"%", lastName+"%")
	if err == sql.ErrNoRows {
		err := errors.New("user not found")
		return nil, err
	} else if err != nil {
		return nil, err
	}
	return users, nil
}

func CreateUser(db *sqlx.DB, user *models.User) (*models.UserRegisterResponse, error) {
	result := new(models.UserRegisterResponse)
	rows, err := db.NamedQuery(`INSERT INTO users (first_name, second_name, birthdate, gender, biography, city) 
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

func GetAuthData(db *sqlx.DB, loginData *models.LoginRequest) (*models.AuthData, error) {
	authData := new(models.AuthData)
	err := db.Get(authData, "SELECT password_hash, salt FROM public.user_data WHERE user_id = $1", loginData.UserID)
	if err == sql.ErrNoRows {
		err := errors.New("user not found")
		return nil, err
	} else if err != nil {
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
