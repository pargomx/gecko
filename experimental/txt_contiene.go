package experimental

import "strings"

// ========================================================================== //
// ========================================================================== //

// Para mejorar lectura en comparaciones repetitivas.
type CompareString string

// Verdadero si substr está dentro de CompareString.
func (s CompareString) Contiene(substr string) bool {
	return strings.Contains(string(s), substr)
}

// Verdadero si comienza por prefix
func (s CompareString) HasPrefix(prefix string) bool {
	return strings.HasPrefix(string(s), prefix)
}

// Verdadero si termina por suffix
func (s CompareString) HasSuffix(suffix string) bool {
	return strings.HasSuffix(string(s), suffix)
}

// ExtraerEnmedio retorna el string que resulta
// de quitar lo que está a partir de izq y der.
func (s CompareString) ExtraerEnmedio(izq, der string) string {
	// str="uno %% hola && dos"
	// izq="% "   der=" &"
	str := string(s)
	spli := strings.Split(str, izq)

	if len(spli) < 2 { // Cuando no hay qué cortar a la izquierda.
		return strings.Split(str, der)[0]
	}
	res := spli[1] // res="hola && dos"

	res = strings.Split(res, der)[0] // res="hola"
	return res
}

// ExtraerEnmedio retorna el string que resulta
// de quitar lo que está a partir de izq y der.
func ExtraerEnmedio(str, izq, der string) string {
	// str="uno %% hola && dos"
	// izq="% "   der=" &"

	spli := strings.Split(str, izq)

	if len(spli) < 2 { // Cuando no hay qué cortar a la izquierda.
		return strings.Split(str, der)[0]
	}
	res := spli[1] // res="hola && dos"

	res = strings.Split(res, der)[0] // res="hola"
	return res
}
