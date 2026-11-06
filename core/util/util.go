package util

import (
	"crypto/rand"
	"net"
	"net/http"

	"github.com/go-chi/httprate"
	twmerge "github.com/Oudwins/tailwind-merge-go"
)

func CryptoRandomBytes(len int) []byte {
	buf := make([]byte, len)
	n, err := rand.Read(buf)
	if err != nil || n != len {
		panic(err)
	}
	return buf
}

func TwMerge(classes ...string) string {
	return twmerge.Merge(classes...)
}


func KeyByIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}
	return httprate.CanonicalizeIP(ip)
}
