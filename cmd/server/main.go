package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nedpranson/the-christmas-breach/core/db"
	"github.com/nedpranson/the-christmas-breach/core/public"
	"github.com/nedpranson/the-christmas-breach/handlers/home"
)

func main() {
	os.Exit(run())
}

func run() int {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	conn, err := db.Open()
	if err != nil {
		slog.Error("db::Open failed", "err", err)
		return 1
	}
	defer conn.Close()

	if err := db.Migrate(conn); err != nil {
		slog.Error("db::Migrate failed", "err", err)
		return 1
	}

	mux := http.NewServeMux()
	mux.Handle("/", public.NewFileServer())

	mux.HandleFunc("GET /{$}", home.Home)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if err != http.ErrServerClosed {
			slog.Error("http::ListenAndServe failed", "err", err)
			return 1
		}
	case <-ctx.Done():
		srv.Close()
	}

	return 0
}
