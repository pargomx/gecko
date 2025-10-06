package gkoid

import (
	"crypto/rand"
	"errors"
	"strings"
)

// Usa alfabeto ASCII 0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ (62 caracteres)
type Alfanum uint64

// ================================================================ //

func (id Alfanum) Uint64() uint64 {
	return uint64(id)
}

func (id Alfanum) String() string {
	if id == 0 {
		return ""
	}
	var b strings.Builder
	n := uint64(id)
	for n > 0 {
		b.WriteByte(byte(alfabeto62[n%62]))
		n /= 62
	}
	// Reverse the string
	s := b.String()
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// ================================================================ //

func ParseAlfanum(s string) (Alfanum, error) {
	if len(s) > 10 {
		return 0, errors.New("gkoid.ParseAlfanum: too long")
	}
	var n uint64
	for _, r := range s {
		idx := strings.IndexRune(string(alfabeto62), r)
		if idx == -1 {
			return 0, errors.New("gkoid.ParseAlfanum: invalid character")
		}
		n = n*62 + uint64(idx)
	}
	return Alfanum(n), nil
}

// ================================================================ //

// Generates a random AlfanumID consisting of the specified number of base62 characters.
// The number of characters (chars) must be between 1 and 10, inclusive.
// This limit is enforced because a uint64 can represent at most 10 base62 digits without overflow.
// Returns an error if chars is outside the valid range or if random data cannot be read.
func NewAlfanum(chars int) (Alfanum, error) {
	if chars < 1 || chars > 10 {
		return 0, errors.New("gkoid: chars must be between 1 and 10")
	}
	bytes := make([]byte, chars)
	_, err := rand.Read(bytes)
	if err != nil {
		return 0, err
	}
	var id uint64
	for i := range chars {
		idx := bytes[i] % 62
		id = id*62 + uint64(idx)
	}
	return Alfanum(id), nil
}

// Genera un ID aleatorio con n caracteres alfanuméricos (entre 1 y 10).
// Panics if error.
func NewAlfanumMust(chars int) Alfanum {
	id, err := NewAlfanum(chars)
	if err != nil {
		panic(err)
	}
	return id
}
