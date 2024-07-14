package db

import (
	"fmt"
	"log"

	"github.com/theborzet/time-tracker/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func Init(cfg *config.Config) *sqlx.DB {
	url := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Pass, cfg.DB.Name)

	db, err := sqlx.Open("postgres", url)

	if err != nil {
		log.Fatalln(err)
	}

	return db
}

func Close(db *sqlx.DB) error {
	if err := db.Close(); err != nil {
		return err
	}
	return nil
}
