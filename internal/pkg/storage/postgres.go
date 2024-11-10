package storage

import (
	"bufio"
	"encoding/csv"
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

func NewConnection() (*sqlx.DB, error) {
	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, USER, PASSWORD, DBNAME,
	)
	var err error
	conn, err := sqlx.Connect("pgx", connString)
	if err != nil {
		slog.Error("Connection error: " + err.Error())
		return nil, err
	}

	conn.SetConnMaxIdleTime(time.Second * 30)
	conn.SetConnMaxLifetime(time.Second * 30)
	conn.SetMaxIdleConns(100)
	conn.SetMaxOpenConns(100)

	if err = conn.Ping(); err != nil {
		slog.Error("Ping error: " + err.Error())
		return nil, err
	}

	return conn, nil
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

	users_index, err := db.Preparex(`
	CREATE INDEX users_names_idx ON users USING gist(second_name gist_trgm_ops, first_name gist_trgm_ops);
	`)
	if err != nil {
		return err
	}
	_, err = users_index.Exec()
	if err != nil {
		return err
	}

	friends_index, err := db.Preparex(`
	CREATE INDEX IF NOT EXISTS fki_friend_user_id
    ON public.friends USING btree
    (friend_user_id ASC NULLS LAST)
    TABLESPACE pg_default;
	`)
	if err != nil {
		return err
	}
	_, err = friends_index.Exec()
	if err != nil {
		return err
	}

	fki_users_index, err := db.Preparex(`
	CREATE INDEX IF NOT EXISTS fki_user_id
    ON public.friends USING btree
    (user_id ASC NULLS LAST)
    TABLESPACE pg_default;
	`)
	if err != nil {
		return err
	}
	_, err = fki_users_index.Exec()
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
