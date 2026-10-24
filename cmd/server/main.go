package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/nedpranson/the-christmas-breach/handlers/home"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", home.Home)

	if err := http.ListenAndServe(":3000", mux); err != nil {
		slog.Error("http::ListenAndServe failed", "err", err)
		os.Exit(1)
	}
}
