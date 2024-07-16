package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/pargomx/gecko"
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

	if socket != "" {
		err := g.IniciarEnSocket(socket)
		if err != nil {
			gecko.LogError(err)
		}

	} else if puerto > 0 {
		err := g.IniciarEnPuerto(puerto)
		if err != nil {
			gecko.LogError(err)
		}

	} else {
		fmt.Println("No se especificó puerto ni socket")
	}
}
