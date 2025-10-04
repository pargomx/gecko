package gecko

import (
	"errors"
	"strings"
	"time"

	"github.com/pargomx/gecko/gkt"
)

// ================================================================ //
// ========== FORM VALUES ========================================= //

// Valor del form sin modificar.
func (c *Context) FormTalCual(name string) string {
	return c.request.FormValue(name)
}

// Valor del form espaciado simple sin saltos de línea.
func (c *Context) FormVal(name string) string {
	return gkt.SinEspaciosExtra(c.request.FormValue(name))
}

// Valor del form espaciado simple sin saltos de línea, sino el default.
func (c *Context) FormValDefault(name, defaultValue string) string {
	value := gkt.SinEspaciosExtra(c.request.FormValue(name))
	if value == "" {
		return defaultValue
	}
	return value
}

// Valor del form convertido a bool.
// Retorna false a menos de que el valor sea: "on", "true", "1".
func (c *Context) FormBool(name string) bool {
	return gkt.ToBool(c.request.FormValue(name))
}

// Valor del form convertido a entero.
func (c *Context) FormIntMust(name string) (int, error) {
	return gkt.ToInt(c.request.FormValue(name))
}

// Valor del form convertido a entero sin verificar error (default 0).
func (c *Context) FormInt(name string) int {
	num, _ := gkt.ToInt(c.request.FormValue(name))
	return num
}

// Valor del form convertido a uint64.
func (c *Context) FormUintMust(name string) (uint, error) {
	return gkt.ToUint(c.request.FormValue(name))
}

// Valor del form convertido a uint64 sin verificar error (default 0).
func (c *Context) FormUint(name string) uint {
	num, _ := gkt.ToUint(c.request.FormValue(name))
	return num
}

// Valor del form convertido a uint64.
func (c *Context) FormUint64Must(name string) (uint64, error) {
	return gkt.ToUint64(c.request.FormValue(name))
}

// Valor del form convertido a uint64 sin verificar error (default 0).
func (c *Context) FormUint64(name string) uint64 {
	num, _ := gkt.ToUint64(c.request.FormValue(name))
	return num
}

// Valor del form convertido a centavos.
func (c *Context) FormCentavos(name string) (int, error) {
	return gkt.ToCentavos(c.request.FormValue(name))
}

// Valor del form convertido a time desde una fecha 28/08/2022 o 2022-02-13.
func (c *Context) FormFecha(name string) (time.Time, error) {
	return gkt.ToFecha(c.request.FormValue(name))
}

// Valor del form formato fecha convertido a time, que puede estar indefinido.
func (c *Context) FormFechaNullable(name string) (*time.Time, error) {
	return gkt.ToFechaNullable(c.request.FormValue(name))
}

// Valor del form convertido a time desde una fecha con hora.
func (c *Context) FormFechaHora(name string) (time.Time, error) {
	return gkt.ToFechaHora(c.request.FormValue(name))
}

// Valor del form formato fecha con hora convertido a time, que puede estar indefinido.
func (c *Context) FormFechaHoraNullable(name string) (*time.Time, error) {
	return gkt.ToFechaHoraNullable(c.request.FormValue(name))
}

// Valor del form convertido a time.
func (c *Context) FormTime(name string, layout string) (time.Time, error) {
	return gkt.ToTime(c.request.FormValue(name), layout)
}

// Valor del form convertido a time, que puede estar indefinido.
func (c *Context) FormTimeNullable(name string, layout string) (*time.Time, error) {
	return gkt.ToTimeNullable(c.request.FormValue(name), layout)
}

// ================================================================ //

// Múltiples valores del form sin modificar.
func (c *Context) MultiFormTalCual(name string) []string {
	if c.request.PostForm == nil && c.request.Form == nil {
		c.request.ParseForm()
	}
	return c.request.Form[name]
}

// Múltiples valores del form espaciado simple sin saltos de línea.
func (c *Context) MultiFormVal(name string) []string {
	res := []string{}
	for _, v := range c.MultiFormTalCual(name) {
		res = append(res, gkt.SinEspaciosExtra(v))
	}
	return res
}

// Múltiples valores del form convertidos a enteros.
// No se agregan los valores que tengan errores en la conversión.
func (c *Context) MultiFormInt(name string) []int {
	res := []int{}
	for _, v := range c.MultiFormTalCual(name) {
		n, err := gkt.ToInt(v)
		if err != nil {
			continue
		}
		res = append(res, n)
	}
	return res
}

// Múltiples valores del form convertidos a enteros.
// Los valores deben ser números válidos todos.
func (c *Context) MultiFormIntMust(name string) ([]int, error) {
	res := []int{}
	for _, v := range c.MultiFormTalCual(name) {
		n, err := gkt.ToInt(v)
		if err != nil {
			return nil, errors.New("el valor [" + v + "] no es un número válido para [" + name + "]")
		}
		res = append(res, n)
	}
	return res, nil
}

// Múltiples valores del form convertidos a enteros.
// No se agregan los valores que tengan errores en la conversión.
func (c *Context) MultiFormUint(name string) []uint64 {
	res := []uint64{}
	for _, v := range c.MultiFormTalCual(name) {
		n, err := gkt.ToUint64(v)
		if err != nil {
			continue
		}
		res = append(res, n)
	}
	return res
}

// Múltiples valores del form convertidos a enteros.
// Los valores deben ser números válidos todos.
func (c *Context) MultiFormUintMust(name string) ([]uint64, error) {
	res := []uint64{}
	for _, v := range c.MultiFormTalCual(name) {
		n, err := gkt.ToUint64(v)
		if err != nil {
			return nil, errors.New("el valor [" + v + "] no es un número válido para [" + name + "]")
		}
		res = append(res, n)
	}
	return res, nil
}

// ================================================================ //

// Deprecated. Transformar texto en capa de aplicación, no en handler.
func (c *Context) FormUpper(name string) string {
	return strings.ToUpper(gkt.SinEspaciosExtra(c.request.FormValue(name)))
}

// Deprecated. Transformar texto en capa de aplicación, no en handler.
func (c *Context) FormLower(name string) string {
	return strings.ToLower(gkt.SinEspaciosExtra(c.request.FormValue(name)))
}
