package gecko

import (
	"errors"
	"fmt"
	"net/http"
)

type Gkerror struct {
	// HTTP status code que define el tipo de error
	codigo int

	// Mensaje para el usuario
	mensaje string

	// Operación que se estaba intentando realizar
	operación string

	// Contexto de la acción a realizar. Ej. "editar usuario_id=1234"
	contexto string

	// Error genérico de una dependencia externa
	err error
}

// Errors
var (
	ErrBadRequest                    = NewErr(http.StatusBadRequest)                    // HTTP 400 Bad Request
	ErrUnauthorized                  = NewErr(http.StatusUnauthorized)                  // HTTP 401 Unauthorized
	ErrPaymentRequired               = NewErr(http.StatusPaymentRequired)               // HTTP 402 Payment Required
	ErrForbidden                     = NewErr(http.StatusForbidden)                     // HTTP 403 Forbidden
	ErrNotFound                      = NewErr(http.StatusNotFound)                      // HTTP 404 Not Found
	ErrMethodNotAllowed              = NewErr(http.StatusMethodNotAllowed)              // HTTP 405 Method Not Allowed
	ErrNotAcceptable                 = NewErr(http.StatusNotAcceptable)                 // HTTP 406 Not Acceptable
	ErrProxyAuthRequired             = NewErr(http.StatusProxyAuthRequired)             // HTTP 407 Proxy AuthRequired
	ErrRequestTimeout                = NewErr(http.StatusRequestTimeout)                // HTTP 408 Request Timeout
	ErrConflict                      = NewErr(http.StatusConflict)                      // HTTP 409 Conflict
	ErrGone                          = NewErr(http.StatusGone)                          // HTTP 410 Gone
	ErrLengthRequired                = NewErr(http.StatusLengthRequired)                // HTTP 411 Length Required
	ErrPreconditionFailed            = NewErr(http.StatusPreconditionFailed)            // HTTP 412 Precondition Failed
	ErrStatusRequestEntityTooLarge   = NewErr(http.StatusRequestEntityTooLarge)         // HTTP 413 Payload Too Large
	ErrRequestURITooLong             = NewErr(http.StatusRequestURITooLong)             // HTTP 414 URI Too Long
	ErrUnsupportedMediaType          = NewErr(http.StatusUnsupportedMediaType)          // HTTP 415 Unsupported Media Type
	ErrRequestedRangeNotSatisfiable  = NewErr(http.StatusRequestedRangeNotSatisfiable)  // HTTP 416 Range Not Satisfiable
	ErrExpectationFailed             = NewErr(http.StatusExpectationFailed)             // HTTP 417 Expectation Failed
	ErrTeapot                        = NewErr(http.StatusTeapot)                        // HTTP 418 I'm a teapot
	ErrMisdirectedRequest            = NewErr(http.StatusMisdirectedRequest)            // HTTP 421 Misdirected Request
	ErrUnprocessableEntity           = NewErr(http.StatusUnprocessableEntity)           // HTTP 422 Unprocessable Entity
	ErrLocked                        = NewErr(http.StatusLocked)                        // HTTP 423 Locked
	ErrFailedDependency              = NewErr(http.StatusFailedDependency)              // HTTP 424 Failed Dependency
	ErrTooEarly                      = NewErr(http.StatusTooEarly)                      // HTTP 425 Too Early
	ErrUpgradeRequired               = NewErr(http.StatusUpgradeRequired)               // HTTP 426 Upgrade Required
	ErrPreconditionRequired          = NewErr(http.StatusPreconditionRequired)          // HTTP 428 Precondition Required
	ErrTooManyRequests               = NewErr(http.StatusTooManyRequests)               // HTTP 429 Too Many Requests
	ErrRequestHeaderFieldsTooLarge   = NewErr(http.StatusRequestHeaderFieldsTooLarge)   // HTTP 431 Request Header Fields Too Large
	ErrUnavailableForLegalReasons    = NewErr(http.StatusUnavailableForLegalReasons)    // HTTP 451 Unavailable For Legal Reasons
	ErrInternalServerError           = NewErr(http.StatusInternalServerError)           // HTTP 500 Internal Server Error
	ErrNotImplemented                = NewErr(http.StatusNotImplemented)                // HTTP 501 Not Implemented
	ErrBadGateway                    = NewErr(http.StatusBadGateway)                    // HTTP 502 Bad Gateway
	ErrServiceUnavailable            = NewErr(http.StatusServiceUnavailable)            // HTTP 503 Service Unavailable
	ErrGatewayTimeout                = NewErr(http.StatusGatewayTimeout)                // HTTP 504 Gateway Timeout
	ErrHTTPVersionNotSupported       = NewErr(http.StatusHTTPVersionNotSupported)       // HTTP 505 HTTP Version Not Supported
	ErrVariantAlsoNegotiates         = NewErr(http.StatusVariantAlsoNegotiates)         // HTTP 506 Variant Also Negotiates
	ErrInsufficientStorage           = NewErr(http.StatusInsufficientStorage)           // HTTP 507 Insufficient Storage
	ErrLoopDetected                  = NewErr(http.StatusLoopDetected)                  // HTTP 508 Loop Detected
	ErrNotExtended                   = NewErr(http.StatusNotExtended)                   // HTTP 510 Not Extended
	ErrNetworkAuthenticationRequired = NewErr(http.StatusNetworkAuthenticationRequired) // HTTP 511 Network Authentication Required

	ErrValidatorNotRegistered = errors.New("validator not registered")
	ErrRendererNotRegistered  = errors.New("renderer not registered")
	ErrInvalidRedirectCode    = errors.New("invalid redirect status code")
	ErrCookieNotFound         = errors.New("cookie not found")
	ErrInvalidCertOrKeyType   = errors.New("invalid cert or key type, must be string or []byte")
	ErrInvalidListenerNetwork = errors.New("invalid listener network")
)

