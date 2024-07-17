package gecko

import (
	"net/http"
	"net/url"
)

// Representa a la solicitud HTTP actual
// y ofrece los medios para responderla.
type Context struct {
	request  *http.Request
	response *Response
	path     string
	query    url.Values
	gecko    *Gecko
	SesionID string
}

func (c *Context) Request() *http.Request {
	return c.request
}

func (c *Context) Response() *Response {
	return c.response
}
