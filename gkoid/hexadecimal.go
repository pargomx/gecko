package gkoid

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

// Identificador que se presenta al usuario como hexadecimal de n caracteres
// pero se almacena y procesa como uint64.
type Hex uint64

// ================================================================ //

func (id Hex) Uint64() uint64 {
	return uint64(id)
}

func (id Hex) String() string {
	return fmt.Sprintf("%x", uint64(id))
}

// func (id Hex) String() string {
// 	var b [8]byte // 64 bits = 8 bytes
// 	n := 8
// 	val := uint64(id)
// 	for i := 7; i >= 0; i-- {
// 		b[i] = byte(val & 0xff)
// 		val >>= 8
// 	}
// 	for n > 1 && b[8-n] == 0 {
// 		n-- // Trim leading zeros
// 	}
// 	return hex.EncodeToString(b[8-n:])
// }

// ================================================================ //

func ParseHex(s string) (Hex, error) {
	if len(s) == 0 {
		return 0, errors.New("gkoid.ParseHex: empty string")
	}
	if len(s) > 16 {
		return 0, errors.New("gkoid.ParseHex: string too long")
	}
	var id uint64
	for _, c := range s {
		var v byte
		switch {
		case '0' <= c && c <= '9':
			v = byte(c - '0')
		case 'a' <= c && c <= 'f':
			v = byte(c - 'a' + 10)
		case 'A' <= c && c <= 'F':
			v = byte(c - 'A' + 10)
		default:
			return 0, errors.New("gkoid.ParseHex: invalid hex digit")
		}
		id = (id << 4) | uint64(v)
	}
	return Hex(id), nil
}

// ================================================================ //

// Genera un ID aleatorio con n dígitos hexadecimales (máximo 16).
func NewHex(digitos int) (Hex, error) {
	if digitos < 1 {
		return 0, fmt.Errorf("gkoid: debe tener mínimo 1 dígito")
	}
	if digitos > 16 {
		return 0, fmt.Errorf("gkoid: no puede tener más de 16 dígitos hex")
	}
	n := digitos * 4 // Cada dígito hexadecimal representa 4 bits
	if n > 64 {
		return 0, fmt.Errorf("gkoid.NewHexID: demasiados bits")
	}
	num, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), uint(n)))
	if err != nil {
		return 0, fmt.Errorf("gkoid.NewHexID: %w", err)
	}
	return Hex(num.Uint64()), nil
}

// Genera un ID aleatorio con n dígitos hexadecimales (máximo 16).
// Panics if error.
func NewHexMust(digitos int) Hex {
	id, err := NewHex(digitos)
	if err != nil {
		panic(err)
	}
	return id
}
