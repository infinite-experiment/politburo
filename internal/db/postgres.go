package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func InitPostgres() error {

	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	user := os.Getenv("PG_USER")
	dbname := os.Getenv("PG_DB")
	password := os.Getenv("PG_PASSWORD")
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

	var err error
	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		return err
	}

	return nil
}
