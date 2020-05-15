package pgsql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"telebot/telebot/CA/pkg/infrastructure/config"
)

func New(p config.DBconfig) (*sql.DB, error) {
	dsn := "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
	dsn = fmt.Sprintf(dsn, p.Host, p.Port, p.User, p.Password, p.DBname)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("Что-то не так в New(p config.DBconfig): %v", err)
		return nil, err
	}

	return db, nil
}
