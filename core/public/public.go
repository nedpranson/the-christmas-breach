package public

import (
	"errors"
	"io"
	"io/fs"
	"net/http"
	"strings"

	"github.com/nedpranson/the-christmas-breach/core/util/view"
	base_templ "github.com/nedpranson/the-christmas-breach/templates/base"
)

type FileServer struct {
	fs fs.FS
}

func (s *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/")

	file, err := s.fs.Open(name)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) || errors.Is(err, fs.ErrInvalid) {
			view.Error(w, base_templ.ErrorNotFound(), http.StatusNotFound)
			return
		}
		view.Error(w, base_templ.ErrorInternalServer(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		view.Error(w, base_templ.ErrorInternalServer(), http.StatusInternalServerError)
		return
	}

	if info.IsDir() {
		view.Error(w, base_templ.ErrorNotFound(), http.StatusNotFound)
		return
	}

	http.ServeContent(w, r, info.Name(), info.ModTime(), file.(io.ReadSeeker))
}
