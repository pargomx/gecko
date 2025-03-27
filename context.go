package gecko

import (
	"net/http"
	"net/url"
	"time"
)

// Representa a la solicitud HTTP actual
// y ofrece los medios para responderla.
type Context struct {
	request  *http.Request
	response *Response
	path     string // Patrón de ruta registrado. Ej: "GET /inicio"
	query    url.Values
	gecko    *Gecko
	SesionID string
	Sesion   any
	time     time.Time // Momento en el que se comenzó a procesar la solicitud, utilizado para el log http.
	Compress bool      // Activar compresión con gzip para Render y RenderOk.
}

func (c *Context) Request() *http.Request {
	return c.request
}

func (c *Context) Response() *Response {
	return c.response
}

func (c *Context) Time() time.Time {
	return c.time
}
