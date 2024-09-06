package gecko

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pargomx/gecko/gko"
)

// Envía un archivo como respuesta o el index.html si es un directorio.
func fsFile(c *Context, fpath string, filesystem fs.FS) error {
	// Si se pide el dir raíz fpath vendrá vacío y es inválido para fs.Stat
	if fpath == "" {
		fpath = "."
	}
	// Verificar si el archivo existe.
	fi, err := fs.Stat(filesystem, fpath)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			gko.Err(err).Op("fsFile.Stat('" + fpath + "')").Log()
		}
		return gko.ErrNoEncontrado()
	}
	// Si es un directorio se sirve el index.html
	if fi.IsDir() {
		return fsDirIndex(c, fpath, fi, filesystem)
	}
	// Abrir el archivo.
	file, err := filesystem.Open(fpath)
	if err != nil {
		gko.Err(err).Op("fsFile.Open('" + fpath + "')").Log()
		return gko.ErrNoEncontrado()
	}
	defer file.Close()
	// Enviar el archivo.
	ff, ok := file.(io.ReadSeeker)
	if !ok {
		return gko.ErrInesperado().Str("file does not implement io.ReadSeeker")
	}
	http.ServeContent(c.response, c.request, fi.Name(), fi.ModTime(), ff)
	return nil
}

// Se servirá el index.html si existe en el directorio fpath.
func fsDirIndex(c *Context, fpath string, fi fs.FileInfo, filesystem fs.FS) error {
	if !fi.IsDir() {
		return gko.ErrInesperado().Str("fpath is not a directory")
	}
	fpath = filepath.ToSlash(filepath.Join(fpath, "index.html"))
	// Abrir el archivo.
	file, err := filesystem.Open(fpath)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			gko.Err(err).Op("fsDirIndex.Open('" + fpath + "')").Log()
		}
		return gko.ErrNoEncontrado()
	}
	defer file.Close()
	// Enviar el archivo.
	ff, ok := file.(io.ReadSeeker)
	if !ok {
		return gko.ErrInesperado().Str("file does not implement io.ReadSeeker")
	}
	http.ServeContent(c.response, c.request, fi.Name(), fi.ModTime(), ff)
	return nil
}

// staticDirectoryHandler creates handler function to serve files from provided
// file system When disablePathUnescaping is set then file name from path is not
// unescaped and is served as is.
func staticDirectoryHandler(filesystem fs.FS) HandlerFunc {
	return func(c *Context) error {
		fpath := c.Param("fpath")
		// Convertir %2F a / por ejemplo.
		fpath, err := url.PathUnescape(fpath)
		if err != nil {
			return fmt.Errorf("failed to unescape path variable: %w", err)
		}
		// Necesario en windows porque fs.FS solo usa slashes.
		fpath = filepath.ToSlash(fpath)
		// fs.Open() asume que fpath es relativa al root y rechaza el prefijo `/`.
		filepath.Clean(strings.TrimPrefix(fpath, "/"))
		// Servir el archivo solicitado.
		return fsFile(c, fpath, filesystem)
	}
}

// ================================================================ //
// ========== SERVIR ARCHIVOS ===================================== //

// Crea un handler para servir un archivo estático desde un filesystem dado.
func StaticFileHandler(file string, filesystem fs.FS) HandlerFunc {
	return func(c *Context) error {
		return fsFile(c, file, filesystem)
	}
}

// Registra una ruta para servir un archivo desde el filesystem dado.
func (g *Gecko) FileFS(path, file string, filesystem fs.FS) {
	g.registrarRuta(http.MethodGet, path, func(c *Context) error {
		return fsFile(c, file, filesystem)
	})
}

// Registra una nueva ruta para servir un archivo.
// Debe ser una ruta relativa por lo regular.
func (g *Gecko) File(path, file string) {
	g.registrarRuta(http.MethodGet, path, func(c *Context) error {
		return fsFile(c, file, c.gecko.Filesystem)
	})
}

