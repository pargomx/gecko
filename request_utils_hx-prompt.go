package gecko

import (
	"time"
)

// ================================================================ //
// ========== HX-Prompt Header ==================================== //

// Valor del header Hx-Prompt sin sanitizar.
func (c *Context) PromptTalCual() string {
	return c.request.Header.Get("HX-Prompt")
}

// Valor del header Hx-Prompt sanitizado.
func (c *Context) PromptVal() string {
	return TxtSanitizar(c.request.Header.Get("HX-Prompt"))
}

// Valor del header Hx-Prompt sanitizado en mayúsculas.
func (c *Context) PromptUpper() string {
	return TxtUpper(c.request.Header.Get("HX-Prompt"))
}

// Valor del header Hx-Prompt sanitizado en minúsculas.
func (c *Context) PromptLower() string {
	return TxtLower(c.request.Header.Get("HX-Prompt"))
}

// Valor del header Hx-Prompt convertido a bool.
// Retorna false a menos de que el valor sea: "on", "true", "1".
func (c *Context) PromptBool() bool {
	return TxtBool(c.request.Header.Get("HX-Prompt"))
}

// Valor del header Hx-Prompt convertido a entero.
func (c *Context) PromptIntMust() (int, error) {
	return TxtInt(c.request.Header.Get("HX-Prompt"))
}

// Valor del header Hx-Prompt convertido a entero sin verificar error (default 0).
func (c *Context) PromptInt() int {
	num, _ := TxtInt(c.request.Header.Get("HX-Prompt"))
	return num
}

// Valor del header Hx-Prompt convertido a uint64.
func (c *Context) PromptUintMust() (uint64, error) {
	return TxtUint64(c.request.Header.Get("HX-Prompt"))
}

// Valor del header Hx-Prompt convertido a uint64 sin verificar error (default 0).
func (c *Context) PromptUint() uint64 {
	num, _ := TxtUint64(c.request.Header.Get("HX-Prompt"))
	return num
}

// Valor del header Hx-Prompt convertido a centavos.
func (c *Context) PromptCentavos() (int, error) {
	return TxtCentavos(c.request.Header.Get("HX-Prompt"))
}

// Valor del header Hx-Prompt convertido a time.
func (c *Context) PromptTime(layout string) (time.Time, error) {
	return TxtTime(c.request.Header.Get("HX-Prompt"), layout)
}

// Valor del header Hx-Prompt convertido a time, que puede estar indefinido.
func (c *Context) PromptTimeNullable(layout string) (*time.Time, error) {
	return TxtTimeNullable(c.request.Header.Get("HX-Prompt"), layout)
}

// Valor del header Hx-Prompt convertido a time desde una fecha 28/08/2022 o 2022-02-13.
func (c *Context) PromptFecha(layout string) (time.Time, error) {
	return TxtFecha(c.request.Header.Get("HX-Prompt"))
}

// Valor del path formato fecha convertido a time, que puede estar indefinido.
func (c *Context) PromptFechaNullable() (*time.Time, error) {
	return TxtFechaNullable(c.request.Header.Get("HX-Prompt"))
}
