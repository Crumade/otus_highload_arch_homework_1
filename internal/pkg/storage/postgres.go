package storage

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
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

type PostgresDB struct {
	Conn *sqlx.DB
}

func (pg *PostgresDB) NewConnection() error {
	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, USER, PASSWORD, DBNAME,
	)
	var err error
	pg.Conn, err = sqlx.Connect("pgx", connString)
	if err != nil {
		slog.Error("Connection error: " + err.Error())
		return err
	}

	pg.Conn.SetConnMaxIdleTime(time.Second * 30)
	pg.Conn.SetConnMaxLifetime(time.Second * 30)
	pg.Conn.SetMaxIdleConns(10)
	pg.Conn.SetMaxOpenConns(10)

	if err = pg.Conn.Ping(); err != nil {
		slog.Error("Ping error: " + err.Error())
		return err
	}

	return nil
}

func MigrateSchema(db *sqlx.DB) {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
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

func MigrateUsers(db *sqlx.DB) error {
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
			_, err = db.Exec(insertStatement, users...)
			if err != nil {
				return err
			}
			users = nil
			placeholders = nil
			index = 0
			slog.Info(fmt.Sprintf("Для добавления юзеров в БД прошло времени %.2f c", time.Since(start).Seconds()))
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
			_, err = db.Exec(insertStatement, users...)
			if err != nil {
				return err
			}
			users = nil
			placeholders = nil
			index = 0
		}
	}

	err = createIndexes(db)
	if err != nil {
		return err
	}
	return nil
}

func MigratePosts(db *sqlx.DB) error {
	file, err := os.Open("posts.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var placeholders []string
	var posts []any
	index := 0
	start := time.Now()
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		user := new(models.User)
		err = db.Get(user, `SELECT id
								FROM users 
								OFFSET floor(random()*8391) 
								LIMIT 1`)
		if err != nil {
			return err
		}
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)",
			index*2+1,
			index*2+2,
		))
		posts = append(posts, user.ID, line)
		index++

	}

	tempInsert := fmt.Sprintf(`
					INSERT INTO posts(user_id, content) VALUES %s;`,
		strings.Join(placeholders, ","))
	_, err = db.Exec(tempInsert, posts...)
	if err != nil {
		return err
	}

	posts = nil
	placeholders = nil
	index = 0
	slog.Info(fmt.Sprintf("Для добавления постов в БД прошло времени %.2f c", time.Since(start).Seconds()))
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

	err = stm.Select(users, lastName+"%", firstName+"%")
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

func GetPostFeed(db *sqlx.DB, offset int, limit int) (*[]models.Post, error) {
	posts := new([]models.Post)
	stm, err := db.Preparex(`SELECT id,
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

func GetPostByID(db *sqlx.DB, id string) (*models.Post, error) {
	post := new(models.Post)

	err := db.Get(post, "SELECT id, user_id, content FROM public.posts WHERE id = $1", id)
	if err == sql.ErrNoRows {
		err := errors.New("post not found")
		return nil, err
	} else if err != nil {
		return nil, err
	}

	return post, nil
}

func DeletePost(db *sqlx.DB, id string) (bool, error) {
	stm, err := db.Preparex("DELETE FROM public.posts WHERE id = $1:")
	if err != nil {
		return false, err
	}

	_, err = stm.Exec(id)
	if err != nil {
		return false, err
	}

	return true, nil
}
