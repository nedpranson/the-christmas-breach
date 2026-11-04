package ratelimit

import (
	"net"
	"net/http"
	"time"

	"github.com/go-chi/httprate"
	"github.com/nedpranson/the-christmas-breach/core/util/view"
	auth_templ "github.com/nedpranson/the-christmas-breach/templates/auth"
)

func LogInRateLimitMiddleware(f http.HandlerFunc) http.Handler {
    return httprate.LimitBy(
        10,
        30*time.Minute,
		keyByIp,
        httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
            view.Error(w, auth_templ.LogIn(auth_templ.LogInProps{
                Username:     r.FormValue("username"),
                StayLoggedIn: r.FormValue("persistent") != "",
                Error:        "Too many log in attempts. Please try again later.",
            }), r.Context(), http.StatusTooManyRequests)
        }),
    )(http.HandlerFunc(f))
}

func SignUpRateLimitMiddleware(f http.HandlerFunc) http.Handler {
    return httprate.LimitBy(
        3,
        30*time.Minute,
		keyByIp,
        httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
            view.Error(w, auth_templ.SignUp(auth_templ.SignUpProps{
                Username:     r.FormValue("username"),
                StayLoggedIn: r.FormValue("persistent") != "",
                Error:        "Too many sign up attempts. Please try again later.",
            }), r.Context(), http.StatusTooManyRequests)
        }),
    )(http.HandlerFunc(f))
}

 func keyByIp(r *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}
	return httprate.CanonicalizeIP(ip), nil
 }
