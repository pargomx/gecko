package gecko

import (
	"errors"
	"time"

	"github.com/pargomx/gecko/gkt"
)

// ================================================================ //
// ========== QUERY PARAMS ======================================== //

// Valor del query sin modificar.
func (c *Context) QueryTalCual(name string) string {
	return c.QueryParam(name)
}

// Valor del query espaciado simple sin saltos de línea.
func (c *Context) QueryVal(name string) string {
	return gkt.SinEspaciosExtra(c.QueryParam(name))
}

// Valor del query espaciado simple sin saltos de línea, sino el default.
func (c *Context) QueryParamDefault(name, defaultValue string) string {
	value := gkt.SinEspaciosExtra(c.QueryParam(name))
	if value == "" {
		return defaultValue
	}
	return value
}

// Valor del query convertido a bool.
// Retorna false a menos de que el valor sea: "on", "true", "1".
func (c *Context) QueryBool(name string) bool {
	return gkt.ToBool(c.QueryParam(name))
}

// Valor del query convertido a entero.
func (c *Context) QueryIntMust(name string) (int, error) {
	return gkt.ToInt(c.QueryParam(name))
}

// Valor del query convertido a entero sin verificar error (default 0).
func (c *Context) QueryInt(name string) int {
	num, _ := gkt.ToInt(c.QueryParam(name))
	return num
}

// Valor del query convertido a uint64.
func (c *Context) QueryUintMust(name string) (uint64, error) {
	return gkt.ToUint64(c.QueryParam(name))
}

// Valor del query convertido a uint64 sin verificar error (default 0).
func (c *Context) QueryUint(name string) uint64 {
	num, _ := gkt.ToUint64(c.QueryParam(name))
	return num
}

// Valor del query convertido a centavos.
func (c *Context) QueryCentavos(name string) (int, error) {
	return gkt.ToCentavos(c.QueryParam(name))
}

// Valor del query convertido a time desde una fecha 28/08/2022 o 2022-02-13.
func (c *Context) QueryFecha(name string, layout string) (time.Time, error) {
	return gkt.ToFecha(c.QueryParam(name))
}

// Valor del query formato fecha convertido a time, que puede estar indefinido.
func (c *Context) QueryFechaNullable(name string) (*time.Time, error) {
	return gkt.ToFechaNullable(c.QueryParam(name))
}

// Valor del query convertido a time desde una fecha con hora.
func (c *Context) QueryFechaHora(name string) (time.Time, error) {
	return gkt.ToFechaHora(c.QueryParam(name))
}

// Valor del query formato fecha con hora convertido a time, que puede estar indefinido.
func (c *Context) QueryFechaHoraNullable(name string) (*time.Time, error) {
	return gkt.ToFechaHoraNullable(c.QueryParam(name))
}

// Valor del query convertido a time.
func (c *Context) QueryTime(name string, layout string) (time.Time, error) {
	return gkt.ToTime(c.QueryParam(name), layout)
}

// Valor del query convertido a time, que puede estar indefinido.
func (c *Context) QueryTimeNullable(name string, layout string) (*time.Time, error) {
	return gkt.ToTimeNullable(c.QueryParam(name), layout)
}

// ================================================================ //

// Múltiples valores obtenidos del query sin modificar.
func (c *Context) MultiQueryTalCual(name string) []string {
	if c.query == nil {
		c.query = c.request.URL.Query()
	}
	return c.query[name]
}

// Múltiples valores obtenidos del query espaciado simple sin saltos de línea.
func (c *Context) MultiQueryVal(name string) []string {
	res := []string{}
	for _, v := range c.MultiQueryTalCual(name) {
		res = append(res, gkt.SinEspaciosExtra(v))
	}
	return res
}

// Múltiples valores obtenidos del query convertidos a enteros.
// No se agregan los valores que tengan errores en la conversión.
func (c *Context) MultiQueryInt(name string) []int {
	res := []int{}
	for _, v := range c.MultiQueryTalCual(name) {
		n, err := gkt.ToInt(v)
		if err != nil {
			continue
		}
		res = append(res, n)
	}
	return res
}

// Múltiples valores obtenidos del query convertidos a enteros.
// Los valores deben ser números válidos todos.
func (c *Context) MultiQueryIntMust(name string) ([]int, error) {
	res := []int{}
	for _, v := range c.MultiQueryTalCual(name) {
		n, err := gkt.ToInt(v)
		if err != nil {
			return nil, errors.New("el valor [" + v + "] no es un número válido para [" + name + "]")
		}
		res = append(res, n)
	}
	return res, nil
}

// Múltiples valores obtenidos del query convertidos a enteros.
// No se agregan los valores que tengan errores en la conversión.
func (c *Context) MultiQueryUint(name string) []uint64 {
	res := []uint64{}
	for _, v := range c.MultiQueryTalCual(name) {
		n, err := gkt.ToUint64(v)
		if err != nil {
			continue
		}
		res = append(res, n)
	}
	return res
}

// Múltiples valores obtenidos del query convertidos a enteros.
// Los valores deben ser números válidos todos.
func (c *Context) MultiQueryUintMust(name string) ([]uint64, error) {
	res := []uint64{}
	for _, v := range c.MultiQueryTalCual(name) {
		n, err := gkt.ToUint64(v)
		if err != nil {
			return nil, errors.New("el valor [" + v + "] no es un número válido para [" + name + "]")
		}
		res = append(res, n)
	}
	return res, nil
}
