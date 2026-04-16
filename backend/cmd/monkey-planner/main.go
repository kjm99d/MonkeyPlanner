package main

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	mphttp "github.com/kjm99d/monkey-planner/backend/internal/http"
	"github.com/kjm99d/monkey-planner/backend/internal/service"
	"github.com/kjm99d/monkey-planner/backend/internal/storage"
	_ "github.com/kjm99d/monkey-planner/backend/internal/storage/postgres"
	_ "github.com/kjm99d/monkey-planner/backend/internal/storage/sqlite"
	"github.com/kjm99d/monkey-planner/backend/web"
)

func main() {
	addr := getenv("MP_ADDR", ":8080")
	dsn := getenv("MP_DSN", defaultDSN())

	repo, err := storage.NewRepo(dsn)
	if err != nil {
		log.Fatalf("storage open: %v", err)
	}
	defer repo.Close()

	svc := service.New(repo, nil)

	var static fs.FS
	if dist, err := web.Dist(); err == nil {
		static = dist
		log.Printf("monkey-planner: prod build — embedded frontend enabled")
	} else {
		log.Printf("monkey-planner: dev build — run Vite dev server at :5173 for UI (%v)", err)
	}

	router := mphttp.NewRouter(svc, static)

	log.Printf("monkey-planner listening on %s (dsn=%s)", addr, dsn)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func defaultDSN() string {
	_ = os.MkdirAll("./data", 0o755)
	return "sqlite://" + filepath.Join("./data", "monkey.db")
}

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
