package gecko

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pargomx/gecko/gko"
)

// El handler centralizado para enviar al cliente errores de gecko.
//
// NOTA: el error se ignora cuando se genera en un middleware después
// de que un handler ya respondió algo al cliente sin error.
func (g *Gecko) responderErrorHTTP(c *Context, err error) {
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

	// Copiar error y agregar contexto http para loggearlo sin agregar info redundante a log http.
	gkerr := *gko.Err(err)
	gkerr.Op(c.path) // Patrón de ruta registrada para ubicar handler.
	if strings.Contains(c.path, "}") {
		gkerr.Ctx("path", c.request.URL.Path) // Si hay parámetros en la ruta se incluyen.
		// c.request.URL.Path     // Ruta sin query.
		// c.request.URL.String() // Ruta con query.
	}
	if len(c.SesionID) > 6 {
		gkerr.Ctx("sesion", c.SesionID[:6]) // Conocer usuario sin exponer sesión.
	}
	gkerr.Log()

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
	if g.Renderer != nil && g.TmplError != "" {
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
		`<html><head><title>Error %d</title></head><body style="background-color:black;color:white;"><h2>%s</h2><a href="/" style="color:aqua;">Ir a inicio</a></body></html>`+"\n",
		gkerr.GetCodigoHTTP(), gkerr.GetMensaje(),
	))
	if err != nil {
		gko.LogAlert("gko.ErrHandler: default response: " + err.Error())
	}
}
