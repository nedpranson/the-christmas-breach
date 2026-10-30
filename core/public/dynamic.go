//go:build dev

package public

import "os"

func NewFileServer() *FileServer {
	return &FileServer{fs: os.DirFS("public")}
}
