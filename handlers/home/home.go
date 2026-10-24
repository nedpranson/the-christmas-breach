package home

import (
	"net/http"

	"github.com/nedpranson/the-christmas-breach/core/util/view"
	home_templ "github.com/nedpranson/the-christmas-breach/templates/home"
)

func Home(w http.ResponseWriter, r *http.Request) {
	view.Render(w, home_templ.Home());
}
