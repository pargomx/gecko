package gko

import (
	"fmt"
	"os"
	"time"
)

// ================================================================ //
// ========== TIME ================================================ //
func timestamp() string {
	if PrintLogTimestamps {
		return time.Now().Format("2006-01-02 15:04:05") + " "
	}
	return ""
}

// Indica si se imprimirá la timestamp al inicio
// de cada entrada de log al Stdout.
var PrintLogTimestamps bool = true

// ================================================================ //
// ========== LOG ================================================= //

// 2020-11-25 18:54:32 [DEBUG] Información detallada cyan.
func LogDebug(a ...any) {
	println(timestamp() + cCyan + "[DEBUG] " + rWhite + fmt.Sprint(a...) + reset)
}
func LogDebugf(format string, a ...any) {
	println(timestamp() + cCyan + "[DEBUG] " + reset + fmt.Sprintf(format, a...) + reset)
}

// Okey LOG LISTO al terminar con éxito una función
func LogOkey(a ...any) {
	println(timestamp() + cGreen + "[LISTO] " + reset + fmt.Sprint(a...))
}
func LogOkeyf(format string, a ...any) {
	println(timestamp() + cGreen + "[LISTO] " + reset + fmt.Sprintf(format, a...))
}

// 2020-11-25 18:54:32 [EVENT] Algo importante sucede cyan bold.
func LogEvento(a ...any) {
	println(timestamp() + cCyan + "[EVENT] " + bold + fmt.Sprint(a...) + reset)
}
func LogEventof(format string, a ...any) {
	println(timestamp() + cCyan + "[EVENT] " + bold + fmt.Sprintf(format, a...) + reset)
}

// 2020-11-25 18:54:32 [INFOR] Algo interesante sucede cyan.
func LogInfo(a ...any) {
	println(timestamp() + cCyan + "[INFOR] " + rWhite + fmt.Sprint(a...) + reset)
}
func LogInfof(format string, a ...any) {
	println(timestamp() + cCyan + "[INFOR] " + rWhite + fmt.Sprintf(format, a...) + reset)
}

// 2020-11-25 18:54:32 [AVISO] Algo está sucediendo yellow.
func LogWarn(a ...any) {
	println(timestamp() + cYellow + "[AVISO] " + reset + fmt.Sprint(a...))
}
func LogWarnf(format string, a ...any) {
	println(timestamp() + cYellow + "[AVISO] " + reset + fmt.Sprintf(format, a...))
}

// 2020-11-25 18:54:32 [ALERT] Algo no está bien yellow.
func LogAlert(a ...any) {
	println(timestamp() + cYellow + "[ALERT] " + reset + fmt.Sprint(a...))
}
func LogAlertf(format string, a ...any) {
	println(timestamp() + cYellow + "[ALERT] " + reset + fmt.Sprintf(format, a...))
}

// 2020-11-25 18:54:32 [ALERT] Falló algo y se canceló algo yellow.
func LogAbort(a ...any) {
	println(timestamp() + cYellow + "[ABORT] " + reset + fmt.Sprint(a...))
}
func LogAbortf(format string, a ...any) {
	println(timestamp() + cYellow + "[ABORT] " + reset + fmt.Sprintf(format, a...))
}

// 2020-11-25 18:54:32 [FATAL] No se puede continuar la ejecución.
func FatalExit(a ...any) {
	println(timestamp() + cRed + "[FATAL] " + bold + fmt.Sprint(a...) + reset)
	os.Exit(1)
}
func FatalExitf(format string, a ...any) {
	println(timestamp() + cRed + "[FATAL] " + bold + fmt.Sprintf(format, a...) + reset)
	os.Exit(1)
}
func FatalError(err error) {
	LogError(err)
	os.Exit(1)
}

// ================================================================ //
// ========== ERROR =============================================== //

// Error satisface la interfaz `error` componiendo el mensaje
// de una manera comprensible y completa para poner en los logs.
// No está pensado para el usuario porque da mucha información.
func (e *Error) Error() string {
	msg := ""
	if e.tipo > 0 {
		msg += fmt.Sprintf("[%d]", e.tipo)
	}
	if e.operación != "" {
		msg += " " + e.operación
	}
	if e.mensaje != "" {
		msg += ": " + e.mensaje + "."
	}
	if e.valores != "" {
		msg += " {" + e.valores + "}"
	}
	if e.texto != "" {
		msg += " " + e.texto
	}
	return msg
}

// Imprime el error como estructura con toda su información.
func (e *Error) Print() {
	fmt.Printf("gko.Error{"+
		"\n\ttipo: "+cPurple+"%d"+reset+
		"\n\tmsg: "+cPurple+"%s"+reset+
		"\n\tops: "+cPurple+"%s"+reset+
		"\n\tctx: "+cPurple+"%s"+reset+
		"\n\ttxt: "+cPurple+"%v"+reset+
		"\n}\n",
		e.tipo, e.mensaje, e.operación, e.valores, e.texto,
	)
}

// Devuelve un mensaje para presentar al usuario.
func (e *Error) GetMensaje() string {
	if e.mensaje != "" {
		return e.mensaje + "."
	} else if e.texto != "" {
		return e.texto
	} else if e.tipo > 0 {
		return fmt.Sprintf("Error %d", e.tipo)
	}
	return "Hubo un error, por favor contacta a soporte."
}

// Imprime el error en la consola.
// Alias para gko.Err(err).Log()
func LogError(err error) {
	Err(err).Log()
}

// Log imprime el error gecko con colores en la consola.
func (e *Error) Log() {

	// 2020-11-25 18:54:32
	msg := timestamp()

	// 2020-11-25 18:54:32 [ERROR] (404)
	if e.tipo == 0 {
		msg += bRed + "[ERROR]" + reset
	} else {
		msg += bRed + fmt.Sprintf("[ERROR] (%d)", e.tipo) + reset
	}

	// 2020-11-25 18:54:32 [ERROR] (404) DoSomething > GetRecord
	if e.operación != "" {
		msg += " " + rYellow + e.operación
	}

	// [ERROR] (404) DoSomething > GetRecord: Usuario no encontrado.
	if e.mensaje != "" {
		msg += " " + bRed + e.mensaje + "." + reset
	}

	// [ERROR] (404) DoSomething > GetRecord: Usuario no encontrado. {id=123}
	if e.valores != "" {
		msg += " " + rPurple + e.valores
	}

	// [ERROR] (404) DoSomething > GetRecord: Usuario no encontrado. {id=123} sql: no rows
	if e.texto != "" {
		msg += " " + rRed + e.texto
	}

	println(msg + reset)
}
