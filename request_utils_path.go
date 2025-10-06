package gecko

import (
	"time"

	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/gkoid"
	"github.com/pargomx/gecko/gkt"
)

// ================================================================ //
// ========== PATH PARAMS ========================================= //

// Valor del path sin modificar.
func (c *Context) PathTalCual(name string) string {
	return c.Param(name)
}

// Valor del path espaciado simple sin saltos de línea.
func (c *Context) PathVal(name string) string {
	return gkt.SinEspaciosExtra(c.Param(name))
}

// Valor del path espaciado simple sin saltos de línea, sino el default.
func (c *Context) PathValDefault(name, defaultValue string) string {
	value := gkt.SinEspaciosExtra(c.Param(name))
	if value == "" {
		return defaultValue
	}
	return value
}

// Valor del path convertido a bool.
// Retorna false a menos de que el valor sea: "on", "true", "1".
func (c *Context) PathBool(name string) bool {
	return gkt.ToBool(c.Param(name))
}

// Valor del path convertido a entero.
func (c *Context) PathIntMust(name string) (int, error) {
	return gkt.ToInt(c.Param(name))
}

// Valor del path convertido a entero sin verificar error (default 0).
func (c *Context) PathInt(name string) int {
	num, _ := gkt.ToInt(c.Param(name))
	return num
}

// Valor del path convertido a uint64.
func (c *Context) PathUintMust(name string) (uint, error) {
	return gkt.ToUint(c.Param(name))
}

// Valor del path convertido a uint64 sin verificar error (default 0).
func (c *Context) PathUint(name string) uint {
	num, _ := gkt.ToUint(c.Param(name))
	return num
}

// Valor del path convertido a uint64.
func (c *Context) PathUint64Must(name string) (uint64, error) {
	return gkt.ToUint64(c.Param(name))
}

// Valor del path convertido a uint64 sin verificar error (default 0).
func (c *Context) PathUint64(name string) uint64 {
	num, _ := gkt.ToUint64(c.Param(name))
	return num
}

// Valor del path convertido a centavos.
func (c *Context) PathCentavos(name string) (int, error) {
	return gkt.ToCentavos(c.Param(name))
}

// Valor del path convertido a time desde una fecha 28/08/2022 o 2022-02-13.
func (c *Context) PathFecha(name string, layout string) (time.Time, error) {
	return gkt.ToFecha(c.Param(name))
}

// Valor del path formato fecha convertido a time, que puede estar indefinido.
func (c *Context) PathFechaNullable(name string) (*time.Time, error) {
	return gkt.ToFechaNullable(c.Param(name))
}

// Valor del path convertido a time desde una fecha con hora.
func (c *Context) PathFechaHora(name string) (time.Time, error) {
	return gkt.ToFechaHora(c.Param(name))
}

// Valor del path formato fecha con hora convertido a time, que puede estar indefinido.
func (c *Context) PathFechaHoraNullable(name string) (*time.Time, error) {
	return gkt.ToFechaHoraNullable(c.Param(name))
}

// Valor del path convertido a time.
func (c *Context) PathTime(name string, layout string) (time.Time, error) {
	return gkt.ToTime(c.Param(name), layout)
}

// Valor del path convertido a time, que puede estar indefinido.
func (c *Context) PathTimeNullable(name string, layout string) (*time.Time, error) {
	return gkt.ToTimeNullable(c.Param(name), layout)
}

func (c *Context) PathDecimal(name string) gkoid.Decimal {
	id, err := gkoid.ParseDecimal(c.request.PathValue(name))
	if err != nil {
		gko.Err(err).Ctx("name", name).Log()
	}
	return id
}
func (c *Context) PathHex(name string) gkoid.Hex {
	id, err := gkoid.ParseHex(c.request.PathValue(name))
	if err != nil {
		gko.Err(err).Ctx("name", name).Log()
	}
	return id
}
func (c *Context) PathAlfanum(name string) gkoid.Alfanum {
	id, err := gkoid.ParseAlfanum(c.request.PathValue(name))
	if err != nil {
		gko.Err(err).Ctx("name", name).Log()
	}
	return id
}
