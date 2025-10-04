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
	if n > 18 {
		return nil, fmt.Errorf("gkoid: id numérico no puede tener más de 18 dígitos")
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

// RandomInt genera un ID numérico aleatorio de n dígitos (int).
func RandomInt(n int) (int, error) {
	num, err := randomBigIntNDigits(n)
	if err != nil {
		return 0, err
	}
	return int(num.Int64()), nil
}

// RandomInt64 genera un ID numérico aleatorio de n dígitos (int64).
func RandomInt64(n int) (int64, error) {
	num, err := randomBigIntNDigits(n)
	if err != nil {
		return 0, err
	}
	return num.Int64(), nil
}

// RandomUint genera un ID numérico aleatorio de n dígitos (uint).
func RandomUint(n int) (uint, error) {
	num, err := randomBigIntNDigits(n)
	if err != nil {
		return 0, err
	}
	return uint(num.Uint64()), nil
}

// RandomUint64 genera un ID numérico aleatorio de n dígitos (uint64).
func RandomUint64(n int) (uint64, error) {
	num, err := randomBigIntNDigits(n)
	if err != nil {
		return 0, err
	}
	return num.Uint64(), nil
}
