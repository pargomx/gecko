# Gecko

Gecko es una librería escrita en Go que proporciona utilidades para manejar solicitudes HTTP y responderlas con plantillas HTML, errores personalizados, y más. Está basada en la librería estándar de Go, lo que la hace fácil de integrar y usar en cualquier proyecto Go.

## Características

- Manejo de solicitudes HTTP
- Respuestas con plantillas HTML
- Errores personalizados
- Y más...

## Changelog

Para ver los cambios realizados en cada versión de Gecko, consulta el archivo [version.md](version.md)

## Uso

Para usar Gecko en tu proyecto, simplemente importa la librería y comienza a usar sus funciones. Aquí hay un ejemplo de cómo hacerlo:

```go
import "github.com/pargomx/gecko"

func main() {
	g := gecko.New()

	g.GET("/", func(c *gecko.Context) error {
		return c.StringOk("Inicio")
	})

	g.GET("/hola", func(c *gecko.Context) error {
		return c.StringOk("Hola mundo")
	})

	g.GET("/saludo/{a}", func(c *gecko.Context) error {
		return c.StringOk("Hola, " + c.Param("a"))
	})

	g.Static("/assets")

	err := g.IniciarServidor()
	if err != nil {
		panic(err)
	}
}
```