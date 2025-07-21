package db

import (
	"book-store/internal/config"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func NewDB(cfg config.Config) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		get(cfg.GetHost()), get(cfg.GetPort()),
		get(cfg.GetUser()), get(cfg.GetPassword()),
		get(cfg.GetName()),
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		logrus.Fatalf("Unable to open DB: %v", err)
	}
	if err = db.Ping(); err != nil {
		logrus.Fatalf("Unable to connect: %v", err)
	}
	return db, nil
}

var get = func(name string) string {
	v := os.Getenv(name)
	if v == "" {
		logrus.Fatalf("db config %s not present in env var", name)
	}
	return v
}
