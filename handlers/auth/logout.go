package auth

import (
	"log/slog"
	"net/http"

	core_auth "github.com/nedpranson/the-christmas-breach/core/auth"
)

func LogOutPost(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(core_auth.SessionCookie)
	if err == nil {
		if err := core_auth.DeleteSession(r.Context(), c.Value); err != nil {
			slog.Error("core_auth::DeleteSession failed", "err", err)
		}
	}

	http.SetCookie(w, &http.Cookie{
        Name:     core_auth.SessionCookie,
        Value:    "",
        Path:     "/",
        MaxAge:   -1,
        HttpOnly: true,
		Secure:   false, // TODO: use secure
        SameSite: http.SameSiteLaxMode,
    })

	http.Redirect(w, r, "/log-in", http.StatusSeeOther)
}