// Registra una nueva ruta para servir un archivo con path absoluta.
// Jamás poner a disposición del usuario!
func (g *Gecko) FileAbs(path, fpath string) {
	if !filepath.IsAbs(fpath) {
		gko.FatalExitf("ruta absoluta inválida: %s", fpath)
	}
	fi, err := os.Stat(fpath)
	if err != nil {
		gko.Err(err).Op("FileAbs.Stat('" + fpath + "')").FatalExit()
	}
	if fi.IsDir() {
		gko.FatalExitf("FileAbs('%s') es un directorio", fpath)
	}
	g.registrarRuta(http.MethodGet, path, func(c *Context) error {
		fi, err = os.Stat(fpath)
		if err != nil {
			gko.Err(err).Op("FileAbs.Stat('" + fpath + "')").FatalExit()
		}
		file, err := os.Open(fpath)
		if err != nil {
			gko.Err(err).Op("FileAbs.Open('" + fpath + "')").Log()
			return gko.ErrNoEncontrado()
		}
		defer file.Close()
		http.ServeContent(c.response, c.request, fi.Name(), fi.ModTime(), file)
		return nil
	})
}

// Envía el contenido de un archivo como respuesta desde un filesystem dado.
// Deduce el ContentType y se encarga del caché gracias a http.ServeContent.
func (c *Context) FileFS(file string, filesystem fs.FS) error {
	return fsFile(c, file, filesystem)
}

// Envía el contenido de un archivo como respuesta.
// Deduce el ContentType y se encarga del caché gracias a http.ServeContent.
func (c *Context) File(file string) error {
	return fsFile(c, file, c.gecko.Filesystem)
}

// FileAttachment es similar a File() excepto que se usa para enviar
// un archivo como adjunto, especificando un nombre para él.
//
// Hace que el navegador descargue el archivo sin visualizarlo.
//
// Content-Disposition = "attachment; filename=<FILE_NAME>"
func (c *Context) FileAttachment(file, name string) error {
	c.response.Header().Set(HeaderContentDisposition, fmt.Sprintf("attachment; filename=%q", name))
	return fsFile(c, file, c.gecko.Filesystem)
}

// FileInline es similar a File() excepto que se usa para enviar
// un archivo como inline, especificando un nombre para él.
//
// Hace que el navegador visualice el archivo sin descargarlo.
//
// Content-Disposition = "inline; filename=<FILE_NAME>"
func (c *Context) FileInline(file, name string) error {
	c.response.Header().Set(HeaderContentDisposition, fmt.Sprintf("inline; filename=%q", name))
	return fsFile(c, file, c.gecko.Filesystem)
}

// ================================================================ //
// ========== SERVIR DIRECTORIOS ================================== //

// Crea un filesystem en donde el `subdir` es la nueva raíz.
// La ruta debe ser un subdirectorio relativo al filesystem dado.
//
// No se permiten rutas absolutas ni fuera de la raíz del fs dado,
// pero no se verifica que los symlink no escapen del nuevo root.
func mustSubFS(currentFs fs.FS, subdir string) fs.FS {
	subdir = filepath.ToSlash(filepath.Clean(subdir))
	if subdir == "." || subdir == "./" || subdir == "" {
		gko.FatalExitf("subdir inválido: utilice Static(...) para servir '%s'", subdir)
	}
	if strings.HasPrefix(subdir, "/") {
		gko.FatalExitf("subdir inválido: '%s' no es relativo", subdir)
	}
	if strings.HasPrefix(subdir, "..") {
		gko.FatalExitf("subdir inválido: no se permite '../' en %s", subdir)
	}
	rootInfo, err := fs.Stat(currentFs, subdir)
	if err != nil {
		gko.FatalExitf("subdir inválido: %v", err)
	}
	if !rootInfo.IsDir() {
		gko.FatalExitf("subdir inválido: %s no es un directorio", subdir)
	}
	subFs, err := fs.Sub(currentFs, subdir)
	if err != nil {
		gko.FatalExitf("imposible crear subFS: subdir inválido: %v", err)
	}
	return subFs
}

