package gecko

// ================================================================ //
// ========== Request HTMX ======================================== //

// Si la solicitud viene de HTMX significa que tiene el header HX-Request = true.
// Cuando es HX-History-Restore-Request se necesita enviar la página entera.
func (c *Context) EsHTMX() bool {
	return c.request.Header.Get("HX-Request") == "true" &&
		c.request.Header.Get("HX-History-Restore-Request") != "true"
}

// ================================================================ //
// ========== Responder HTMX ====================================== //

// Devuelve un estatus "204 No Content" e instruye a HTMX para que
// vuelva a cargar la página entera con el header "HX-Refresh".
//
// Conveniente como respuesta a una solicitud PUT.
func (c *Context) RefreshHTMX() error {
	c.response.Header().Set("HX-Refresh", "true")
	return c.NoContent(204)
}

// Agrega un evento al HX-Trigger
func (c *Context) TriggerEventoHTMX(evento string) {
	c.response.Header().Set("HX-Trigger", evento)
}
