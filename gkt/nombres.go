package gkt

import (
	"strings"
)

// nombrePropio convierte texto a formato de nombre propio. Si el texto está todo
// en mayúsculas o minúsculas, aplica mayúscula solo a la primera letra de cada
// palabra excepto "la".
func NombrePropio(txt string) string {
	if txt == "" {
		return ""
	}

	// Si todo está en mayúsculas o minúsculas, convertir a formato de nombre propio
	if strings.ToUpper(txt) == txt || strings.ToLower(txt) == txt {
		words := strings.Fields(strings.ToLower(txt))
		for i, word := range words {
			if word != "la" {
				words[i] = PrimeraMayusc(word)
			}
		}
		return strings.Join(words, " ")
	}

	return PrimeraMayusc(txt)
}
