package gko

func (e *Error) GetCodigoHTTP() int {
	switch e.tipo {
	case tipoErrInesperado:
		return 500
	case tipoErrNoEncontrado:
		return 404
	case tipoErrYaExiste:
		return 409
	case tipoErrHayHuerfanos:
		return 409
	case tipoErrTooManyReq:
		return 429
	case tipoErrTooBig:
		return 400
	case tipoErrTooLong:
		return 400
	case tipoErrDatoIndef:
		return 400
	case tipoErrDatoInvalido:
		return 400
	case tipoErrNoSoportado:
		return 415
	case tipoErrNoAutorizado:
		return 403
	case tipoErrTimeout:
		return 408
	case tipoErrNoDisponible:
		return 503
	case tipoErrNoSpaceLeft:
		return 507
	case tipoErrAlEscribir:
		return 503
	case tipoErrAlLeer:
		return 503
	default:
		return 500
	}
}

func EsErrNotFound(err error) bool {
	return Err(err).tipo == tipoErrNoEncontrado
}
func EsErrAlreadyExists(err error) bool {
	return Err(err).tipo == tipoErrYaExiste
}
