package gecko

import (
	"io"
)

// Agregar un header a la respuesta.
func (c *Context) HeaderAdd(key string, value string) {
	c.response.Header().Add(key, value)
}

// Responder con status code y MIME "text/html" UTF8.
func (c *Context) HTML(code int, html string) (err error) {
	return c.Blob(code, MIMETextHTMLCharsetUTF8, []byte(html))
}

// Responder con status code y MIME "text/html" UTF8.
func (c *Context) HTMLBlob(code int, b []byte) (err error) {
	return c.Blob(code, MIMETextHTMLCharsetUTF8, b)
}

// Responder con status code y MIME "text/plain" UTF8.
//
// Preferible usar c.StringOk("msg") o retornar un error
// para que sí quede registrado en el log por el ErrHandler.
func (c *Context) String(code int, s string) (err error) {
	return c.Blob(code, MIMETextPlainCharsetUTF8, []byte(s))
}

func (c *Context) ContentOk(contentType string, content []byte) (err error) {
	c.writeContentType(contentType)
	c.response.WriteHeader(200)
	_, err = c.response.Write(content)
	return
}

// Responder con status code y MIME especificados. Ver gecko.MIME...
func (c *Context) Blob(code int, contentType string, b []byte) (err error) {
	c.writeContentType(contentType)
	c.response.WriteHeader(code)
	_, err = c.response.Write(b)
	return
}

// Responder con status code y MIME especificados. Ver gecko.MIME...
func (c *Context) Stream(code int, contentType string, r io.Reader) (err error) {
	c.writeContentType(contentType)
	c.response.WriteHeader(code)
	_, err = io.Copy(c.response, r)
	return
}

// Responder con un body vacío y un status code.
// Útil para enviar un Status304NotModified.
func (c *Context) NoContent(code int) error {
	c.response.WriteHeader(code)
	return nil
}

// ================================================================ //

func (c *Context) writeContentType(value string) {
	header := c.response.Header()
	if header.Get(HeaderContentType) == "" {
		header.Set(HeaderContentType, value)
	}
}
