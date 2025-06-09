package gko

import "crypto/rand"

var alfabeto62 = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

const mask62 = 61

func newID(size int) (string, error) {
	if size <= 0 {
		return "", Op("gkoid.New62").Str("size must be positive")
	}
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", Op("gkoid.New62").Err(err)
	}
	id := make([]rune, size)
	for i := 0; i < size; i++ {
		id[i] = alfabeto62[bytes[i]&mask62]
	}
	return string(id[:size]), nil
}
