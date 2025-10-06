package gkt

import (
	"fmt"
	"strings"
	"time"

	"github.com/pargomx/gecko/gko"
)

const (
	FormatoFecha     = "2006-01-02"
	FormatoFechaHora = "2006-01-02 15:04:05"
	FormatoHora      = "15:04:05"
)

// Timezone America/Mexico_City
var TzMexico = mustLoadLocation("America/Mexico_City")

// Ahora en México.
func Now() time.Time {
	return time.Now().In(TzMexico)
}

// mustLoadLocation is a helper function that panics if the location cannot be loaded.
func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(fmt.Sprintf("failed to load location %s: %v", name, err))
	}
	return loc
}

// ================================================================ //
// ========== Fecha simple ======================================== //

const msgErrFormatoFecha string = "Utilice formato 25/12/2024 para fecha"

// Convierte el txt en tiempo usando el layout de fecha y
// zona horaria America/Mexico_City.
//
// Si no hay un valor para la fecha, se considera error.
//
// Utiliza time.ParseInLocation porque sino MySQL recibe la
// fecha en UTC cuando la espera en hora local, y esto hace que
// la fecha guardada sea un día anterior (6h) al esperado.
//
// Acepta formatos así:
//
//	"2006-01-02" "2006/01/02" "28-01-2006" "28/01/2006".
func ToFecha(txt string) (time.Time, error) {
	txt = SinEspaciosNinguno(txt)
	if txt == "" {
		return time.Time{}, gko.ErrDatoIndef
	}
	if len(txt) != 10 {
		return time.Time{}, gko.ErrDatoInvalido.Msg(msgErrFormatoFecha)
	}
	txt = strings.ReplaceAll(txt, "/", "-") // Permitir "2006/01/02"
	if txt[2:3] == "-" {                    // Permitir fecha volteada: "28-01-2006"
		txt = txt[6:] + "-" + txt[3:5] + "-" + txt[:2] // Siempre dejar layout "2006-01-02"
	}
	// ParseInLocation porque MySQL la espera en hora local.
	tiempo, err := time.ParseInLocation(FormatoFecha, txt, TzMexico)
	if err != nil {
		return time.Time{}, gko.ErrDatoInvalido.Err(err).Msg(msgErrFormatoFecha)
	}
	return tiempo, nil
}

// Convierte el txt en tiempo usando el layout de fecha y
// zona horaria America/Mexico_City.
//
// Si no hay valor para la fecha, se devuelve nil sin error.
//
// Utiliza time.ParseInLocation porque sino MySQL recibe la
// fecha en UTC cuando la espera en hora local, y esto hace que
// la fecha guardada sea un día anterior (6h) al esperado.
//
// Acepta formatos así:
//
//	"2006-01-02" "2006/01/02" "28-01-2006" "28/01/2006".
func ToFechaNullable(txt string) (*time.Time, error) {
	txt = SinEspaciosNinguno(txt)
	if txt == "" {
		return nil, nil
	}
	if len(txt) != 10 {
		return nil, gko.ErrDatoInvalido.Msg(msgErrFormatoFecha)
	}
	txt = strings.ReplaceAll(txt, "/", "-") // Permitir "2006/01/02"
	if txt[2:3] == "-" {                    // Permitir fecha volteada: "28-01-2006"
		txt = txt[6:] + "-" + txt[3:5] + "-" + txt[:2] // Siempre dejar layout "2006-01-02"
	}
	// ParseInLocation porque MySQL la espera en hora local.
	tiempo, err := time.ParseInLocation(FormatoFecha, txt, TzMexico)
	if err != nil {
		return nil, gko.ErrDatoInvalido.Err(err).Msg(msgErrFormatoFecha)
	}
	return &tiempo, nil
}

// ================================================================ //
// ========== Fecha y hora ======================================== //

const msgErrFormatoFechaHora string = "Utilice formato 2024-12-25 16:30:00 para fecha y hora"

