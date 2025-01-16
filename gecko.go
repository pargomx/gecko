package gecko

import (
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pargomx/gecko/gko"
)

// Gecko es un servidor web simple basado en la librería estándar de Go 1.22.
//
// El router tiene algunas particularidades:
//
//   - Las rutas siempre son para un método específico.
//   - Las rutas no pueden contener espacios en blanco.
//   - Las rutas deben comenzar con slash.
//   - Las rutas con trailing slash son tratadas como si no lo tuvieran.
//
// Ejemplos:
//   - Solicitudes "/hola" y "/hola/" usarán el mismo handler.
//   - Solicitud "/hola/x/y/z" no usará el handler de "/hola/".
type Gecko struct {
	mux         *http.ServeMux
	IPExtractor IPExtractor
	Renderer    Renderer
	HTTPLogger  HTTPLogger

	Filesystem fs.FS // Utilizado por los file handlers.

	TmplBaseLayout string // Nombre de la plantilla base.
	TmplError      string // Nombre de la plantilla para errores.
}

// Implementa la interfaz http.Handler.
func (g *Gecko) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Middleware global.
	quitarTrailingSlash(r)
	// fmt.Println("Sirviendo", r.Method, r.URL.Path)
	// Proceder con el router de la librería estándar.
	g.mux.ServeHTTP(w, r)
}

// Nuevo servidor gecko.
func New() *Gecko {
	pwd, err := os.Getwd()
	if err != nil {
		gko.FatalError(err)
	}
	g := &Gecko{
		mux: http.NewServeMux(),

		Filesystem: os.DirFS(pwd),

		TmplBaseLayout: "base_layout",
		TmplError:      "",
	}
	// NotFound handler consistente para rutas GET no registradas.
	g.mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		g.ResponderHTTPHandlerError(gko.ErrNoEncontrado(), &Context{
			request:  r,
			response: NewResponse(w, g),
			path:     "GET /{...}",
			gecko:    g,
		})
	})
	return g
}

// Nuevo router cuando se registra una ruta en "GET /".
func NewSinRoot404() *Gecko {
	pwd, err := os.Getwd()
	if err != nil {
		gko.FatalError(err)
	}
	g := &Gecko{
		mux: http.NewServeMux(),

		Filesystem: os.DirFS(pwd),

		TmplBaseLayout: "base_layout",
		TmplError:      "",
	}
	return g
}

// Inicia el servidor HTTP.
func (g *Gecko) IniciarEnPuerto(port int) error {
	if port < 1 || port > 65535 {
		return gko.ErrDatoInvalido().Msg("puerto TCP inválido")
	}
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: g,
	}
	gko.LogEventof("Escuchando en tcp/%d", port)
	return srv.ListenAndServe()
}

// Create a Unix domain sock and listen for incoming connections.
func (g *Gecko) IniciarEnSocket(socket string) error {
	if socket == "" {
		return gko.ErrDatoIndef().Msg("socket path indefinido")
	}
	sock, err := net.Listen("unix", socket)
	if err != nil {
		return err
	}
	// Cleanup the socket file.
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-exitChan
		err = os.Remove(socket)
		if err != nil {
			gko.LogError(err)
		} else {
			gko.LogInfof("socket removido %v", socket)
		}
		os.Exit(0)
	}()
	gko.LogEventof("Escuchando en unix %v", socket)
	return http.Serve(sock, g)
}
