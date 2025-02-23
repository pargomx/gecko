package plantillas

import (
	"html/template"
	"regexp"
	"strings"
)

var (
	reBold   = regexp.MustCompile(`\*(.*?)\*`)
	reItalic = regexp.MustCompile(`_(.*?)_`)
	reCode   = regexp.MustCompile("`(.*?)`")
)

// Enfatizar busca texto entre asteriscos y lo convierte en negritas,
// además de texto entre guiones bajos y lo convierte en cursivas.
func Enfatizar(text string) template.HTML {
	escapedText := template.HTMLEscapeString(text)
	escapedText = reBold.ReplaceAllStringFunc(escapedText, func(match string) string {
		return "<strong>" + strings.Trim(match, "*") + "</strong>"
	})
	escapedText = reItalic.ReplaceAllStringFunc(escapedText, func(match string) string {
		return "<em>" + strings.Trim(match, "_") + "</em>"
	})
	escapedText = reCode.ReplaceAllStringFunc(escapedText, func(match string) string {
		return "<code>" + strings.Trim(match, "`") + "</code>"
	})
	return template.HTML(escapedText)
}
