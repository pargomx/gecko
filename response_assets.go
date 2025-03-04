package gecko

import (
	"crypto/sha1"
	_ "embed"
	"fmt"

	"github.com/pargomx/gecko/gko"
)

// Recurso estático
type staticResource struct {
	content  []byte
	etag     string // sha1 entre dobles comillas
	mimeType string
	lenght   string // content lenght in bytes
}

// Agrega un recurso estático a gecko para servirlo
// utilizando su ETag obtenido con sha1.
func (g *Gecko) AgregarRecurso(name string, content []byte, mimeType string) {
	err := g.agregarRecurso(name, content, mimeType)
	if err != nil {
		gko.FatalError(err)
	}
}

func (g *Gecko) agregarRecurso(name string, content []byte, mimeType string) error {
	op := gko.Op("gecko.AgregarRecurso")
	if name == "" {
		return op.Str("nombre no especificado")
	}
	if len(content) == 0 {
		return op.Str("recurso sin contenido")
	}
	if mimeType == "" {
		return op.Str("mimeType sin contenido")
	}
	h := sha1.New()
	_, err := h.Write(content)
	if err != nil {
		return op.Err(err)
	}
	if g.staticFiles == nil {
		g.staticFiles = make(map[string]staticResource)
	}
	if _, exists := g.staticFiles[name]; exists {
		return op.Strf("recurso '%s' registrado doble", name)
	}
	g.staticFiles[name] = staticResource{
		content:  content,
		etag:     fmt.Sprintf("\"%x\"", h.Sum(nil)),
		mimeType: mimeType,
		lenght:   fmt.Sprintf("%d", len(content)),
	}
	return nil
}

// Sirve un recurso estático registrado anteriormente en gecko
// utilizando su ETag para el control de caché.
func (g *Gecko) ServirRecurso(name string) HandlerFunc {
	res, exists := g.staticFiles[name]
	if !exists {
		gko.FatalExitf("gecko.ServirRecurso: recurso '%v' no registrado", name)
	}
	return func(c *Context) error {

		// Enviar Content-Lenght incluso en 304.
		c.response.Header().Set(HeaderContentLength, res.lenght)

		// Enviar 304 si el Etag coincide con 'If-None-Match'.
		if match := c.request.Header.Get("If-None-Match"); match != "" {
			if match == res.etag {
				return c.NoContent(Status304NotModified)
			}
		}

		// Enviar contenido con su Etag y Cache-Control.
		c.response.Header().Set("Etag", res.etag)
		c.response.Header().Set(HeaderCacheControl, "public, max-age=3600, stale-while-revalidate=18000, stale-if-error=18000")
		c.response.Header().Set(HeaderContentType, res.mimeType)
		return c.ContentOk(res.mimeType, res.content)
	}
}

// ================================================================ //
// ========== Recursos estándar =================================== //

//go:embed javascript/gecko.js
var geckoJsExtension []byte

//go:embed javascript/htmx.js
var htmxJs []byte

//go:embed javascript/htmx.min.js
var htmxMinJs []byte

// Extensión gecko para htmx.
// Colocar en una ruta como "/gecko.js" y activar con hx-ext="gecko".
func (g *Gecko) ServirGeckoJS() HandlerFunc {
	g.AgregarRecurso("gecko.js", geckoJsExtension,
		MIMEApplicationJavaScriptCharsetUTF8)
	return g.ServirRecurso("gecko.js")
}

// Librería HTMX.js
func (g *Gecko) ServirHtmxJS() HandlerFunc {
	g.AgregarRecurso("htmx.js", htmxJs,
		MIMEApplicationJavaScriptCharsetUTF8)
	return g.ServirRecurso("htmx.js")
}

// Librería HTMX.js
func (g *Gecko) ServirHtmxMinJS() HandlerFunc {
	g.AgregarRecurso("htmx.min.js", htmxMinJs,
		MIMEApplicationJavaScriptCharsetUTF8)
	return g.ServirRecurso("htmx.min.js")
}
