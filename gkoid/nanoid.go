package gkoid

import (
	"crypto/rand"
	"math"

	"github.com/pargomx/gecko/gko"
)

// ============================================================ //
//             Generador de ids aleatorios - nanoid				//
// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++	//
//  Código original obtenido de:								//
//  Matous Dzivjak <matousdzivjak@gmail.com>					//
//  MIT License: https://github.com/matoous/go-nanoid			//
// ============================================================ //

var (
	alfabetoPwd = []rune("_-0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!#$&%")
	alfabeto64  = []rune("_-0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	alfabeto62  = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	alfabeto36  = []rune("0123456789abcdefghijklmnopqrstuvwxyz")
	alfabeto16  = []rune("0123456789ABCDEF")
	alfabetoNum = []rune("0123456789")
)

// By performing a bitwise AND with 63 (00111111), the result of bytes[i]&63
// will always be a value between 0 and 63 to limit the possible values of
// bytes[i] to the range of indices valid for alfabeto64 with 64 elements.
const (
	maskPwd = 68
	mask64  = 63
	mask62  = 61
	mask36  = 35
	mask16  = 15
	maskNum = 9
)

// Genera una cadena aleatoria del largo especificado con el alfabeto:
//
//	"_-0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!#$&%" // 69 chars
func NewPwd(size int) (string, error) {
	if size <= 0 {
		return "", gko.Op("gkoid.NewPwd").Str("size must be positive")
	}
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", gko.Op("gkoid.NewPwd").Err(err)
	}
	id := make([]rune, size)
	for i := range size {
		id[i] = alfabetoPwd[bytes[i]&maskPwd]
	}
	return string(id[:size]), nil
}

// Genera una cadena aleatoria del largo especificado con el alfabeto:
//
//	"_-0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // 64 chars
func New64(size int) (string, error) {
	if size <= 0 {
		return "", gko.Op("gkoid.New64").Str("size must be positive")
	}
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", gko.Op("gkoid.New64").Err(err)
	}
	id := make([]rune, size)
	for i := range size {
		id[i] = alfabeto64[bytes[i]&mask64]
	}
	return string(id[:size]), nil
}

// Genera una cadena aleatoria del largo especificado con el alfabeto:
//
//	"0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // 62 chars
func New62(size int) (string, error) {
	if size <= 0 {
		return "", gko.Op("gkoid.New62").Str("size must be positive")
	}
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", gko.Op("gkoid.New62").Err(err)
	}
	id := make([]rune, size)
	for i := range size {
		id[i] = alfabeto62[bytes[i]&mask62]
	}
	return string(id[:size]), nil
}

// Genera una cadena aleatoria del largo especificado con el alfabeto:
//
//	"0123456789abcdefghijklmnopqrstuvwxyz" // 36 chars
func New36(size int) (string, error) {
	if size <= 0 {
		return "", gko.Op("gkoid.New36").Str("size must be positive")
	}
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", gko.Op("gkoid.New36").Err(err)
	}
	id := make([]rune, size)
	for i := range size {
		id[i] = alfabeto36[bytes[i]&mask36]
	}
	return string(id[:size]), nil
}

// Genera una cadena aleatoria del largo especificado con el alfabeto:
//
//	"0123456789ABCDEF" // 16 chars
func New16(size int) (string, error) {
	if size <= 0 {
		return "", gko.Op("gkoid.New16").Str("size must be positive")
	}
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", gko.Op("gkoid.New16").Err(err)
	}
	id := make([]rune, size)
	for i := range size {
		id[i] = alfabeto16[bytes[i]&mask16]
	}
	return string(id[:size]), nil
}

// Genera una cadena aleatoria del largo especificado con el alfabeto:
//
//	"0123456789" // 10 chars
func NewNum(size int) (string, error) {
	if size <= 0 {
		return "", gko.Op("gkoid.NewNum").Str("size must be positive")
	}
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", gko.Op("gkoid.NewNum").Err(err)
	}
	id := make([]rune, size)
	for i := range size {
		id[i] = alfabetoNum[bytes[i]&maskNum]
	}
	return string(id[:size]), nil
}

// ================================================================ //
// ========== ALFABETO CUSTOM ===================================== //

// bit mask usada para obtener índice válido de un caracter en alfabeto
// a partir de los random bytes generados.
func getMask(alphabetSize int) int {
	for i := 1; i <= 8; i++ {
		mask := (2 << uint(i)) - 1
		if mask >= alphabetSize-1 {
			return mask
		}
	}
	return 0
}

// Genera una cadena aleatoria del largo especificado con un alfabeto
// que debe tener de 1 a 255 caracteres.
func NewCustom(alphabet string, size int) (string, error) {
	chars := []rune(alphabet)

	if len(alphabet) == 0 || len(alphabet) > 255 {
		return "", gko.Op("gkoid.NewCustom").Str("alphabet must not be empty and contain no more than 255 chars")
	}
	if size <= 0 {
		return "", gko.Op("gkoid.NewCustom").Str("size must be positive")
	}

	mask := getMask(len(chars))
	// estimate how many random bytes we will need for the ID, we might actually need more but this is tradeoff
	// between average case and worst case
	ceilArg := 1.6 * float64(mask*size) / float64(len(alphabet))
	step := int(math.Ceil(ceilArg))

	id := make([]rune, size)
	bytes := make([]byte, step)
	for j := 0; ; {
		_, err := rand.Read(bytes)
		if err != nil {
			return "", err
		}
		for i := 0; i < step; i++ {
			currByte := bytes[i] & byte(mask)
			if currByte < byte(len(chars)) {
				id[j] = chars[currByte]
				j++
				if j == size {
					return string(id[:size]), nil
				}
			}
		}
	}
}
