package gecko

import (
	"fmt"
	"net/http"

	"github.com/pargomx/gecko/gko"
)

// ================================================================ //
// ========== Respuestas satisfactorias (2xx) ===================== //

func (c *Context) StringOk(msg string) error {
	c.response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.response.WriteHeader(200)
	c.response.Writer.Write([]byte(msg))
	return nil
}

// String sends a string response with status code 200 OK.
func (c *Context) StatusOk(msg string) (err error) {
	return c.Blob(http.StatusOK, MIMETextPlainCharsetUTF8, []byte(msg))
}

// String sends a string response with status code 200 OK.
func (c *Context) StatusOkf(format string, a ...any) (err error) {
	return c.Blob(http.StatusOK, MIMETextPlainCharsetUTF8, []byte(fmt.Sprintf(format, a...)))
}

// Retorna un estatus 202 aceptado con el mensaje dado.
func (c *Context) StatusAccepted(msg string) error {
	c.response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.response.WriteHeader(202)
	_, err := c.response.Writer.Write([]byte(msg))
	return err
}

// ================================================================ //
// ========== Redirecciones (3xx) ================================= //

// Redirect the request to a provided URL with status code.
func (c *Context) Redirect(code int, url string) error {
	if code < 300 || code > 308 {
		return gko.ErrInesperado().Str("redirect inválido").Ctx("code", code)
	}
	c.response.Header().Set(HeaderLocation, url)
	c.response.WriteHeader(code)
	return nil
}

// Redirige a la URL usando fmt.Sprintf para construir el path.
//
// Normal: 303 StatusSeeOther & header Location.
//
// HTMX: 200 OK "Redirigiendo a..." & header HX-Redirect.
func (c *Context) Redir(url string) error {
	if c.EsHTMX() {
		c.response.Header().Set("HX-Redirect", url)
		return c.StatusOk("Redirigiendo a " + url)
	}
	c.response.Header().Set(HeaderLocation, url)
	c.response.WriteHeader(303)
	return nil
}

// Redirige a la URL usando fmt.Sprintf para construir el path.
//
// Normal: 303 StatusSeeOther & header Location.
//
// HTMX: 200 OK "Redirigiendo a..." & header HX-Redirect.
func (c *Context) Redirf(format string, a ...any) error {
	if c.EsHTMX() {
		c.response.Header().Set("HX-Redirect", fmt.Sprintf(format, a...))
		return c.StatusOk("Redirigiendo a " + fmt.Sprintf(format, a...))
	}
	c.response.Header().Set(HeaderLocation, fmt.Sprintf(format, a...))
	c.response.WriteHeader(303)
	return nil
}
