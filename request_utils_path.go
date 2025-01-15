package gecko

import (
	"time"
)

func (c *Context) PathParam(name string) string {
	return c.Param(name)
}

// ================================================================ //
// ========== PATH PARAMS ========================================= //

// Valor del path sin sanitizar.
func (c *Context) PathTalCual(name string) string {
	return c.PathParam(name)
}

// Valor del path sanitizado.
func (c *Context) PathVal(name string) string {
	return TxtSanitizar(c.PathParam(name))
}

// Valor del path sanitizado en mayúsculas.
func (c *Context) PathUpper(name string) string {
	return TxtUpper(c.PathParam(name))
}

// Valor del path sanitizado en minúsculas.
func (c *Context) PathLower(name string) string {
	return TxtLower(c.PathParam(name))
}

// Valor del path convertido a bool.
// Retorna false a menos de que el valor sea: "on", "true", "1".
func (c *Context) PathBool(name string) bool {
	return TxtBool(c.PathParam(name))
}

// Valor del path convertido a entero.
func (c *Context) PathIntMust(name string) (int, error) {
	return TxtInt(c.PathParam(name))
}

// Valor del path convertido a entero sin verificar error (default 0).
func (c *Context) PathInt(name string) int {
	num, _ := TxtInt(c.PathParam(name))
	return num
}

// Valor del path convertido a uint64.
func (c *Context) PathUintMust(name string) (uint64, error) {
	return TxtUint64(c.PathParam(name))
}

// Valor del path convertido a uint64 sin verificar error (default 0).
func (c *Context) PathUint(name string) uint64 {
	num, _ := TxtUint64(c.PathParam(name))
	return num
}

// Valor del path convertido a centavos.
func (c *Context) PathCentavos(name string) (int, error) {
	return TxtCentavos(c.PathParam(name))
}

// Valor del path convertido a time.
func (c *Context) PathTime(name string, layout string) (time.Time, error) {
	return TxtTime(c.PathParam(name), layout)
}

// Valor del path convertido a time, que puede estar indefinido.
func (c *Context) PathTimeNullable(name string, layout string) (*time.Time, error) {
	return TxtTimeNullable(c.PathParam(name), layout)
}

// Valor del path convertido a time desde una fecha 28/08/2022 o 2022-02-13.
func (c *Context) PathFecha(name string, layout string) (time.Time, error) {
	return TxtFecha(c.PathParam(name))
}

// Valor del path formato fecha convertido a time, que puede estar indefinido.
func (c *Context) PathFechaNullable(name string) (*time.Time, error) {
	return TxtFechaNullable(c.PathParam(name))
}
