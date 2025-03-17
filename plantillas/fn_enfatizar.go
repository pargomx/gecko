package plantillas

import (
	"html/template"
	"net/url"
	"regexp"
	"strings"

	"github.com/pargomx/gecko/gkt"
)

var (
	reBold           = regexp.MustCompile(`\*(.*?)\*`)
	reItalic         = regexp.MustCompile(`_(.*?)_`)
	reCode           = regexp.MustCompile("`(.*?)`")
	reCodeMultilinea = regexp.MustCompile("```([\\s\\S]*?)```")
	reLink           = regexp.MustCompile(`\[(.*?)\]\((.*?)\)`) // Link markdown [title](https://example.com)
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
	escapedText = reCodeMultilinea.ReplaceAllStringFunc(escapedText, func(match string) string {
		return "<code style=\"white-space: pre; display: inline-block;\">" + strings.Trim(match, "`") + "</code>"
	})
	escapedText = reCode.ReplaceAllStringFunc(escapedText, func(match string) string {
		return "<code>" + strings.Trim(match, "`") + "</code>"
	})
	escapedText = reLink.ReplaceAllStringFunc(escapedText, func(match string) string {
		matches := reLink.FindStringSubmatch(match)
		if len(matches) == 3 {
			desc := matches[1]
			unsafeLink := gkt.SinEspaciosNinguno(matches[2])
			if unsafeLink == "" {
				return `<a href="#bad_url" class="hover:underline">` + desc + ` [undefined link]</a>`
			}
			URL, err := url.Parse(unsafeLink)
			if err != nil {
				return `<a href="#bad_url" class="hover:underline">` + desc + " [bad link: " + err.Error() + `]</a>`
			}
			// Si no es absoluta del mismo origen entonces debe llevar https o http.
			if !strings.HasPrefix(unsafeLink, "/") && URL.Scheme != "http" && URL.Scheme != "https" {
				URL.Scheme = "https"
			}
			if desc == "" {
				desc = URL.String()
			}
			return `<a href="` + URL.String() + `" class="hover:underline">` + desc + `</a>`
		}
		return match
	})
	return template.HTML(escapedText)
}
