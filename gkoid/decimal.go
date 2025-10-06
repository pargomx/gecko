package gkoid

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
)

type Decimal uint64

// ================================================================ //

func (id Decimal) Uint64() uint64 {
	return uint64(id)
}

func (id Decimal) String() string {
	return fmt.Sprintf("%d", id)
}

// ================================================================ //

func ParseDecimal(s string) (Decimal, error) {
	num, err := strconv.ParseUint(s, 10, 64)
	return Decimal(num), err
}

// ================================================================ //

// NewDecimal genera un ID aleatorio con n dígitos decimales.
// Mínimo 3 y máximo 19.
//
//	NewDecimal(4) rango 1000 - 9999
//	NewDecimal(12) rango 100,000,000,000 - 999,999,999,999
//
// La idea es sustituir con un value-object el simplemente hacer:
//
//	rand.IntN(9000000) + 1000000
func NewDecimal(digitos int) (Decimal, error) {
	if digitos > 19 {
		return 0, fmt.Errorf("gkoid: debe tener máximo 19 dígitos decimales")
	}
	if digitos < 3 {
		return 0, fmt.Errorf("gkoid: debe tener mínimo 3 dígitos decimales")
	}
	// Para n > 1, rango 10^(n-1) a 10^n - 1
	min := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(digitos-1)), nil)
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(digitos)), nil)
	diff := new(big.Int).Sub(max, min)
	// Ej. NewDecimalID(3): min=100 max=1000 diff=900 max_rand=899 added=999
	num, err := rand.Int(rand.Reader, diff)
	if err != nil {
		return 0, fmt.Errorf("gkoid.NewDecimalID: %w", err)
	}
	num.Add(num, min)
	return Decimal(num.Uint64()), nil
}

// NewDecimal genera un ID aleatorio con n dígitos decimales.
// Mínimo 3 y máximo 19.
// Panics if error.
func NewDecimalMust(digitos int) Decimal {
	id, err := NewDecimal(digitos)
	if err != nil {
		panic(err)
	}
	return id
}
