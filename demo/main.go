package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/logsqlite"
)

func makeLogEntry(host, req string) gecko.LogEntry {
	return gecko.LogEntry{
		Timestamp: time.Now(),
		Host:      host,
		Ruta:      "/fake/" + host + "/" + req,
		Status:    200,
		BytesIn:   45,
		BytesOut:  1024,
	}
}

func main() {

	var mensaje, socket string
	var puerto int
	flag.IntVar(&puerto, "p", 8080, "Puerto TCP en el que escuchará el servidor")
	flag.StringVar(&socket, "s", "", "Socket Unix en el que escuchará el servidor")
	flag.StringVar(&mensaje, "m", "Hola mundo", "Mensaje que retornará el servidor")
	flag.Parse()

	fmt.Println("Preparando servidor")
	g := gecko.New()

	// HTTP Logger
	logger0, err := logsqlite.NewLogger("demo/httplog.sql", time.Second*3)
	if err != nil {
		gko.FatalError(err)
	}
	g.HTTPLogger = logger0
	defer logger0.Close()

	logger1, err := logsqlite.NewLogger("demo/httplog.sql", time.Second*3)
	if err != nil {
		gko.FatalError(err)
	}
	logger2, err := logsqlite.NewLogger("demo/httplog.sql", time.Second*3)
	if err != nil {
		gko.FatalError(err)
	}
	logger3, err := logsqlite.NewLogger("demo/httplog.sql", time.Second*3)
	if err != nil {
		gko.FatalError(err)
	}

	gko.LogInfo("Logging")
	for i := 0; i < 100000; i++ {
		logger0.SaveLog(makeLogEntry("logger0", "1/"+strconv.Itoa(i)))
		logger0.SaveLog(makeLogEntry("logger0", "2/"+strconv.Itoa(i)))
		logger1.SaveLog(makeLogEntry("logger1", "3/"+strconv.Itoa(i)))
		logger1.SaveLog(makeLogEntry("logger1", "3/"+strconv.Itoa(i)))
		logger2.SaveLog(makeLogEntry("logger2", "4/"+strconv.Itoa(i)))
		logger3.SaveLog(makeLogEntry("logger3", "5/"+strconv.Itoa(i)))
		logger2.SaveLog(makeLogEntry("logger2", "6/"+strconv.Itoa(i)))
		logger1.SaveLog(makeLogEntry("logger1", "7/"+strconv.Itoa(i)))
		logger0.SaveLog(makeLogEntry("logger0", "8/"+strconv.Itoa(i)))
	}
	gko.LogInfo("Logged")

	g.GET("/", func(c *gecko.Context) error {
		return c.StringOk(mensaje + "\n")
	})
	g.GET("/teapot", func(c *gecko.Context) error {
		return c.String(http.StatusTeapot, "I'm a teapot\n")
	})

	g.GET("/error", func(c *gecko.Context) error {
		return gko.ErrNoEncontrado().Msg("Error de prueba").Str("huehuehue")
	})

	g.GET("/o", func(c *gecko.Context) error {
		// Este método no registra el error en log.
		return c.String(200, "Hummm")
	})

	g.GET("/u", func(c *gecko.Context) error {
		return gko.ErrDatoIndef().Msg("Datos indefinidos")
	})

	tmpl, err := template.New("hola").Parse(`hola {{ .Nombre }}, mi sesion es {{ .Sesion.Nombre }} y mi edad es {{ .Sesion.Edad }}`)
	if err != nil {
		gko.LogError(err)
	}
	tmpl.Parse(`{{ define "base_layout" }}<!DOCTYPE html><html><head><title>Gecko</title></head><body>BASE {{ .Contenido }}</body></html>{{ end }}`)
	tmpl.Parse(`{{ define "error" }}<!DOCTYPE html><html><head><title>Gecko</title></head><body>ERROR {{ .Contenido }}</body></html>{{ end }}`)
	tmpl.Parse(`{{ define "mybold" }}<!DOCTYPE html><html><head><title>Gecko</title></head><body>Bolded {{ bolded .Titulo }}</body></html>{{ end }}`)
	g.Renderer = &renderer{tmpl: tmpl}

	g.GET("/bold", func(c *gecko.Context) error {
		data := map[string]any{
			"Titulo": "Hola *Mundo*",
		}
		return c.RenderOk("mybold", data)
	})

	g.GET("/render", func(c *gecko.Context) error {
		c.Sesion = &Sesion{
			Nombre: "Juan",
			Edad:   30,
		}
		data := map[string]any{
			"Nombre": "Mundo",
		}
		return c.RenderOk("hola", data)
	})

	if socket != "" {
		err := g.IniciarEnSocket(socket)
		if err != nil {
			gko.LogError(err)
		}

	} else if puerto > 0 {
		err := g.IniciarEnPuerto(puerto)
		if err != nil {
			gko.LogError(err)
		}

	} else {
		fmt.Println("No se especificó puerto ni socket")
	}
}

type Sesion struct {
	Nombre string
	Edad   int
}

type renderer struct {
	tmpl *template.Template
}

func (r *renderer) Render(w io.Writer, nombre string, data any, c *gecko.Context) error {
	return r.tmpl.ExecuteTemplate(w, nombre, data)
}
