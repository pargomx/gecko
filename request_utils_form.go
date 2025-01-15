package gecko

import (
	"errors"
	"time"
)

// ================================================================ //
// ========== FORM VALUES ========================================= //

// Valor del form sin sanitizar.
func (c *Context) FormTalCual(name string) string {
	return c.request.FormValue(name)
}

// Valor del form sanitizado.
func (c *Context) FormVal(name string) string {
	return TxtSanitizar(c.request.FormValue(name))
}

// Valor del form sanitizado en mayúsculas.
func (c *Context) FormUpper(name string) string {
	return TxtUpper(c.request.FormValue(name))
}

// Valor del form sanitizado en minúsculas.
func (c *Context) FormLower(name string) string {
	return TxtLower(c.request.FormValue(name))
}

// Valor del form convertido a bool.
// Retorna false a menos de que el valor sea: "on", "true", "1".
func (c *Context) FormBool(name string) bool {
	return TxtBool(c.request.FormValue(name))
}

// Valor del form convertido a entero.
func (c *Context) FormIntMust(name string) (int, error) {
	return TxtInt(c.request.FormValue(name))
}

// Valor del form convertido a entero sin verificar error (default 0).
func (c *Context) FormInt(name string) int {
	num, _ := TxtInt(c.request.FormValue(name))
	return num
}

// Valor del form convertido a uint64.
func (c *Context) FormUintMust(name string) (uint64, error) {
	return TxtUint64(c.request.FormValue(name))
}

// Valor del form convertido a uint64 sin verificar error (default 0).
func (c *Context) FormUint(name string) uint64 {
	num, _ := TxtUint64(c.request.FormValue(name))
	return num
}

// Valor del form convertido a centavos.
func (c *Context) FormCentavos(name string) (int, error) {
	return TxtCentavos(c.request.FormValue(name))
}

// Valor del form convertido a time.
func (c *Context) FormTime(name string, layout string) (time.Time, error) {
	return TxtTime(c.request.FormValue(name), layout)
}

// Valor del form convertido a time, que puede estar indefinido.
func (c *Context) FormTimeNullable(name string, layout string) (*time.Time, error) {
	return TxtTimeNullable(c.request.FormValue(name), layout)
}

// Valor del form convertido a time desde una fecha 28/08/2022 o 2022-02-13.
func (c *Context) FormFecha(name string) (time.Time, error) {
	return TxtFecha(c.request.FormValue(name))
}

// Valor del path formato fecha convertido a time, que puede estar indefinido.
func (c *Context) FormFechaNullable(name string) (*time.Time, error) {
	return TxtFechaNullable(c.FormValue(name))
}

// ================================================================ //

// Múltiples valores sin sanitizar obtenidos del form.
func (c *Context) MultiFormTalCual(key string) []string {
	c.request.ParseForm()
	return c.request.Form[key]
}

// Múltiples valores sanitizados obtenidos del form.
func (c *Context) MultiFormVal(name string) []string {
	res := []string{}
	c.request.ParseForm()
	for _, v := range c.request.Form[name] {
		res = append(res, TxtSanitizar(v))
	}
	return res
}

// Múltiples valores sanitizados en mayúsculas obtenidos del form.
func (c *Context) MultiFormUpper(name string) []string {
	res := []string{}
	c.request.ParseForm()
	for _, v := range c.request.Form[name] {
		res = append(res, TxtUpper(v))
	}
	return res
}

// Múltiples valores sanitizados en minúsculas obtenidos del form.
func (c *Context) MultiFormLower(name string) []string {
	res := []string{}
	c.request.ParseForm()
	for _, v := range c.request.Form[name] {
		res = append(res, TxtLower(v))
	}
	return res
}

// Múltiples valores convertidos a enteros obtenidos del form.
// No se agregan los valores que tengan errores en la conversión.
func (c *Context) MultiFormInt(name string) []int {
	res := []int{}
	c.request.ParseForm()
	for _, v := range c.request.Form[name] {
		n, err := TxtInt(v)
		if err != nil {
			continue
		}
		res = append(res, n)
	}
	return res
}

// Múltiples valores convertidos a enteros obtenidos del form.
// Los valores deben ser números válidos todos.
func (c *Context) MultiFormIntMust(name string) ([]int, error) {
	res := []int{}
	c.request.ParseForm()
	for _, v := range c.request.Form[name] {
		n, err := TxtInt(v)
		if err != nil {
			return nil, errors.New("el valor [" + v + "] no es un número válido para [" + name + "]")
		}
		res = append(res, n)
	}
	return res, nil
}

// Múltiples valores convertidos a enteros obtenidos del form.
// No se agregan los valores que tengan errores en la conversión.
func (c *Context) MultiFormUint(name string) []uint64 {
	res := []uint64{}
	c.request.ParseForm()
	for _, v := range c.request.Form[name] {
		n, err := TxtUint64(v)
		if err != nil {
			continue
		}
		res = append(res, n)
	}
	return res
}

// Múltiples valores convertidos a enteros obtenidos del form.
// Los valores deben ser números válidos todos.
func (c *Context) MultiFormUintMust(name string) ([]uint64, error) {
	res := []uint64{}
	c.request.ParseForm()
	for _, v := range c.request.Form[name] {
		n, err := TxtUint64(v)
		if err != nil {
			return nil, errors.New("el valor [" + v + "] no es un número válido para [" + name + "]")
		}
		res = append(res, n)
	}
	return res, nil
}

// FormValDefault returns the form field value or default value
// for the provided name. Note: it does not distinguish if form
// had no value by that name or value was empty string.
func (c *Context) FormValDefault(name, defaultValue string) string {
	if c.query == nil {
		c.query = c.request.URL.Query()
	}
	value := c.query.Get(name)
	if value == "" {
		return defaultValue
	}
	return value
}
