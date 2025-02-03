package gecko

import (
	"net/http"
	"strings"
	"time"

	"github.com/pargomx/gecko/gko"
)

// Tratar rutas con trailing slash como si no lo tuvieran.
// Utilizado como middleware global antes del router.
func quitarTrailingSlash(r *http.Request) {
	url := r.URL
	path := url.Path
	queryString := r.URL.RawQuery
	l := len(path) - 1
	if l > 0 && strings.HasSuffix(path, "/") {
		path = path[:l]
		uri := path
		if queryString != "" {
			uri += "?" + queryString
		}
		r.RequestURI = uri
		url.Path = path
	}
}

// ================================================================ //
// ========== RUTAS Y HANDLERS ==================================== //

// HandlerFunc defines a function to serve HTTP requests.
type HandlerFunc func(c *Context) error

// Registrar una nueva ruta con un http.HandlerFunc
// que prepare el gecko.Context y ejecute el gecko.HandlerFunc.
func (g *Gecko) registrarRuta(método string, ruta string, handler HandlerFunc) {
	patrón := toMuxPattern(método, ruta)
	g.mux.HandleFunc(patrón, func(w http.ResponseWriter, r *http.Request) {
		c := &Context{
			request:  r,
			response: NewResponse(w, g),
			path:     patrón,
			gecko:    g,
			time:     time.Now(),
		}
		err := handler(c)
		if err != nil {
			g.responderErrorHTTP(c, err)
		}
		if g.HTTPLogger != nil {
			g.logHTTP(c, err)
		}
	})
	// fmt.Println("RUTA:", patrón)
}

// NotFound handler para rutas GET no registradas evitando usar el de *http.ServeMux.
func (g *Gecko) registrarNotFoundHandler() {
	g.mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		c := &Context{
			request:  r,
			response: NewResponse(w, g),
			path:     "GET /{...}",
			gecko:    g,
			time:     time.Now(),
		}
		err := gko.ErrNoEncontrado()
		g.responderErrorHTTP(c, err)
		if g.HTTPLogger != nil {
			g.logHTTP(c, err)
		}
	})
}

// Necesario para validar patrón de ruta con método y las reglas de gecko.
func toMuxPattern(método string, ruta string) string {
	// Validar la ruta.
	if strings.Contains(ruta, " ") {
		gko.FatalExitf("gecko.Router: ruta no puede contener espacios en blanco: '%s'", ruta)
	}
	ruta = strings.TrimSuffix(ruta, "/")
	if ruta == "" {
		ruta = "/{$}"
	}
	if ruta[0] != '/' {
		gko.FatalExitf("gecko.Router: ruta debe comenzar con slash: '%s'", ruta)
	}
	// Validar método.
	if método == "" {
		gko.FatalExitf("gecko.Router: método no puede estar indefinido: '%s'", ruta)
	}
	return método + " " + ruta
}

// ================================================================ //
// ========== Registrar handlers con métodos ====================== //

func (g *Gecko) GET(path string, handler HandlerFunc) {
	g.registrarRuta(http.MethodGet, path, handler)
}
func (g *Gecko) POST(path string, handler HandlerFunc) {
	g.registrarRuta(http.MethodPost, path, handler)
}
func (g *Gecko) PUT(path string, handler HandlerFunc) {
	g.registrarRuta(http.MethodPut, path, handler)
}
func (g *Gecko) PATCH(path string, handler HandlerFunc) {
	g.registrarRuta(http.MethodPatch, path, handler)
}
func (g *Gecko) DELETE(path string, handler HandlerFunc) {
	g.registrarRuta(http.MethodDelete, path, handler)
}

func (g *Gecko) POS(path string, handler HandlerFunc) {
	g.registrarRuta(http.MethodPost, path, handler)
}
func (g *Gecko) PCH(path string, handler HandlerFunc) {
	g.registrarRuta(http.MethodPatch, path, handler)
}
func (g *Gecko) DEL(path string, handler HandlerFunc) {
	g.registrarRuta(http.MethodDelete, path, handler)
}

// Registra un handler que redirige con StatusSeeOther (303) a la URL dada.
func (g *Gecko) Redir(path string, redirURL string) {
	g.registrarRuta(http.MethodGet, path, func(c *Context) error {
		return c.RedirOtro(redirURL)
	})
}

/*
func (g *Gecko) OPTIONS(path string, handler HandlerFunc) {
	g.registrarRuta(http.MethodOptions, path, handler)
}
func (g *Gecko) HEAD(path string, handler HandlerFunc) {
	g.registrarRuta(http.MethodHead, path, handler)
}
func (g *Gecko) CONNECT(path string, handler HandlerFunc) {
	g.registrarRuta(http.MethodConnect, path, handler)
}
func (g *Gecko) TRACE(path string, handler HandlerFunc) {
	g.registrarRuta(http.MethodTrace, path, handler)
}
*/

// ================================================================ //
