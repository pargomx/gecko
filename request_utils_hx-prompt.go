package gecko

import (
	"net/url"
	"time"

	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/gkt"
)

// ================================================================ //
// ========== HX-Prompt Header ==================================== //

// Valor del header Hx-Prompt-Encoded decodificado, o sino Hx-Prompt sin modificar.
func (c *Context) PromptTalCual() string {
	prompt, err := url.QueryUnescape(c.request.Header.Get(HxPromptEncoded))
	if err != nil {
		gko.Err(err).Strf("decoding HxPromptEncoded '%s'", c.request.Header.Get(HxPromptEncoded)).Log()
	}
	if prompt == "" {
		prompt = c.request.Header.Get(HxPrompt)
	}
	return prompt
}

// Valor del header Hx-Prompt espaciado simple sin saltos de línea.
func (c *Context) PromptVal() string {
	return gkt.SinEspaciosExtra(c.PromptTalCual())
}

// Valor del header Hx-Prompt espaciado simple sin saltos de línea, sino el default.
func (c *Context) PromptValDefault(name, defaultValue string) string {
	value := gkt.SinEspaciosExtra(c.PromptTalCual())
	if value == "" {
		return defaultValue
	}
	return value
}

// Valor del header Hx-Prompt convertido a bool.
// Retorna false a menos de que el valor sea: "on", "true", "1".
func (c *Context) PromptBool() bool {
	return gkt.ToBool(c.request.Header.Get(HxPrompt))
}

// Valor del header Hx-Prompt convertido a entero.
func (c *Context) PromptIntMust() (int, error) {
	return gkt.ToInt(c.request.Header.Get(HxPrompt))
}

// Valor del header Hx-Prompt convertido a entero sin verificar error (default 0).
func (c *Context) PromptInt() int {
	num, _ := gkt.ToInt(c.request.Header.Get(HxPrompt))
	return num
}

// Valor del header Hx-Prompt convertido a uint64.
func (c *Context) PromptUintMust() (uint64, error) {
	return gkt.ToUint64(c.request.Header.Get(HxPrompt))
}

// Valor del header Hx-Prompt convertido a uint64 sin verificar error (default 0).
func (c *Context) PromptUint() uint64 {
	num, _ := gkt.ToUint64(c.request.Header.Get(HxPrompt))
	return num
}

// Valor del header Hx-Prompt convertido a centavos.
func (c *Context) PromptCentavos() (int, error) {
	return gkt.ToCentavos(c.request.Header.Get(HxPrompt))
}

// Valor del header Hx-Prompt convertido a time desde una fecha 28/08/2022 o 2022-02-13.
func (c *Context) PromptFecha(layout string) (time.Time, error) {
	return gkt.ToFecha(c.request.Header.Get(HxPrompt))
}

// Valor del header Hx-Prompt formato fecha convertido a time, que puede estar indefinido.
func (c *Context) PromptFechaNullable() (*time.Time, error) {
	return gkt.ToFechaNullable(c.request.Header.Get(HxPrompt))
}

// Valor del header Hx-Prompt convertido a time desde una fecha con hora.
func (c *Context) PromptFechaHora(name string) (time.Time, error) {
	return gkt.ToFechaHora(c.request.Header.Get(HxPrompt))
}

// Valor del header Hx-Prompt formato fecha con hora convertido a time, que puede estar indefinido.
func (c *Context) PromptFechaHoraNullable(name string) (*time.Time, error) {
	return gkt.ToFechaHoraNullable(c.request.Header.Get(HxPrompt))
}

// Valor del header Hx-Prompt convertido a time.
func (c *Context) PromptTime(layout string) (time.Time, error) {
	return gkt.ToTime(c.request.Header.Get(HxPrompt), layout)
}

// Valor del header Hx-Prompt convertido a time, que puede estar indefinido.
func (c *Context) PromptTimeNullable(layout string) (*time.Time, error) {
	return gkt.ToTimeNullable(c.request.Header.Get(HxPrompt), layout)
}