// Registra una ruta para servir archivos estáticos desde un directorio
// absoluto desde el root filesystem del sistema operativo.
// Usar con precaución y usar `StaticSub` o `StaticPwd` si es posible.
//
// Por ejemplo:
//
//	g.StaticAbs("/x", "/home/user/static")
//
// Servirá "/home/user/static/hola.html" en la ruta "/x/hola.html"
func (g *Gecko) StaticAbs(rutaWeb string, absolutePath string) {
	if !filepath.IsAbs(absolutePath) {
		gko.FatalExitf("ruta absoluta inválida: %s", absolutePath)
	}
	rootInfo, err := os.Stat(absolutePath)
	if err != nil {
		gko.FatalExitf("ruta absoluta inválida: %v", err)
	}
	if !rootInfo.IsDir() {
		gko.FatalExitf("ruta absoluta inválida: %s no es un directorio", absolutePath)
	}
	handler := staticDirectoryHandler(os.DirFS(absolutePath))
	g.registrarRuta(http.MethodGet, rutaWeb, handler)
	g.registrarRuta(http.MethodGet, path.Join(rutaWeb, "{fpath...}"), handler)
}

// Se registra tanto /files como /files/*. para que el mux no redireccione
// /files a /files/ en un loop infinito debido a quitarTrailingSlash(r).

// Registra una ruta para servir archivos en el directorio actual
// desde la ruta dada.
//
// Por ejemplo, si el directorio actual es "/home/user" y se llama
//
//	g.StaticDir("/static")
//
// Servirá "/home/user/hola.html" en la ruta "/static/hola.html"
//
// Cuando la ruta solicitada es un directorio se intentará servir
// el archivo index.html en ese directorio. Si no existe retorna 404.
func (g *Gecko) StaticPwd(rutaWeb string) {
	handler := staticDirectoryHandler(g.Filesystem)
	g.registrarRuta(http.MethodGet, rutaWeb, handler)
	g.registrarRuta(http.MethodGet, path.Join(rutaWeb, "{fpath...}"), handler)
}

// Registra una ruta para servir archivos en el subdirectorio dado.
//
// La ruta debe ser relativa al directorio actual, por ejemplo:
//
//	g.StaticSub("/x", "static")
//	g.StaticSub("/y", "./static/img")
//
//	La ruta "/x" servirá el archivo "static/index.html"
//	La ruta "/x/img/1.png" servirá el archivo "static/img/1.png"
//	La ruta "/y/1.png" servirá el archivo "static/img/1.png"
//
// Una ruta absoluta es inválida, así como la ruta de un archivo.
// Para servir el directorio actual usar `g.Static("/")`.
func (g *Gecko) StaticSub(rutaWeb string, subdir string) {
	handler := staticDirectoryHandler(mustSubFS(g.Filesystem, subdir))
	g.registrarRuta(http.MethodGet, rutaWeb, handler)
	g.registrarRuta(http.MethodGet, path.Join(rutaWeb, "{fpath...}"), handler)
}

// Registra una ruta para servir archivos estáticos desde un filesystem dado.
// Es útil para servir archivos embebidos, por ejemplo:
//
//	//go:embed static/img
//	var fs embed.FS
//	fs := gecko.StaticFS("/x", fs)
//
// La ruta "/x/static/img/1.png" servirá el archivo "static/img/1.png
func (g *Gecko) StaticFS(pathPrefix string, filesystem fs.FS) {
	handler := staticDirectoryHandler(filesystem)
	g.registrarRuta(http.MethodGet, pathPrefix, handler)
	g.registrarRuta(http.MethodGet, path.Join(pathPrefix, "{fpath...}"), handler)
}

// Registra una ruta para servir archivos estáticos desde un filesystem dado
// quitando el prefijo de la ruta.
// Es útil para servir archivos embebidos, por ejemplo:
//
//	//go:embed static/img
//	var fs embed.FS
//	fs := gecko.StaticSubFS("/x", "static/img", fs)
//
// La ruta "/x/1.png" servirá el archivo "static/img/1.png
//
// Es necesario poner el subdir fsRoot porque `//go:embed assets/images`
// embebe archivos con el path incluyendo `assets/images` como prefijo.
func (g *Gecko) StaticSubFS(pathPrefix string, fsRoot string, filesystem fs.FS) {
	handler := staticDirectoryHandler(mustSubFS(filesystem, fsRoot))
	g.registrarRuta(http.MethodGet, pathPrefix, handler)
	g.registrarRuta(http.MethodGet, path.Join(pathPrefix, "{fpath...}"), handler)
}

// ================================================================ //
