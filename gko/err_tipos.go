package gko

func (e *Error) GetCodigoHTTP() int {
	switch {
	case e.Contiene(ErrInesperado):
		return 500
	case e.Contiene(ErrNoEncontrado):
		return 404
	case e.Contiene(ErrYaExiste):
		return 409
	case e.Contiene(ErrHayHuerfanos):
		return 409
	case e.Contiene(ErrTooManyReq):
		return 429
	case e.Contiene(ErrTooBig):
		return 400
	case e.Contiene(ErrTooLong):
		return 400
	case e.Contiene(ErrDatoIndef):
		return 400
	case e.Contiene(ErrDatoInvalido):
		return 400
	case e.Contiene(ErrNoSoportado):
		return 415
	case e.Contiene(ErrNoAutorizado):
		return 403
	case e.Contiene(ErrTimeout):
		return 408
	case e.Contiene(ErrNoDisponible):
		return 503
	case e.Contiene(ErrNoSpaceLeft):
		return 507
	case e.Contiene(ErrAlEscribir):
		return 503
	case e.Contiene(ErrAlLeer):
		return 503
	default:
		return 500
	}
}

func EsErrNotFound(err error) bool {
	return Err(err).Contiene(ErrNoEncontrado)
}
func EsErrAlreadyExists(err error) bool {
	return Err(err).Contiene(ErrYaExiste)
}