func FatalErr(err error) {
	if err != nil {
		panic(err)
	}
}

func FatalFmt(msg string, args ...any) {
	panic(fmt.Sprintf(msg, args...))
}

// Error satisface la interfaz `error` componiendo el mensaje
// de una manera comprensible y completa para poner en los logs.
//
// Evitar visibilizar al usuario porque da todo el contexto.
func (e *Gkerror) Error() string {
	msg := ""
	if e.codigo > 0 {
		msg += fmt.Sprintf("[%d]", e.codigo)
	}
	if e.operación != "" {
		msg += " " + e.operación
	}
	if e.err != nil {
		msg += " " + e.err.Error()
	}
	if e.mensaje != "" && MostrarMensajeEnErrores {
		msg += ": " + e.mensaje + "."
	}
	if e.contexto != "" {
		msg += " {" + e.contexto + "}"
	}
	return msg
}

// ================================================================ //
// ========== S E T T E R S ======================================= //

// Define un nuevo status code para el error.
//
// Subsecuentes llamadas sustituyen el código anterior.
func (e *Gkerror) Code(code int) *Gkerror {
	e.codigo = code
	return e
}

func (e *Gkerror) GetCode() int {
	return e.codigo
}

// Mensaje dirigido al usuario.
// Subsecuentes llamadas se concatenan con `: `.
func (e *Gkerror) Msg(msg string) *Gkerror {
	if e.mensaje == "" {
		e.mensaje = msg
	} else {
		e.mensaje += ": " + msg
	}
	return e
}

// Mensaje dirigido al usuario.
// Subsecuentes llamadas se concatenan con `: `.
func (e *Gkerror) Msgf(format string, a ...any) *Gkerror {
	if e.mensaje == "" {
		e.mensaje = fmt.Sprintf(format, a...)
	} else {
		e.mensaje += ": " + fmt.Sprintf(format, a...)
	}
	return e
}

// Operación que se intentaba realizar.
// Subsecuentes llamadas se concatenan con ` > `.
func (e *Gkerror) Op(op string) *Gkerror {
	if e.operación == "" {
		e.operación = op
	} else {
		e.operación = op + " > " + e.operación
	}
	return e
}

// Contexto en forma de "clave=valor".
// Subsecuentes llamadas se concatenan con ` `.
func (e *Gkerror) Ctx(key string, val any) *Gkerror {
	if e.contexto == "" {
		e.contexto = fmt.Sprintf("%s=%v", key, val)
	} else {
		e.contexto += fmt.Sprintf(" %s=%v", key, val)
	}
	return e
}

func (e *Gkerror) Err(err error) *Gkerror {
	if err == nil {
		if e.err == nil {
			e.err = errors.New("err nil")
		} else {
			e.err = errors.Join(e.err, errors.New("err nil"))
		}
		return e
	}

	ne, esGecko := err.(*Gkerror)
	if !esGecko {
		if e.err == nil {
			e.err = err
		} else {
			e.err = errors.Join(e.err, err)
		}
		return e
	}

	// es gecko
	e.codigo = ne.codigo

	if ne.mensaje != "" {
		if ne.mensaje == "" {
			e.mensaje = ne.mensaje
		} else {
			e.mensaje += ": " + ne.mensaje
		}
	}

	if ne.operación != "" {
		if e.operación == "" {
			e.operación = ne.operación
		} else {
			e.operación = ne.operación + " > " + e.operación
		}
	}

	if ne.contexto != "" {
		if e.contexto == "" {
			e.contexto = ne.contexto
		} else {
			e.contexto += " " + ne.contexto
		}
	}

	if ne.err != nil {
		if e.err == nil {
			e.err = ne.err
		} else {
			e.err = errors.Join(e.err, ne.err)
		}
		return e
	}

	// fmt.Println("gk: " + err.Error())

	return e
}

// ================================================================ //
// ========== C O N S T R U C T O R E S =========================== //

// Nuevo error gecko con http status code.
func NewErr(code int) *Gkerror {
	return &Gkerror{
		codigo: code,
	}
}

// Nuevo error gecko con la operación que se intenta realizar.
func NewOp(op string) *Gkerror {
	return &Gkerror{
		operación: op,
	}
}

// ================================================================ //
// ========== A S S E R T I O N S ================================= //

// Convierte un interface error en *GeckoError usando type assertion.
func Err(err error) *Gkerror {
	if err == nil {
		return &Gkerror{
			err: errors.New("err nil"),
		}
	}
	if errGecko, ok := err.(*Gkerror); ok {
		return errGecko
	}
	return &Gkerror{
		err: err,
	}
}

// Reporta si el código del error es 404 NotFound.
func (e *Gkerror) EsNotFound() bool {
	if e == nil {
		return false
	}
	return e.codigo == http.StatusNotFound
}

// Reporta si el código del error es 409 Conflict.
func (e *Gkerror) EsAlreadyExists() bool {
	if e == nil {
		return false
	}
	return e.codigo == http.StatusConflict
}

// Reporta si el código del error concreto es 404 NotFound.
func EsErrNotFound(err error) bool {
	if err == nil {
		return false
	}
	e, ok := err.(*Gkerror)
	if !ok {
		return false
	}
	if e == nil {
		return false
	}
	return e.codigo == http.StatusNotFound
}

// Reporta si el código del error concreto es 409 Conflict.
func EsErrAlreadyExists(err error) bool {
	if err == nil {
		return false
	}
	e, ok := err.(*Gkerror)
	if !ok {
		return false
	}
	if e == nil {
		return false
	}
	return e.codigo == http.StatusConflict
}
