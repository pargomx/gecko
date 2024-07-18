package gko

import (
	"fmt"
)

type Error struct {
	tipo      int    // Define el tipo de error.
	mensaje   string // Mensaje amigable para el usuario.
	operación string // Funciones que se estaban ejecutando.
	valores   string // Claves=Valor que dan contexto a la operación.
	texto     string // Error técnico que no verá el usuario.
}

const (
	tipoErrInesperado   = iota // Error desconocido, normalmente de una dependencia externa.
	tipoErrNoEncontrado        // No se encuentra un registro por su ID.
	tipoErrYaExiste            // Ya existe un recurson con el mismo ID.
	tipoErrHayHuerfanos        // No se puede borrar porque tiene hijos.
	tipoErrTooManyReq          // Se esperaba un solo registro y se encontraron muchos.
	tipoErrTooBig              // Un archivo es demasiado pesado.
	tipoErrTooLong             // Un string es demasiado largo.
	tipoErrDatoIndef           // Un dato es obligatorio y no se recibió.
	tipoErrDatoInvalido        // Un dato no cumple con las reglas de validación.
	tipoErrNoSoportado         // Un formato de archivo o dato no es soportado por el sistema.
	tipoErrNoAutorizado        // Un usuario no tiene permisos para realizar una acción.
	tipoErrTimeout             // Una operación tarda más de lo esperado.
	tipoErrNoDisponible        // Un servicio no está disponible.
	tipoErrNoSpaceLeft         // Se alcanzó la capacidad máxima.
	tipoErrAlEscribir          // Error al escribir en un archivo.
	tipoErrAlLeer              // Error al leer un archivo.
)

// ========================================================== //
// ========== C O N S T R U C T O R E S ===================== //

// Definir la operación para tener el contexto en caso de error.
func Op(op string) *Error { return &Error{operación: op} }

func ErrInesperado() *Error   { return &Error{tipo: tipoErrInesperado} }
func ErrNoEncntrado() *Error  { return &Error{tipo: tipoErrNoEncontrado} }
func ErrYaExiste() *Error     { return &Error{tipo: tipoErrYaExiste} }
func ErrHayHuerfanos() *Error { return &Error{tipo: tipoErrHayHuerfanos} }
func ErrTooManyReq() *Error   { return &Error{tipo: tipoErrTooManyReq} }
func ErrTooBig() *Error       { return &Error{tipo: tipoErrTooBig} }
func ErrTooLong() *Error      { return &Error{tipo: tipoErrTooLong} }
func ErrDatoIndef() *Error    { return &Error{tipo: tipoErrDatoIndef} }
func ErrDatoInvalido() *Error { return &Error{tipo: tipoErrDatoInvalido} }
func ErrNoSoportado() *Error  { return &Error{tipo: tipoErrNoSoportado} }
func ErrNoAutorizado() *Error { return &Error{tipo: tipoErrNoAutorizado} }
func ErrTimeout() *Error      { return &Error{tipo: tipoErrTimeout} }
func ErrNoDisponible() *Error { return &Error{tipo: tipoErrNoDisponible} }
func ErrNoSpaceLeft() *Error  { return &Error{tipo: tipoErrNoSpaceLeft} }
func ErrAlEscribir() *Error   { return &Error{tipo: tipoErrAlEscribir} }
func ErrAlLeer() *Error       { return &Error{tipo: tipoErrAlLeer} }

// Convierte cualquier error en el tipo de gecko
// para poder usar sus métodos. NUNCA retorna nil.
func Err(err error) *Error {
	// Si no hay error, retornar uno vacío.
	if err == nil {
		return &Error{}
	}
	// Si ya es un error gecko, retornarlo.
	if errGk, ok := err.(*Error); ok {
		return errGk
	}
	// Si es un error normal, wrappearlo.
	return &Error{
		texto: err.Error(),
	}
}

// ========================================================== //
// ========== S E T T E R S ================================= //

// Define un nuevo status setTipo para el error.
// Subsecuentes llamadas sustituyen el código anterior.
func (e *Error) setTipo(code int) *Error {
	if code > 0 {
		e.tipo = code
	}
	return e
}

// Definir un error dirigido al desarrollador.
// Subsecuentes llamadas se concatenan con ":".
func (e *Error) Str(txt string) *Error {
	if txt == "" {
		LogWarn("err.Str() con mensaje vacío")
		return e
	}
	if e.texto == "" {
		e.texto = txt
	} else {
		e.texto = txt + ": " + e.texto
	}
	return e
}

// Definir un mensaje dirigido al usuario.
// Subsecuentes llamadas se concatenan con ":".
func (e *Error) Msg(msg string) *Error {
	if msg == "" {
		LogWarn("err.Msg() con mensaje vacío")
		return e
	}
	if e.mensaje == "" {
		e.mensaje = msg
	} else {
		e.mensaje = msg + ": " + e.mensaje
	}
	return e
}

// Definir un mensaje dirigido al usuario con fmt.Sprintf().
// Subsecuentes llamadas se concatenan con ":".
func (e *Error) Msgf(format string, a ...any) *Error {
	e.Msg(fmt.Sprintf(format, a...))
	return e
}

// Definir operación que se intenta ejecutar.
// Subsecuentes llamadas se concatenan con ">".
func (e *Error) Op(op string) *Error {
	if op == "" {
		LogWarn("err.Op() con operación vacía")
		return e
	}
	if e.operación == "" {
		e.operación = op
	} else {
		e.operación = op + " > " + e.operación
	}
	return e
}

// Agregar contexto en forma de "clave=valor".
// Subsecuentes llamadas se concatenan con " ".
func (e *Error) Ctx(key string, val any) *Error {
	if e.valores == "" {
		e.valores = fmt.Sprintf("%s=%v", key, val)
	} else {
		e.valores += fmt.Sprintf(" %s=%v", key, val)
	}
	return e
}

// Agregar error genérico o castear y combinar un error Gecko.
func (e *Error) Err(err error) *Error {

	// si no hay error nuevo, no hacer nada
	if err == nil {
		return e
	}

	// si el error no es de gecko solo wrappearlo
	errGk, ok := err.(*Error)
	if !ok {
		e.Str(err.Error())
		return e
	}

	// si también es de gecko hay que combinarlos
	if errGk.tipo > 0 {
		e.setTipo(errGk.tipo)
	}
	if errGk.texto != "" {
		e.Str(errGk.texto)
	}
	if errGk.mensaje != "" {
		e.Msg(errGk.mensaje)
	}
	if errGk.operación != "" {
		e.Op(errGk.operación)
	}
	if errGk.valores != "" {
		if e.valores == "" {
			e.valores = errGk.valores
		} else {
			e.valores += " " + errGk.valores
		}
	}
	return e
}
