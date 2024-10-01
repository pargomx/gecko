package plantillas

import (
	"fmt"
	"time"
)

func timestamp() string {
	return time.Now().Format("2006-01-02 15:04:05 MST")
}

// ================================================================ //

// FechaText. TODO: mejorar y poner en español.
func FechaText(f time.Time) string {
	if f.IsZero() {
		return "Sin fecha"
	}
	return f.Format("02 Jan 2006 03:04 PM")
}

// ================================================================ //

func FechaEsp(f time.Time) string {
	mes := ""
	switch f.Month() {
	case time.January:
		mes = "enero"
	case time.February:
		mes = "febrero"
	case time.March:
		mes = "marzo"
	case time.April:
		mes = "abril"
	case time.May:
		mes = "mayo"
	case time.June:
		mes = "junio"
	case time.July:
		mes = "julio"
	case time.August:
		mes = "agosto"
	case time.September:
		mes = "septiembre"
	case time.October:
		mes = "octubre"
	case time.November:
		mes = "noviembre"
	case time.December:
		mes = "diciembre"
	default:
		mes = "error_en_mes"
	}
	return fmt.Sprintf("%v de %v de %v", f.Day(), mes, f.Year())
}

// ================================================================ //

// Convirte duración en hh:mm.
// Ejemplo: 1h1m45s -> "01:02"
func fmtHorasMinutos(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}

func MinutosToString(mins int) string {
	if mins == 0 {
		return ""
	}
	if mins < 60 {
		return fmt.Sprintf("%dm", mins)
	}
	h := mins / 60
	m := mins % 60
	if m == 0 {
		return fmt.Sprintf("%d h", h)
	}
	return fmt.Sprintf("%d:%02dm", h, m)
}

func SegundosToString(segs int) string {
	if segs == 0 {
		return "-"
	}
	if segs < 60 {
		return fmt.Sprintf("%ds", segs)
	}
	h := segs / 3600
	m := (segs % 3600) / 60
	s := segs % 60
	if h == 0 {
		return fmt.Sprintf("%dm %ds", m, s)
	}
	if m == 0 {
		return fmt.Sprintf("%dh", h)
	}
	if s == 0 {
		return fmt.Sprintf("%d:%02dm", h, m)
	}
	return fmt.Sprintf("%d:%02d:%02ds", h, m, s)
}
