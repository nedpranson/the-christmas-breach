package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/nedpranson/the-christmas-breach/core/db"
	"github.com/nedpranson/the-christmas-breach/core/public"
	"github.com/nedpranson/the-christmas-breach/handlers/home"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	conn, err := db.Open()
	if err != nil {
		slog.Error("db::Open failed", "err", err)
		os.Exit(1)
	}
	defer conn.Close()

	if err := db.Migrate(conn); err != nil {
		slog.Error("db::Migrate failed", "err", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.Handle("/", public.NewFileServer())

	mux.HandleFunc("GET /{$}", home.Home)

	if err := http.ListenAndServe(":3000", mux); err != nil {
		slog.Error("http::ListenAndServe failed", "err", err)
		os.Exit(1)
	}
}
