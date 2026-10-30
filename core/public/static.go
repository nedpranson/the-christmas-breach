//go:build !dev

package public

import public_embed "github.com/nedpranson/the-christmas-breach/public"

func NewFileServer() *FileServer {
	return &FileServer{fs: public_embed.FS}
}
