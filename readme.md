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


## Errores

Un error se puede definir como cualquier desviación o fallo en el comportamiento esperado de un programa o sistema. Estas fallas pueden variar en impacto y frecuencia dependiendo la causa y naturaleza:

- Errores de sintaxis: evitan la compilación o ejecución.

- Errores de tiempo de ejecución: condiciones inválidas como división entre 0 o dereferenciar pointers nulos; causan panics o errores fatales.

- Errores lógicos o semánticos: no se producen errores como tal pero el resultado de la programación no es el correcto.

- Errores de validación: cuando los datos introducidos al sistema no cumplen con el formato, rango o características requeridas.

- Errores de consistencia: cuando el comando o datos dados provocaría inconsistencias en el estado del sistema según las reglas de negocio.

- Errores de recurso no encontrado: not found, no rows, etc.

- Errores de autorización: por falta de privilegios para el recurso especificado, o un rol equivocado.

- Errores de protección: ya sea por límite de tiempo de espera, rápida repetición de solicitudes, tamaño del mensaje, etc.

- Errores de disponibilidad: cuando el sistema no tiene todos los componentes disponibles, no puede acceder a recursos necesarios en el backend, ya no hay espacio de almacenamiento, o directamente no está implementada la función.

- Errores de conexión: cuando el cliente no se puede conectar al sistema por razones particulares al usuario.

En cada caso varía la necesidad de transmitir y registrar información sobre el error tanto al usuario como al desarrollador:

- Errores graves que deben informarse de inmediato al desarrollador con el mayor contexto posible.

- Errores esperados por parte del usuario que deben dar retroalimentación clara al usuario sobre qué debe hacer diferente para solucionar el error. El desarrollador se beneficia del contexto para mejorar la interfaz o las validaciones.

- Errores esperados por parte de bots que intentan encontrar vulnerabilidades en el servicio.

