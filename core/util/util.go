package util

import "crypto/rand"

func CryptoRandomBytes(len int) []byte {
	buf := make([]byte, len)
	n, err := rand.Read(buf)
	if err != nil || n != len {
		panic(err)
	}
	return buf
}
