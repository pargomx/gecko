package gecko

import (
	"fmt"
	"net/http"
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
