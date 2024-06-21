package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ConfingDB struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func Connect(cnf ConfingDB) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cnf.Host,
		cnf.Port,
		cnf.User,
		cnf.Password,
		cnf.Name,
	)
	db, err := sqlx.Connect("postgres", dsn)
	return db, err
}
