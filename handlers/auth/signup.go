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

var signUpRateLimiter = httprate.NewRateLimiter(3, 30*time.Minute)

func SignUp(w http.ResponseWriter, r *http.Request) {
	view.Render(w, auth_templ.SignUp(), r.Context())
}

func SignUpPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	passwordVerify := r.FormValue("password-verify")
	persistent := r.FormValue("persistent") != ""

	if username == "" || password == "" || passwordVerify == "" {
		view.Error(w, base_templ.Error(base_templ.ErrorProps{
			Message: "Username and password are required.",
			Status:  http.StatusBadRequest,
		}), r.Context(), http.StatusBadRequest)
		return
	}

	if len(password) > 72 {
		view.Error(w, auth_templ.SignUp(auth_templ.SignUpProps{
			Username:     username,
			StayLoggedIn: persistent,
			Error:        "Password is too long.",
		}), r.Context(), http.StatusBadRequest)
		return
	}

	if password != passwordVerify {
		view.Error(w, auth_templ.SignUp(auth_templ.SignUpProps{
			Username:     username,
			StayLoggedIn: persistent,
			Error:        "Passwords don't match.",
		}), r.Context(), http.StatusBadRequest)
		return
	}

	if signUpRateLimiter.OnLimit(w, r, util.KeyByIP(r)) {
		view.Error(w, auth_templ.SignUp(auth_templ.SignUpProps{
			Username:     username,
			StayLoggedIn: persistent,
			Error:        "Too many sign up attempts. Please try again later.",
		}), r.Context(), http.StatusTooManyRequests)
		return
	}

	u, err := user.Create(r.Context(), &user.CreateUserParams{
		Username: username,
		Pasword:  password,
	})

	if err == user.ErrInvalidUsername {
		view.Error(w, auth_templ.SignUp(auth_templ.SignUpProps{
			Username:     username,
			StayLoggedIn: persistent,
			Error:        "Username can only contain English letters, numbers, and symbols.",
		}), r.Context(), http.StatusBadRequest)
		return
	}

	if err == user.ErrUsernameTaken {
		view.Error(w, auth_templ.SignUp(auth_templ.SignUpProps{
			Username:     username,
			StayLoggedIn: persistent,
			Error:        "Username already taken.",
		}), r.Context(), http.StatusConflict)
		return
	}

	if err != nil {
		slog.Error("user::Create failed", "err", err)
		view.Error(w, auth_templ.SignUp(auth_templ.SignUpProps{
			Username:     username,
			StayLoggedIn: persistent,
			Error:        "Something broke on our end. Please try again later.",
		}), r.Context(), http.StatusInternalServerError)
		return
	}

	session, err := u.CreateSesion(r.Context())
	if err != nil {
		// todo: notify that account was created just we could not login
		slog.Error("user::Create failed", "err", err)
		view.Error(w, auth_templ.SignUp(auth_templ.SignUpProps{
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
