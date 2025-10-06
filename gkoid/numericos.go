package gkoid

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// randomBigIntNDigits genera un número aleatorio de n dígitos como *big.Int.
func randomBigIntNDigits(n int) (*big.Int, error) {
	if n <= 0 {
		return nil, fmt.Errorf("gkoid: n debe ser mayor que 0")
	}
	if n > 19 {
		return nil, fmt.Errorf("gkoid: id numérico no puede tener más de 19 dígitos")
	}
	if n == 1 {
		// Para un solo dígito, rango 1-9
		num, err := rand.Int(rand.Reader, big.NewInt(9))
		if err != nil {
			return nil, fmt.Errorf("gkoid: %w", err)
		}
		num.Add(num, big.NewInt(1))
		return num, nil
	}
	// Para n > 1, rango 10^(n-1) a 10^n - 1
	min := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(n-1)), nil)
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(n)), nil)
	diff := new(big.Int).Sub(max, min)
	num, err := rand.Int(rand.Reader, diff)
	if err != nil {
		return nil, fmt.Errorf("gkoid: %w", err)
	}
	num.Add(num, min)
	return num, nil
}

// GenerateInt genera un ID numérico aleatorio de n dígitos (int).
func GenerateInt(n int) (int, error) {
	var maxDigits int
	// Detect architecture: 32-bit or 64-bit
	if ^uint(0)>>32 == 0 {
		maxDigits = 9 // max int32 is 2147483647 (10 digits, but first digit can't be 2 for random)
	} else {
		maxDigits = 18 // max int64 is 9223372036854775807 (19 digits, but first digit can't be 9 for random)
	}
	if n > maxDigits {
		return 0, fmt.Errorf("gkoid.GenerateInt: no puede tener más de %d dígitos en esta arquitectura", maxDigits)
	}
	num, err := randomBigIntNDigits(n)
	if err != nil {
		return 0, err
	}
	return int(num.Int64()), nil
}

// GenerateInt64 genera un ID numérico aleatorio de n dígitos (int64).
func GenerateInt64(n int) (int64, error) {
	if n > 18 {
		return 0, fmt.Errorf("gkoid.GenerateInt: no puede tener más de 18 dígitos")
	}
	num, err := randomBigIntNDigits(n)
	if err != nil {
		return 0, err
	}
	return num.Int64(), nil
}

// GenerateUint genera un ID numérico aleatorio de n dígitos (uint).
func GenerateUint(n int) (uint, error) {
	var maxDigits int
	// Detect architecture: 32-bit or 64-bit
	if ^uint(0)>>32 == 0 {
		maxDigits = 9 // max uint32 is 4294967295 (10 digits, pero 9 para garantizar 999,999,999)
	} else {
		maxDigits = 19 // max uint64 is 18446744073709551615 (20 digits, but first digit can't be 1 for random)
	}
	if n > maxDigits {
		return 0, fmt.Errorf("gkoid.GenerateUint: no puede tener más de %d dígitos en esta arquitectura", maxDigits)
	}
	num, err := randomBigIntNDigits(n)
	if err != nil {
		return 0, err
	}
	return uint(num.Uint64()), nil
}

// GenerateUint64 genera un ID numérico aleatorio de n dígitos (uint64).
func GenerateUint64(n int) (uint64, error) {
	num, err := randomBigIntNDigits(n)
	if err != nil {
		return 0, err
	}
	return num.Uint64(), nil
}
