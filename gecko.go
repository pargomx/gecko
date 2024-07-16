package gecko

import (
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	mux              *http.ServeMux
	IPExtractor      IPExtractor
	Renderer         Renderer
	HTTPErrorHandler func(err error, c *Context)

	// Filesystem is file system used by Static and File handlers to access
	// files. Defaults to os.DirFS(".")
	//
	// When dealing with `embed.FS` use `fs := gecko.MustSubFS(fs,
	// "rootDirectory") to create sub fs which uses necessary prefix for
	// directory path. This is necessary as `//go:embed assets/images` embeds
	// files with paths including `assets/images` as their prefix.
	Filesystem fs.FS

	TmplBaseLayout string
}

// Implementa la interfaz http.Handler.
func (g *Gecko) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Middleware global.
	quitarTrailingSlash(r)
	// fmt.Println("Sirviendo", r.Method, r.URL.Path)
	// Proceder con el router de la librería estándar.
	g.mux.ServeHTTP(w, r)
}

// Nuevo servidor escuchando en :8080.
func New() *Gecko {
	pwd, err := os.Getwd()
	if err != nil {
		FatalErr(err)
	}
	return &Gecko{
		mux: http.NewServeMux(),

		Filesystem: os.DirFS(pwd),

		HTTPErrorHandler: errorHandler,

		IPExtractor: ExtractIPFromRealIPHeader(),

		TmplBaseLayout: "base_layout",
	}
}

// Inicia el servidor HTTP.
func (g *Gecko) IniciarEnPuerto(port int) error {
	if port < 1 || port > 65535 {
		return ErrNotAcceptable.Msg("puerto TCP inválido")
	}
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: g,
	}
	LogEventof("Escuchando en tcp/%d", port)
	return srv.ListenAndServe()
}

// Create a Unix domain sock and listen for incoming connections.
func (g *Gecko) IniciarEnSocket(socket string) error {
	if socket == "" {
		return ErrNotAcceptable.Msg("socket path indefinido")
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
			LogError(err)
		} else {
			LogInfof("socket removido %v", socket)
		}
		os.Exit(0)
	}()
	LogEventof("Escuchando en unix %v", socket)
	return http.Serve(sock, g)
}
