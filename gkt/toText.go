package gkt

import (
	"regexp"
	"strings"
	"unicode"
)

/*
	Todas estas funciones quitan los espacios extra y dejan
	el resultado en una sola línea, pero conservan los diacríticos.
*/

// ================================================================ //
// ========== Change Case ========================================= //

// Símbolos comunes que pueden separar una palabra, como espacios,
// guiones y otros símbolos de puntuación.
var RegexGuiones = regexp.MustCompile(`[-_—–―]+`)

const guion = "-"

// Sustituye espacios, guiones y símbolos de puntuación para dejar
// solamente un guión "-" entre cada palabra.
//
// Para controlar mayúsculas o acentos encadenar funciones correspondientes.
//
//	Ej: "palabas separadas así! _humm" -> "palabas-separadas-así!--humm"
func ToKebab(txt string) string {
	txt = RegexEspacios.ReplaceAllLiteralString(txt, guion)
	txt = RegexPuntuacion.ReplaceAllLiteralString(txt, guion)
	txt = RegexGuiones.ReplaceAllLiteralString(txt, guion)
	txt = strings.Trim(txt, guion)
	// txt = strings.ReplaceAll(txt, " ", "-")
	// txt = strings.ReplaceAll(txt, "_", "-")
	return txt
}

const underscore = "_"

// Sustituye espacios, guiones y símbolos de puntuación para dejar
// solamente un underscore "_" entre cada palabra.
//
// Para controlar mayúsculas o acentos encadenar funciones correspondientes.
//
//	Ej: "palabas separadas así!  -humm" -> "palabas_separadas_así_humm"
func ToSnake(txt string) string {
	txt = RegexEspacios.ReplaceAllLiteralString(txt, underscore)
	txt = RegexPuntuacion.ReplaceAllLiteralString(txt, underscore)
	txt = RegexGuiones.ReplaceAllLiteralString(txt, underscore)
	txt = strings.Trim(txt, underscore)
	return txt
}

// Une todas las palabras y marca los inicios de cada
// uno con mayúsculas.
//
//	Ej: "el pájaro azul -dice" -> "ElPájaroAzulDice"
func ToCamel(txt string) string {
	palabras := strings.Split(ToSnake(txt), "_")
	txt = ""
	for _, palabra := range palabras {
		if palabra == "" {
			continue
		}
		txt = txt + strings.ToUpper(palabra[0:1]) + palabra[1:]
	}
	return txt
}

// ================================================================ //
// ========== Mayúsculas ========================================== //

// Retorna txt en una sola línea sin espacios extras en mayúsculas.
func ToUpper(txt string) string {
	return strings.ToUpper(SinEspaciosExtra(txt))
}

// Retorna txt en una sola línea sin espacios extras en minúsculas.
func ToLower(txt string) string {
	return strings.ToLower(SinEspaciosExtra(txt))
}

// Transforma la primera letra en mayúscula.
func PrimeraMayusc(txt string) string {
	if txt == "" {
		return ""
	}
	r := []rune(txt)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// Transforma la primera letra en minúscula
// a menos que la segunda también sea mayúscula
// como es el caso de un acrónimo.
//
//	Ej: "El barco" -> "el barco". "MXN si" -> "MXN si"
func SinPrimeraMayusc(txt string) string {
	if txt == "" {
		return ""
	}
	r := []rune(txt)
	// Omitir si la segunda letra es mayúscula.
	if len(r) > 2 &&
		unicode.ToUpper(r[1]) == r[1] {
		return txt
	}
	r[0] = unicode.ToLower(r[0])
	return string(r)
}
