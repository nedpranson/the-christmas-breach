package auth

import (
	"log/slog"
	"net/http"

	"github.com/nedpranson/the-christmas-breach/core/user"
)

const SessionCookie = "session"

func LoadUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(SessionCookie)
		if err != nil || c.Value == "" {
			next.ServeHTTP(w, r)
			return
		}

		u, err := user.FromSession(r.Context(), c.Value)
		if err == user.ErrInvalidSession {
			next.ServeHTTP(w, r)
			return
		}

		if err != nil {
			slog.Error("user::FromSession failed", "err", err)
			next.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r.WithContext(user.WithUser(r.Context(), u)))
	})
}
