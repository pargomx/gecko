package gko

import (
	"fmt"
	"slices"
)

// ================================================================ //
// ========== gko.Error =========================================== //

// El gko.Error facilita la construcción de errores con tres destinatarios
// diferentes a la vez: el usuario, el desarrollador, el sistema.
// Para crear un nuevo error se puede...
//
// Declarar una operación al inicio de la función:
//
//	op := gko.Op("MyFunc")
//	if err != nil { return op.Err(err) }
//
// Recibir y convertir un error normal:
//
//	return gko.Err(err).Msg("Algo salió mal")
//
// Usar un ErrorKey que se declare en el mismo paquete o en otro.
//
//	var ErrUserNotFound gko.ErrorKey = "user_not_found"
//	return ErrUserNotFound.Msgf("El usuario %v no existe", userID)
//
// Intenar llenar toda la información necesaria sin ser redundante.
//
// Es un tipo exportado para poder hacer type casting fuera del paquete
// y acceder a sus métodos de lectura.
type Error struct {
	// Identificadores del error.
	// Del más genérico/low_level hasta el más específico/high_level
	errKeys []ErrorKey

	// Mensaje amigable para el usuario.
	// Del más general al más específico.
	mensaje string

	// Error técnico para el desarrollador que no verá el usuario.
	// Desde de más alto nivel y general hasta el más bajo y específico.
	texto string

	// Stack de funciones que se estaban ejecutando.
	// Desde la primera invocada (high_level) a la última (low_level).
	operación string

	// Pares de claves=valor que dan contexto extra a la operación.
	valores string
}

// ================================================================ //
// ========== Instanciar errores ================================== //

// Convierte cualquier error en el tipo de gecko para poder usar sus métodos.
// NUNCA retorna nil. Tampoco hace Wrap al error, por lo que errors.Is() ni
// errors.As() funcionan con el error contenido. El gko.Error ofrece una
// funcionalidad similar y extendia con sus métodos.
func Err(err error) *Error {
	// Si no hay error, retornar uno vacío.
	if err == nil {
		return &Error{}
	}
	// Si ya es un error gecko, solo convertirlo.
	if errGk, ok := err.(*Error); ok {
		return errGk
	}
	// Si es un error normal, transformarlo.
	return &Error{
		texto: err.Error(),
	}
}

// Crea un gko.Error a partir de un gko.ErrorKey y el error proporcionado.
func (k ErrorKey) Err(err error) *Error {
	e := &Error{errKeys: []ErrorKey{k}}
	return e.Err(err)
}

// Genera una copia del error para poder agregar datos sin afectar a la
// variable original, por ejemplo cuando se está dentro de un loop.
//
//	op := op.Copy().Op("InLoop").Strf("loop %v", i)
func (e *Error) Copy() *Error {
	newErr := *e
	return &newErr
}

