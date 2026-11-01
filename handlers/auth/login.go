package auth

import (
	"net/http"

	"github.com/nedpranson/the-christmas-breach/core/util/view"
	auth_templ "github.com/nedpranson/the-christmas-breach/templates/auth"
)

func LogIn(w http.ResponseWriter, r *http.Request) {
	view.Render(w, auth_templ.LogIn());
}