// Convierte el txt en tiempo usando el layout de fecha y
// zona horaria America/Mexico_City.
//
// Si no hay un valor para la fecha, se considera error.
//
// Utiliza time.ParseInLocation porque sino MySQL recibe la
// fecha en UTC cuando la espera en hora local, y esto hace que
// la fecha guardada sea un día anterior (6h) al esperado.
//
// Acepta formatos así:
//
//	"2006-01-02 23:18:00", "2006/01/02 23:18:00", "2006-01-02T23:18:00"
func ToFechaHora(txt string) (time.Time, error) {
	txt = SinEspaciosExtra(txt)
	if txt == "" {
		return time.Time{}, gko.ErrDatoIndef
	}
	if len(txt) != 19 {
		return time.Time{}, gko.ErrDatoInvalido.Msg(msgErrFormatoFechaHora)
	}
	txt = strings.ReplaceAll(txt, "/", "-") // Permitir "2006/01/02 23:18:00"
	txt = strings.ReplaceAll(txt, "T", " ") // Permitir "2006-01-02T23:18:00"
	// ParseInLocation porque MySQL la espera en hora local.
	tiempo, err := time.ParseInLocation(FormatoFechaHora, txt, TzMexico)
	if err != nil {
		return time.Time{}, gko.ErrDatoInvalido.Err(err).Msg(msgErrFormatoFechaHora)
	}
	return tiempo, nil
}

// Convierte el txt en tiempo usando el layout de fecha y
// zona horaria America/Mexico_City.
//
// Si no hay valor para la fecha, se devuelve nil sin error.
//
// Utiliza time.ParseInLocation porque sino MySQL recibe la
// fecha en UTC cuando la espera en hora local, y esto hace que
// la fecha guardada sea un día anterior (6h) al esperado.
//
// Acepta formatos así:
//
//	"2006-01-02 23:18:00", "2006/01/02 23:18:00", "2006-01-02T23:18:00"
func ToFechaHoraNullable(txt string) (*time.Time, error) {
	txt = SinEspaciosExtra(txt)
	if txt == "" {
		return nil, nil
	}
	if len(txt) != 19 {
		return nil, gko.ErrDatoInvalido.Msg(msgErrFormatoFechaHora)
	}
	txt = strings.ReplaceAll(txt, "/", "-") // Permitir "2006/01/02 23:18:00"
	txt = strings.ReplaceAll(txt, "T", " ") // Permitir "2006-01-02T23:18:00"
	// ParseInLocation porque MySQL la espera en hora local.
	tiempo, err := time.ParseInLocation(FormatoFechaHora, txt, TzMexico)
	if err != nil {
		return nil, gko.ErrDatoInvalido.Err(err).Msg(msgErrFormatoFechaHora)
	}
	return &tiempo, nil
}

// ================================================================ //
// ========== Fecha custom ======================================== //

// Convierte el txt en tiempo usando el layout dado y la
// zona horaria America/Mexico_City.
//
// Si no hay un valor para la fecha, se considera error.
//
// Utiliza time.ParseInLocation porque sino MySQL recibe la
// fecha en UTC cuando la espera en hora local, y esto hace que
// la fecha guardada sea un día anterior (6h) al esperado.
//
// Admite separadores "2025-01-01" y "2025/01/01".
func ToTime(txt string, layout string) (time.Time, error) {
	txt = SinEspaciosNinguno(txt)
	if txt == "" {
		return time.Time{}, gko.ErrDatoIndef
	}
	txt = strings.ReplaceAll(txt, "/", "-") // Permitir "2006/01/02"
	// ParseInLocation porque MySQL la espera en hora local.
	tiempo, err := time.ParseInLocation(layout, txt, TzMexico)
	if err != nil {
		return time.Time{}, gko.ErrDatoInvalido.Err(err)
	}
	return tiempo, nil
}

// Convierte el txt en tiempo usando el layout dado y la
// zona horaria America/Mexico_City.
//
// Si no hay valor para la fecha, se devuelve nil sin error.
//
// Utiliza time.ParseInLocation porque sino MySQL recibe la
// fecha en UTC cuando la espera en hora local, y esto hace que
// la fecha guardada sea un día anterior (6h) al esperado.
//
// Admite separadores "2025-01-01" y "2025/01/01".
func ToTimeNullable(txt, layout string) (*time.Time, error) {
	txt = SinEspaciosExtra(txt)
	if txt == "" {
		return nil, nil
	}
	txt = strings.ReplaceAll(txt, "/", "-") // Permitir "2006/01/02"
	// ParseInLocation porque MySQL la espera en hora local.
	tiempo, err := time.ParseInLocation(layout, txt, TzMexico)
	if err != nil {
		return nil, gko.ErrDatoInvalido.Err(err)
	}
	return &tiempo, nil
}

// TimeToStringOrEmpty returns the formatted string of a *time.Time or empty if nil.
func TimeToStringOrEmpty(t *time.Time, format string) string {
	if t == nil {
		return ""
	}
	return t.Format(format)
}
