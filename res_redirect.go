package gecko

import (
	"fmt"

	"github.com/pargomx/gecko/gko"
)

// Devuelve "204 No Content" con "HX-Refresh" para que HTMX
// vuelva a cargar la página actual por completo.
func (c *Context) RefreshHTMX() error {
	c.response.Header().Set("HX-Refresh", "true")
	return c.NoContent(204)
}

// ================================================================ //

// Redirigir al cliente con un código 300-308 y header "Location".
func (c *Context) RedirCode(code int, url string) error {
	if code < 300 || code > 308 {
		return gko.ErrInesperado.Str("redirect inválido").Ctx("code", code)
	}
	c.response.Header().Set(HeaderLocation, url)
	c.response.WriteHeader(code)
	return nil
}

// Redirigir al cliente con un código 300-308 y header "Location".
// Utilizar fmt.Sprintf para construir el url.
func (c *Context) RedirCodef(code int, format string, a ...any) error {
	if code < 300 || code > 308 {
		return gko.ErrInesperado.Str("redirect inválido").Ctx("code", code)
	}
	c.response.Header().Set(HeaderLocation, fmt.Sprintf(format, a...))
	c.response.WriteHeader(code)
	return nil
}

// ================================================================ //

// Redirigir al cliente con "303 SeeOther" y header "Location".
func (c *Context) RedirOtro(url string) error {
	c.response.Header().Set(HeaderLocation, url)
	c.response.WriteHeader(Status303SeeOther)
	return nil
}

// Redirigir al cliente con "303 SeeOther" y header "Location".
// Utilizar fmt.Sprintf para construir el url.
func (c *Context) RedirOtrof(format string, a ...any) error {
	c.response.Header().Set(HeaderLocation, fmt.Sprintf(format, a...))
	c.response.WriteHeader(Status303SeeOther)
	return nil
}

// ================================================================ //

// Redirigir al cliente a cargar por completo una nueva página.
// Utiliza "200 OK" con "HX-Redirect" si es solicitud HTMX.
// Utiliza "303 SeeOther" con "Location" si no es HTMX.
func (c *Context) RedirFull(url string) error {
	if c.EsHTMX() {
		c.response.Header().Set(HxRedirect, url)
		return c.StatusOk("Redirigiendo a " + url)
	}
	c.response.Header().Set(HeaderLocation, url)
	c.response.WriteHeader(Status303SeeOther)
	return nil
}

func (c *Context) RedirFullf(format string, a ...any) error {
	return c.RedirFull(fmt.Sprintf(format, a...))
}

// ================================================================ //

// func (c *Context) Redirect(code int, url string) error {
// 	if code < 300 || code > 308 {
// 		return gko.ErrInesperado().Str("redirect inválido").Ctx("code", code)
// 	}
// 	c.response.Header().Set(HeaderLocation, url)
// 	c.response.WriteHeader(code)
// 	return nil
// }

// Redirige a la URL usando fmt.Sprintf para construir el path.
// Normal: 303 StatusSeeOther & header Location.
// HTMX: 200 OK "Redirigiendo a..." & header HX-Redirect.
// func (c *Context) Redir(url string) error {
// 	if c.EsHTMX() && !(c.request.Method == "GET" || c.request.Method == "HEAD") {
// 		c.response.Header().Set("HX-Redirect", url)
// 		return c.StatusOk("Redirigiendo a " + url)
// 	}
// 	c.response.Header().Set(HeaderLocation, url)
// 	c.response.WriteHeader(303)
// 	return nil
// }

// Redirige a la URL usando fmt.Sprintf para construir el path.
//
// Normal: 303 StatusSeeOther & header Location.
//
// HTMX: 200 OK "Redirigiendo a..." & header HX-Redirect.
// func (c *Context) Redirf(format string, a ...any) error {
// 	if c.EsHTMX() && !(c.request.Method == "GET" || c.request.Method == "HEAD") {
// 		c.response.Header().Set("HX-Redirect", fmt.Sprintf(format, a...))
// 		return c.StatusOk("Redirigiendo a " + fmt.Sprintf(format, a...))
// 	}
// 	c.response.Header().Set(HeaderLocation, fmt.Sprintf(format, a...))
// 	c.response.WriteHeader(303)
// 	return nil
// }
