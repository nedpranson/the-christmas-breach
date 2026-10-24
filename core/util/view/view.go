package view

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
)

func Render(w http.ResponseWriter, base templ.Component) {
	h := w.Header()

	h.Set("Content-Type", "text/html; charset=utf-8")
	base.Render(context.Background(), w)
}

func WithChildren(w http.ResponseWriter, base templ.Component, children templ.Component) {
	h := w.Header()

	h.Set("Content-Type", "text/html; charset=utf-8")
	base.Render(templ.WithChildren(context.Background(), children), w)
}

func RenderEx(w http.ResponseWriter, base templ.Component, ctx context.Context) {
	h := w.Header()

	h.Set("Content-Type", "text/html; charset=utf-8")
	base.Render(ctx, w)
}

func ErrorEx(w http.ResponseWriter, base templ.Component, ctx context.Context, code int) {
	h := w.Header()

	// Delete the Content-Length header, which might be for some other content.
	// Assuming the error string fits in the writer's buffer, we'll figure
	// out the correct Content-Length for it later.
	//
	// We don't delete Content-Encoding, because some middleware sets
	// Content-Encoding: gzip and wraps the ResponseWriter to compress on-the-fly.
	// See https://go.dev/issue/66343.
	h.Del("Content-Length")
	h.Set("Content-Type", "text/html; charset=utf-8")

	w.WriteHeader(code)
	base.Render(ctx, w)
}

func Error(w http.ResponseWriter, base templ.Component, code int) {
	h := w.Header()

	// Delete the Content-Length header, which might be for some other content.
	// Assuming the error string fits in the writer's buffer, we'll figure
	// out the correct Content-Length for it later.
	//
	// We don't delete Content-Encoding, because some middleware sets
	// Content-Encoding: gzip and wraps the ResponseWriter to compress on-the-fly.
	// See https://go.dev/issue/66343.
	h.Del("Content-Length")
	h.Set("Content-Type", "text/html; charset=utf-8")

	w.WriteHeader(code)
	base.Render(context.Background(), w)
}

func ErrorWithChildren(w http.ResponseWriter, base templ.Component, children templ.Component, code int) {
	h := w.Header()

	// Delete the Content-Length header, which might be for some other content.
	// Assuming the error string fits in the writer's buffer, we'll figure
	// out the correct Content-Length for it later.
	//
	// We don't delete Content-Encoding, because some middleware sets
	// Content-Encoding: gzip and wraps the ResponseWriter to compress on-the-fly.
	// See https://go.dev/issue/66343.
	h.Del("Content-Length")
	h.Set("Content-Type", "text/html; charset=utf-8")

	w.WriteHeader(code)
	base.Render(templ.WithChildren(context.Background(), children), w)
}
