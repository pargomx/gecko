package gecko

import (
	"strings"

	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/gkt"
)

// ================================================================ //
// ========== Request HTMX ======================================== //

// Si la solicitud viene de HTMX significa que tiene el header HX-Request = true.
// Cuando es HX-History-Restore-Request se necesita enviar la página entera.
func (c *Context) EsHTMX() bool {
	return c.request.Header.Get("HX-Request") == "true" &&
		c.request.Header.Get("HX-History-Restore-Request") != "true"
}

// Agrega un evento al HX-Trigger
func (c *Context) TriggerEventoHTMX(evento string) {
	c.response.Header().Set("HX-Trigger", evento)
}

// Responder con lo especificado en el cliente mediante el
// atributo hg-askfor. Se pueden solicitar tres cosas:
//
// Redirección simple al recurso:
//  'hx-put="/some" hg-askfor="/recurso" ...'
//
// Recarga de la página con esta url.
//  'hx-put="/some" hg-askfor="full:/recurso" '
//
// Lanzar este evento con htmx.
//  'hx-put="/some" hg-askfor="event:somethingHappend" '
//
// Si la solicitud no tiene el header Hg-Askfor se responde con fallback con c.RedirOtrof().
//
// Para que funcione el cliente debe tener htmx y este eventListener:
/*
	document.addEventListener("htmx:configRequest", function (event) {
		if (event.target.hasAttribute("hg-askfor")) {
			const askforVal = event.target.getAttribute("hg-askfor");
			if (askforVal && askforVal.length > 0) {
				event.detail.headers["Hg-Askfor"] = askforVal;
			}
		}
	});
*/
func (c *Context) AskedForFallback(fallbackRedir string, a ...any) error {
	askfor := c.Request().Header.Get("Hg-Askfor")
	askfor = gkt.SinEspaciosNinguno(askfor)

	// Fallback si no se pidió algo específico.
	if askfor == "" {
		return c.RedirOtrof(fallbackRedir, a...)
	}

	// Mayoritariamente un recurso redirect.
	if strings.HasPrefix(askfor, "/") {
		return c.RedirOtro(askfor)
	}

	// A veces se pide un evento.
	evento, askEvent := strings.CutPrefix(askfor, "event:")
	if askEvent {
		c.TriggerEventoHTMX(evento)
		return c.StringOk(evento)
	}

	// O quizá un full page reload. Solo permitir redirecciones al mismo sitio.
	urlFullRedir, askFullRedir := strings.CutPrefix(askfor, "full:/")
	if askFullRedir {
		return c.RedirFull("/" + urlFullRedir)
	}

	return gko.ErrDatoInvalido().Strf("askfor invalid: %v", askfor)
}

// Responder con lo especificado en el cliente mediante el
// atributo hg-askfor. Se pueden solicitar tres cosas:
//
// hg-askfor="/recurso" Redirección simple al recurso.
//
// hg-askfor="full:/recurso" Recarga de la página con esta url.
//
// hg-askfor="event:somethingHappend" Lanzar este evento con htmx.
//
// Si la solicitud no tiene el header Hg-Askfor se responde con c.StringOk("Ok")
func (c *Context) AskedFor() error {
	askfor := c.Request().Header.Get("Hg-Askfor")
	askfor = gkt.SinEspaciosNinguno(askfor)

	// Fallback si no se pidió algo específico.
	if askfor == "" {
		return c.StringOk("Ok")
	}

	// Mayoritariamente un recurso redirect.
	if strings.HasPrefix(askfor, "/") {
		return c.RedirOtro(askfor)
	}

	// A veces se pide un evento.
	evento, askEvent := strings.CutPrefix(askfor, "event:")
	if askEvent {
		c.TriggerEventoHTMX(evento)
		return c.StringOk(evento)
	}

	// O quizá un full page reload. Solo permitir redirecciones al mismo sitio.
	urlFullRedir, askFullRedir := strings.CutPrefix(askfor, "full:/")
	if askFullRedir {
		return c.RedirFull("/" + urlFullRedir)
	}

	return gko.ErrDatoInvalido().Strf("askfor invalid: %v", askfor)
}
