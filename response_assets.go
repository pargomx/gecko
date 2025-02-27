package gecko

import _ "embed"

//go:embed gecko.js
var geckoJsExtension string

// Extensión gecko para htmx.
// Colocar en una ruta como "/gecko.js" y activar con hx-ext="gecko".
func GeckoJS(c *Context) error {
	c.response.Header().Set(HeaderContentType, MIMEApplicationJavaScriptCharsetUTF8)
	return c.StringOk(geckoJsExtension)
}
