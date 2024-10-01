package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
)

func main() {

	var mensaje, socket string
	var puerto int
	flag.IntVar(&puerto, "p", 8080, "Puerto TCP en el que escuchará el servidor")
	flag.StringVar(&socket, "s", "", "Socket Unix en el que escuchará el servidor")
	flag.StringVar(&mensaje, "m", "Hola mundo", "Mensaje que retornará el servidor")
	flag.Parse()

	fmt.Println("Preparando servidor")
	g := gecko.New()

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
