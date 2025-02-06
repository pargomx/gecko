package gkt

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (

	// Acepta espacios, letras y acentos.
	RegexLetrasAcentosEspacios = regexp.MustCompile(`^[a-zA-ZáéíóúÁÉÍÓÚ ]+$`)

	// Acepta letras minúsculas y guiones.
	RegexMinusculasGuiones = regexp.MustCompile(`^[a-z-]+$`)

	// Acepta letras minúsculas sin espacios.
	RegexMinusculas = regexp.MustCompile(`^[a-z]+$`)

	// Acepta letras minúsculas y mayúsculas, números y guiones.
	RegexAlfanumericoGuiones = regexp.MustCompile(`^[a-zA-Z1-9-]+$`)
)

// ================================================================ //
// ========== Espacios ============================================ //

// Cualquier tipo de espacio, incluyendo tabs y saltos de línea.
var RegexEspacios = regexp.MustCompile(`\s+`)

// SinEspaciosExtra sustituye todos los tabs, saltos de línea y espacios
// dobles por espacios sencillos entre todas las palabras, además de
// cortar todos los espacios al principio y al final de txt.
func SinEspaciosExtra(txt string) string {
	return strings.TrimSpace(RegexEspacios.ReplaceAllLiteralString(txt, " "))
}

// SinEspaciosNinguno elimina todos los tabs, saltos de línea y espacios.
func SinEspaciosNinguno(txt string) string {
	return strings.TrimSpace(RegexEspacios.ReplaceAllLiteralString(txt, " "))
}

// Cualquier tipo de espacio menos saltos de línea.
var RegexEspaciosNoLinebreak = regexp.MustCompile(`[^\S\n]+`)

// SinEspaciosExtra sustituye todos los tabs y espacios dobles
// por espacios sencillos entre todas las palabras, además de
// cortar todos los espacios al principio y al final de txt.
//
// Conserva los saltos de línea.
func SinEspaciosExtraConSaltos(txt string) string {
	return strings.TrimSpace(RegexEspaciosNoLinebreak.ReplaceAllLiteralString(txt, " "))
}

// ================================================================ //
// ========== Diacríticos ========================================= //

// Remueve la clase unicode Mark nonspacing que tiene diacríticos (acentos y más).
var diacriticsTransformer = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

// SinDiacriticos transforma todas las letras especiales o con diacríticos
// en su equivalente latín estándar. https://go.dev/blog/normalization
//
//	Ej. "Él miró al pingüino ò.ó del Barça"
//	 -> "El miro al pinguino o.o del Barca"
func SinDiacriticos(txt string) string {
	output, _, err := transform.String(diacriticsTransformer, txt)
	if err != nil {
		panic(err)
	}
	return output
}

// ================================================================ //
// ========== Puntuación ========================================== //

// Símbolos de puntuación como: [!"#$%&'()*+,-./:;<=>?@[\]^_{|}~¡¿¿`]
var RegexPuntuacion = regexp.MustCompilePOSIX(`[[:punct:]]+`)

func SinPuntuacion(txt string) string {
	return RegexPuntuacion.ReplaceAllLiteralString(txt, "")
}

// ================================================================ //

var RegxAlfanum = regexp.MustCompile(`[^a-zA-Z0-9_]`)

// SinRiesgoParaID *elimina* todos los caracteres que
// no sean alfanuméricos ascii o guión bajo.
//
// Último recurso para asegurar que el txt se transforme
// en algo seguro y apropiado para su en con cualquier
// sitema sin sorpresas unicode.
func SinRiesgoParaID(txt string) string {
	return RegxAlfanum.ReplaceAllLiteralString(txt, "")
}

// TODO: Agregar función que conserve espacios?
// TODO: Hacer que SinRiesgo haga primero SinDiacriticos y ToSnake?
