package gecko

import (
	"fmt"
	"net/http"

	"github.com/pargomx/gecko/gko"
)

// El handler centralizado para enviar al cliente errores de gecko.
//
// NOTA: el error se ignora cuando se genera en un middleware después
// de que un handler ya respondió algo al cliente sin error.
func (g *Gecko) ResponderHTTPHandlerError(err error, c *Context) {
	if err == nil {
		gko.LogAlertf("gko.ErrHandler: err nil: %s", c.path)
		return
	}
	if c == nil {
		gko.LogAlertf("gko.ErrHandler: context nil: %v", err)
		return
	}
	if c.response.Committed {
		gko.LogAlertf("gko.ErrHandler: err returned after response: %s %s", c.path, err)
		return
	}

	// Agregar contexto al error y loggearlo.
	gkerr := gko.Err(err)
	gkerr.Op(c.path) // Patrón de ruta registrada para ubicar handler.
	if len(c.SesionID) > 6 {
		gkerr.Ctx("sesion", c.SesionID[:6]) // Saber usuario sin exponer sesión.
	}
	gkerr.Log()

	// gkerr.Op(c.request.Method + " " + c.request.URL.Path) // Ruta sin query.
	// gkerr.Op(c.request.Method + " " + c.request.URL.String()) // Ruta con query.

	// Método HEAD debe responder sin body.
	if c.request.Method == http.MethodHead {
		err := c.NoContent(gkerr.GetCodigoHTTP())
		if err != nil {
			gko.LogAlert("gko.ErrHandler: head response: " + err.Error())
		}
		return
	}

	// HTMX solo necesita un string.
	if c.EsHTMX() {
		err = c.String(gkerr.GetCodigoHTTP(), gkerr.GetMensaje())
		if err != nil {
			gko.LogAlert("gko.ErrHandler: htmx response: " + err.Error())
		}
		return
	}

	// Mandar plantilla con el error.
	if g.Renderer != nil {
		data := map[string]any{
			"Mensaje":    gkerr.GetMensaje(),
			"StatusCode": gkerr.GetCodigoHTTP(),
			"Titulo":     "Ups: " + gkerr.GetMensaje(),
		}
		if c.Sesion != nil {
			data["Sesion"] = c.Sesion
		}
		err = c.Render(gkerr.GetCodigoHTTP(), g.TmplError, data)
		if err != nil {
			gko.LogAlert("gko.ErrHandler: render err: " + err.Error())
		}
		return
	}

	// Default: responder con texto.
	err = c.HTML(gkerr.GetCodigoHTTP(), fmt.Sprintf(
		`<html><head><title>Error %d</title></head><body style="background-color:black;color:white;"><h2>%s</h2><a href="/" style="color:aqua;">Ir a inicio</a></body></html>`,
		gkerr.GetCodigoHTTP(), gkerr.GetMensaje(),
	))
	if err != nil {
		gko.LogAlert("gko.ErrHandler: default response: " + err.Error())
	}
}