// Agregar error genérico a la cola, o combinar un *gko.Error de manera
// adecuada, suponiendo que el error recibido viene de la capa de ejecución
// inmediatamente inferior al error sobre el que se llama este método.
func (e *Error) Err(err error) *Error {
	// si no hay error nuevo, no hacer nada
	if err == nil {
		return e
	}
	// si el error no es de gecko solo agregar el texto.
	errGk, ok := err.(*Error)
	if !ok {
		e.Str(err.Error())
		return e
	}
	// si también es de gecko hay que combinarlos
	for _, key := range errGk.errKeys {
		e.Key(key)
	}
	if errGk.mensaje != "" {
		e.Msg(errGk.mensaje)
	}
	if errGk.texto != "" {
		e.Str(errGk.texto)
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

// ================================================================ //
// ========== ErrorKey ============================================ //

// ErrorKey identifica un tipo de error sin importar sus parámetros específicos.
// Por ejemplo:
//
//	var (
//	    UserNotFound gko.ErrorKey = "user_not_found"
//	    AgeTooLow    gko.ErrorKey = "age_low"
//	    AgeZero      gko.ErrorKey = "age_zero"
//	)
type ErrorKey string

// Errores comunes genéricos que están pensados para utilizarse como punto de
// partida para crear un nuevo gko.Error y traducir al dominio los errores de
// las capas más bajas de la aplicación (librerías externas) y pueda ser
// identificado por las capas más altas (comandos de aplicación, handlers, UI).
const (
	ErrInesperado   ErrorKey = "inesperado"    // Error desconocido, normalmente de una dependencia externa.
	ErrNoEncontrado ErrorKey = "not_found"     // No se encuentra un registro por su ID.
	ErrYaExiste     ErrorKey = "ya_existe"     // Ya existe un recurson con el mismo ID.
	ErrHayHuerfanos ErrorKey = "hay_huerfanos" // No se puede borrar porque tiene hijos.
	ErrTooManyReq   ErrorKey = "too_many_req"  // Se esperaba un solo registro y se encontraron muchos.
	ErrTooBig       ErrorKey = "too_big"       // Un archivo es demasiado pesado.
	ErrTooLong      ErrorKey = "too_long"      // Un string es demasiado largo.
	ErrDatoIndef    ErrorKey = "dato_indef"    // Un dato es obligatorio y no se recibió.
	ErrDatoInvalido ErrorKey = "dato_invalido" // Un dato no cumple con las reglas de validación.
	ErrNoSoportado  ErrorKey = "no_soportado"  // Un formato de archivo o dato no es soportado por el sistema.
	ErrNoAutorizado ErrorKey = "no_autorizado" // Un usuario no tiene permisos para realizar una acción.
	ErrTimeout      ErrorKey = "timeout"       // Una operación tarda más de lo esperado.
	ErrNoDisponible ErrorKey = "no_disponible" // Un servicio no está disponible.
	ErrNoSpaceLeft  ErrorKey = "no_space_left" // Se alcanzó la capacidad máxima.
	ErrAlEscribir   ErrorKey = "al_escribir"   // Error al escribir en un archivo.
	ErrAlLeer       ErrorKey = "al_leer"       // Error al leer un archivo.
)

// ErrorKey implementa la interfaz error
// devolviendo el string con el que fue declarada.
func (k ErrorKey) Error() string {
	return string(k)
}

// ================================================================ //
// ========== ErrorKey - Identificar errores ====================== //

// Agregar clave al error para que se pueda identificar
// en capas superiores de la aplicación.
func (e *Error) E(key ErrorKey) *Error {
	e.errKeys = append(e.errKeys, key)
	return e
}

// Agregar clave al error para que se pueda identificar
// en capas superiores de la aplicación.
func (e *Error) Key(key ErrorKey) *Error {
	e.errKeys = append(e.errKeys, key)
	return e
}

// Devuelve la clave de error, que es la última que se agregó y debería ser la
// más específica y cercana al dominio.
//
// Por ejemplo: [ErrDatoInvalido, ErrEdadInvalida, ErrEdadNegativa]
// devolvería ErrEdadNegativa.
//
// Si no hay ninguna, devuelve gko.ErrInesperado.
func (e *Error) ErrorKey() ErrorKey {
	if len(e.errKeys) == 0 {
		return ErrInesperado
	}
	return e.errKeys[len(e.errKeys)-1]
}

// Reporta si el error contiene tal clave que lo identifique.
// Similar a errors.Is pero en la implementación de gko.
func (e *Error) Contiene(key ErrorKey) bool {
	return slices.Contains(e.errKeys, key)
}

// Reporta si el error contiene tal clave que lo identifique.
// Similar a errors.Is pero en la implementación de gko.
//
// Nota: solo funciona cuando err es *gko.Error.
func Is(err error, key ErrorKey) bool {
	if errGk, ok := err.(*Error); ok {
		return slices.Contains(errGk.errKeys, key)
	}
	return false
}

// ================================================================ //
// ========== Al usuario ========================================== //

// Crea un nuevo gko.Error con este mensaje para el usuario.
func (k ErrorKey) Msg(txt string) *Error {
	e := &Error{errKeys: []ErrorKey{k}}
	e.Msg(txt)
	return e
}

// Crea un nuevo gko.Error con este mensaje para el usuario usando fmt.Sprintf().
func (k ErrorKey) Msgf(format string, a ...any) *Error {
	e := &Error{errKeys: []ErrorKey{k}}
	e.Msg(fmt.Sprintf(format, a...))
	return e
}

// Define o agrega un mensaje de error dirigido al usuario.
// Futuras llamadas se concatenan al inicio con ":".
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

// Define o agrega un mensaje de error dirigido al usuario usando fmt.Sprintf().
// Futuras llamadas se concatenan al inicio con ":".
func (e *Error) Msgf(format string, a ...any) *Error {
	e.Msg(fmt.Sprintf(format, a...))
	return e
}

// ================================================================ //
// ========== Al desarrollador ==================================== //

// Crea un nuevo gko.Error con este mensaje para el desarrollador.
func (k ErrorKey) Str(txt string) *Error {
	e := &Error{errKeys: []ErrorKey{k}}
	e.Str(txt)
	return e
}

// Crea un nuevo gko.Error con este mensaje para el desarrollador usando fmt.Sprintf().
func (k ErrorKey) Strf(format string, a ...any) *Error {
	e := &Error{errKeys: []ErrorKey{k}}
	e.Str(fmt.Sprintf(format, a...))
	return e
}

// Define o agrega un mensaje de error dirigido al desarrollador.
// Futuras llamadas se concatenan al inicio con ":".
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

// Define o agrega un mensaje de error dirigido al desarrollador usando fmt.Sprintf().
// Futuras llamadas se concatenan al inicio con ":".
func (e *Error) Strf(format string, a ...any) *Error {
	e.Str(fmt.Sprintf(format, a...))
	return e
}

// ================================================================ //
// ========== Al desarrollador - gko.Op =========================== //

// Identifica la operación que se está realizando (función o bloque) para que en
// caso de error, el desarrollador pueda conocer el stack de invocaciones que
// provocó el error.
func Op(op string) *Error {
	return &Error{operación: op}
}

// Definir la operación que se intenta ejecutar. Subsecuentes llamadas se
// concatenan a la derecha con " > " para formar un stack.
//
// El orden en que se llama .Op() es relevante: espera que la primera operación
// sea la función de más alto nivel, y las consiguientes sean operaciones cada
// vez más específicas o de capas inferiores. Tener esto en cuenta para
// encadenar operaciones y errores dentro de la en la misma línea o función.
//
// Por ejemplo:
//
//	op := gko.Op("app.CambiarEdad")
//	err := repoGetUsuario(userId) // returns .Op("repo.GetUsuario")
//	if err != nil {
//	    return op.Op("verifUsuario").Err(err).Op("enSerio")
//	}
//	// Resulta: "app.CambiarEdad > verifUsuario > repo.GetUsuario > enSerio"
func (e *Error) Op(op string) *Error {
	if op == "" {
		LogWarn("err.Op() con operación vacía")
		return e
	}
	if e.operación == "" {
		e.operación = op
	} else {
		e.operación = e.operación + " > " + op
	}
	return e
}

// ================================================================ //
// ========== Al desarrollador - Ctx ============================== //

// Agregar contexto en forma de "clave=valor".
// Subsecuentes llamadas se concatenan con un espacio.
func (e *Error) Ctx(key string, val any) *Error {
	if e.valores == "" {
		e.valores = fmt.Sprintf("%s=%v", key, val)
	} else {
		e.valores += fmt.Sprintf(" %s=%v", key, val)
	}
	return e
}
