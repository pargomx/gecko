package gko

import (
	"fmt"
)

type Error struct {
	tipo      int    // Define el tipo de error.
	mensaje   string // Mensaje amigable para el usuario.
	operación string // Funciones que se estaban ejecutando.
	contexto  string // Claves=Valor que dan contexto a la operación.
	err       string // Error técnico que no verá el usuario.
}

const (
	codeErrInesperado   = iota // Error desconocido, normalmente de una dependencia externa.
	codeErrNoEncntrado         // No se encuentra un registro por su ID.
	codeErrYaExiste            // Ya existe un recurson con el mismo ID.
	codeErrHayHuerfanos        // No se puede borrar porque tiene hijos.
	codeErrTooManyReq          // Se esperaba un solo registro y se encontraron muchos.
	codeErrTooBig              // Un archivo es demasiado pesado.
	codeErrTooLong             // Un string es demasiado largo.
	codeErrDatoIndef           // Un dato es obligatorio y no se recibió.
	codeErrDatoInvalido        // Un dato no cumple con las reglas de validación.
	codeErrNoSoportado         // Un formato de archivo o dato no es soportado por el sistema.
	codeErrNoAutorizado        // Un usuario no tiene permisos para realizar una acción.
	codeErrTimeout             // Una operación tarda más de lo esperado.
	codeErrNoDisponible        // Un servicio no está disponible.
	codeErrNoSpaceLeft         // Se alcanzó la capacidad máxima.
	codeErrAlEscribir          // Error al escribir en un archivo.
	codeErrAlLeer              // Error al leer un archivo.
)

// ========================================================== //
// ========== C O N S T R U C T O R E S ===================== //

// Definir la operación para tener el contexto en caso de error.
func Op(op string) *Error { return &Error{operación: op} }

func ErrInesperado() *Error   { return &Error{tipo: codeErrInesperado} }
func ErrNoEncntrado() *Error  { return &Error{tipo: codeErrNoEncntrado} }
func ErrYaExiste() *Error     { return &Error{tipo: codeErrYaExiste} }
func ErrHayHuerfanos() *Error { return &Error{tipo: codeErrHayHuerfanos} }
func ErrTooManyReq() *Error   { return &Error{tipo: codeErrTooManyReq} }
func ErrTooBig() *Error       { return &Error{tipo: codeErrTooBig} }
func ErrTooLong() *Error      { return &Error{tipo: codeErrTooLong} }
func ErrDatoIndef() *Error    { return &Error{tipo: codeErrDatoIndef} }
func ErrDatoInvalido() *Error { return &Error{tipo: codeErrDatoInvalido} }
func ErrNoSoportado() *Error  { return &Error{tipo: codeErrNoSoportado} }
func ErrNoAutorizado() *Error { return &Error{tipo: codeErrNoAutorizado} }
func ErrTimeout() *Error      { return &Error{tipo: codeErrTimeout} }
func ErrNoDisponible() *Error { return &Error{tipo: codeErrNoDisponible} }
func ErrNoSpaceLeft() *Error  { return &Error{tipo: codeErrNoSpaceLeft} }
func ErrAlEscribir() *Error   { return &Error{tipo: codeErrAlEscribir} }
func ErrAlLeer() *Error       { return &Error{tipo: codeErrAlLeer} }

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
		err: err.Error(),
	}
}

// ========================================================== //
// ========== S E T T E R S ================================= //

// Define un nuevo status code para el error.
// Subsecuentes llamadas sustituyen el código anterior.
func (e *Error) code(code int) *Error {
	if code > 0 {
		e.tipo = code
	}
	return e
}

// Definir un error dirigido al desarrollador.
// Subsecuentes llamadas se concatenan con ":".
func (e *Error) Str(err string) *Error {
	if err == "" {
		LogWarn("err.Str() con mensaje vacío")
		return e
	}
	if e.err == "" {
		e.err = err
	} else {
		e.err = err + ": " + e.err
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
	if e.contexto == "" {
		e.contexto = fmt.Sprintf("%s=%v", key, val)
	} else {
		e.contexto += fmt.Sprintf(" %s=%v", key, val)
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
		e.code(errGk.tipo)
	}
	if errGk.mensaje != "" {
		e.Msg(errGk.mensaje)
	}
	if errGk.operación != "" {
		e.Op(errGk.operación)
	}
	if errGk.contexto != "" {
		if e.contexto == "" {
			e.contexto = errGk.contexto
		} else {
			e.contexto += " " + errGk.contexto
		}
	}
	if errGk.err != "" {
		e.Str(errGk.err)
	}

	return e
}
