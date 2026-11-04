package auth

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/httprate"
	core_auth "github.com/nedpranson/the-christmas-breach/core/auth"
	"github.com/nedpranson/the-christmas-breach/core/user"
	"github.com/nedpranson/the-christmas-breach/core/util"
	"github.com/nedpranson/the-christmas-breach/core/util/view"
	auth_templ "github.com/nedpranson/the-christmas-breach/templates/auth"
	base_templ "github.com/nedpranson/the-christmas-breach/templates/base"
)

var loginRateLimiter = httprate.NewRateLimiter(10, 30*time.Minute)

func LogIn(w http.ResponseWriter, r *http.Request) {
	view.Render(w, auth_templ.LogIn(), r.Context())
}

func LogInPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	persistent := r.FormValue("persistent") != ""

	if username == "" || password == "" {
		view.Error(w, base_templ.Error(base_templ.ErrorProps{
			Message: "Username and password are required.",
			Status:  http.StatusBadRequest,
		}), r.Context(), http.StatusBadRequest)
		return
	}

	if loginRateLimiter.OnLimit(w, r, util.KeyByIP(r)) {
		view.Error(w, auth_templ.LogIn(auth_templ.LogInProps{
			Username:     username,
			StayLoggedIn: persistent,
			Error:        "Too many log in attempts. Please try again later.",
		}), r.Context(), http.StatusTooManyRequests)
		return
	}

	u, err := user.FromCredentials(r.Context(), &user.Credentials{
		Username: username,
		Password: password,
	})

	if err == user.ErrInvalidCredentials {
		view.Error(w, auth_templ.LogIn(auth_templ.LogInProps{
			Username:     username,
			StayLoggedIn: persistent,
			Error:        "Invalid username or password.",
		}), r.Context(), http.StatusUnauthorized)
		return
	}

	if err != nil {
		slog.Error("user::FromCredentials failed", "err", err)
		view.Error(w, auth_templ.LogIn(auth_templ.LogInProps{
			Username:     username,
			StayLoggedIn: persistent,
			Error:        "Something broke on our end. Please try again later.",
		}), r.Context(), http.StatusInternalServerError)
		return
	}

	session, err := u.CreateSesion(r.Context())
	if err != nil {
		slog.Error("user::CreateSesion failed", "err", err, "user", u)
		view.Error(w, auth_templ.LogIn(auth_templ.LogInProps{
			Username:     username,
			StayLoggedIn: persistent,
			Error:        "Something broke on our end. Please try again later.",
		}), r.Context(), http.StatusInternalServerError)
		return
	}

	c := &http.Cookie{
		Name:     core_auth.SessionCookie,
		Value:    session.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // TODO: turn secure when available
		SameSite: http.SameSiteLaxMode,
	}

	if persistent {
		c.Expires = session.ExpiresAt
	}

	http.SetCookie(w, c)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
