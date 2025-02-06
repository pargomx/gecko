package gkt

import (
	"strconv"
	"strings"

	"github.com/pargomx/gecko/gko"
)

// Retorna false a menos de que el valor sea:
// "on", "true", "1".
func ToBool(txt string) bool {
	txt = strings.ToLower(SinEspaciosNinguno(txt))
	return txt == "on" || txt == "true" || txt == "1"
}

// Retorna el valor en tipo entero.
// Retorna error si no es un número válido.
func ToInt(txt string) (int, error) {
	return strconv.Atoi(SinEspaciosNinguno(txt))
}

// Retorna el valor en tipo entero positivo de 8 bytes.
// Valor máximo aceptado: 18446744073709551615.
func ToUint64(txt string) (uint64, error) {
	return strconv.ParseUint(SinEspaciosNinguno(txt), 10, 64)
}

// Devuelve los centavos a partir de un string de dinero
// que puede ser "$200.00", "200", "200.0" por ejemplo.
//
// El valor recibido debe estar en unidades de pesos.
// Se puede incluir centavos pero deben estar como decimales.
func ToCentavos(txt string) (int, error) {
	str := SinEspaciosNinguno(txt)
	str = strings.ReplaceAll(str, " ", "")
	str = strings.ReplaceAll(str, ",", "")
	str = strings.TrimLeft(str, "$")

	partes := strings.Split(str, ".")
	if len(partes) > 2 {
		return 0, gko.ErrDatoInvalido().Msgf("más de 1 punto para centavos: '%s'", str)
	}
	pesos := partes[0]
	centavos := ""
	if len(partes) == 2 {
		centavos = partes[1]
	}

	switch len(centavos) {
	case 0:
		str = pesos + "00"
	case 1:
		str = pesos + centavos + "0"
	case 2:
		str = pesos + centavos
	default:
		return 0, gko.ErrDatoInvalido().
			Msgf("solo puede haber centavos luego del punto: '%s'", centavos)
	}

	res, err := strconv.Atoi(str)
	if err != nil {
		return 0, gko.ErrDatoInvalido().Msgf("número inválido: '%s'", str)
	}
	return res, nil
}
