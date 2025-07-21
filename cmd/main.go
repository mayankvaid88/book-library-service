package cmd

import (
	"book-store/internal/config"
	"book-store/internal/db"
	appHttp "book-store/internal/http"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func Main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		logrus.Fatalf("config load: %v", err)
	}

	db, err := db.NewDB(cfg)
	if err != nil {
		logrus.Fatalf("db init: %v", err)
	}
	defer db.Close()

	r := mux.NewRouter()
	appHttp.RegisterRoutes(r, db)
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		logrus.Fatalf("error while starting the server. error: %s", err.Error())
	}
}
