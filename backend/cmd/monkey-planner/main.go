package main

import (
	"context"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	mphttp "github.com/kjm99d/MonkeyPlanner/backend/internal/http"
	"github.com/kjm99d/MonkeyPlanner/backend/internal/service"
	"github.com/kjm99d/MonkeyPlanner/backend/internal/storage"
	_ "github.com/kjm99d/MonkeyPlanner/backend/internal/storage/postgres"
	_ "github.com/kjm99d/MonkeyPlanner/backend/internal/storage/sqlite"
	"github.com/kjm99d/MonkeyPlanner/backend/web"
)

// version is set at build time via -ldflags "-X main.version=v1.2.0"
var version = "dev"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "mcp" {
		// Subcommands under `mcp`: `mcp` (default: run server), `mcp install ...`.
		if len(os.Args) > 2 && os.Args[2] == "install" {
			runMCPInstall(os.Args[3:])
			return
		}
		runMCP()
		return
	}

	addr := getenv("MP_ADDR", ":8080")
	dsn := getenv("MP_DSN", defaultDSN())

	repo, err := storage.NewRepo(dsn)
	if err != nil {
		log.Fatalf("storage open: %v", err)
	}
	defer repo.Close()

	svc := service.New(repo, nil)

	// First-run bootstrap: seed a Welcome board + issue if the DB is empty.
	// Idempotent — on every subsequent start this is a cheap SELECT + early return.
	if err := svc.SeedWelcomeIfEmpty(context.Background()); err != nil {
		log.Printf("monkey-planner: welcome seed skipped: %v", err)
	}

	var static fs.FS
	if dist, err := web.Dist(); err == nil {
		static = dist
		log.Printf("monkey-planner: prod build — embedded frontend enabled")
	} else {
		log.Printf("monkey-planner: dev build — run Vite dev server at :5173 for UI (%v)", err)
	}

	router := mphttp.NewRouter(svc, static, version)

	srv := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second, // mitigates Slowloris (gosec G112)
		ReadTimeout:       30 * time.Second,
		// WriteTimeout is deliberately long: SSE streams hold the response
		// open until the client disconnects or the server shuts down.
		WriteTimeout: 0,
		IdleTimeout:  120 * time.Second,
	}

	// Listen for SIGINT/SIGTERM so Docker/k8s rolling updates drain in-flight
	// requests instead of killing the process mid-response.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("monkey-planner listening on %s (dsn=%s)", addr, dsn)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	select {
	case err := <-serverErr:
		if err != nil {
			log.Fatalf("server error: %v", err)
		}
	case <-ctx.Done():
		log.Printf("monkey-planner: shutdown signal received, draining...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("monkey-planner: graceful shutdown failed: %v", err)
	} else {
		log.Printf("monkey-planner: stopped cleanly")
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
