package gecko

import (
	"fmt"
	"net/http"

	"github.com/pargomx/gecko/gko"
)

// ================================================================ //
// ========== Respuestas satisfactorias (2xx) ===================== //

func (c *Context) StringOk(msg string) error {
	c.Response().Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Response().WriteHeader(200)
	c.Response().Writer.Write([]byte(msg))
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
	c.Response().Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Response().WriteHeader(202)
	_, err := c.Response().Writer.Write([]byte(msg))
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

// Redirige a la URL usando fmt.Sprintf con código 303 TemporaryRedirect.
func (c *Context) Redir(format string, a ...any) error {
	c.response.Header().Set(HeaderLocation, fmt.Sprintf(format, a...))
	c.response.WriteHeader(303)
	return nil
}

// ================================================================ //
// ================================================================ //

// Retorna un error 400 Bad Request.
func (c *Context) StatusBadRequest(msg string) error {
	if msg == "" {
		msg = "Solicitud no aceptada"
	}
	return c.String(http.StatusBadRequest, msg)
}

// Retorna un error 401 Unauthorized, para usuario no autenticado.
func (c *Context) StatusUnauthorized(msg string) error {
	if msg == "" {
		msg = "No autorizado"
	}
	return c.String(http.StatusUnauthorized, msg)
}

// Retorna un error 402 Payment Required.
func (c *Context) StatusPaymentRequired(msg string) error {
	if msg == "" {
		msg = "Pago requerido"
	}
	return c.String(http.StatusPaymentRequired, msg)
}

// Retorna un error 403 Forbidden, para privilegios insuficientes.
func (c *Context) StatusForbidden(msg string) error {
	if msg == "" {
		msg = "No permitido"
	}
	return c.String(http.StatusForbidden, msg)
}

// Retorna un error 404 Not Found.
func (c *Context) StatusNotFound(msg string) error {
	if msg == "" {
		msg = "Recurso no encontrado"
	}
	return c.String(http.StatusNotFound, msg)
}

// Retorna un error 409 Conflict, para already exists u otros conflictos.
func (c *Context) StatusConflict(msg string) error {
	if msg == "" {
		msg = "Conflicto con recurso existente"
	}
	return c.String(http.StatusConflict, msg)
}

// Retorna un error 415 Unsupported Media Type.
func (c *Context) StatusUnsupportedMedia(msg string) error {
	if msg == "" {
		msg = "Tipo de media no soportado"
	}
	return c.String(http.StatusUnsupportedMediaType, msg)
}

// Retorna un error 429.
func (c *Context) StatusTooManyRequests(msg string) error {
	if msg == "" {
		msg = "Demasiadas solicitudes"
	}
	return c.String(http.StatusTooManyRequests, msg)
}

// Retorna un error 500 Internal Server Error.
func (c *Context) StatusServerError(msg string) error {
	if msg == "" {
		msg = "Error en servidor"
	}
	return c.String(http.StatusInternalServerError, msg)
}

// Retorna un error 500.
func (c *Context) ServerError(err error) error {
	return c.String(http.StatusInternalServerError, err.Error())
}
