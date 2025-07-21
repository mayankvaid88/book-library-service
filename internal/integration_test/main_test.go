package integrationtest

import (
	"book-store/internal/config"
	"book-store/internal/db"
	appHttp "book-store/internal/http"
	"database/sql"
	"os"
	"testing"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var (
	router   *mux.Router
	sharedDB *sql.DB
)

func TestMain(m *testing.M) {
	cfg, err := config.LoadConfig("./config.json")
	if err != nil {
		logrus.Fatalf("config load failed: %v", err)
	}
	sharedDB, err = db.NewDB(cfg)
	if err != nil {
		logrus.Fatalf("db connect failed: %v", err)
	}
	router = mux.NewRouter()
	appHttp.RegisterRoutes(router, sharedDB)
	code := m.Run()
	sharedDB.Close()
	os.Exit(code)
}
