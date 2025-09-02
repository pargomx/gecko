package gko

// ================================================================ //
// ========== TX ================================================== //

// Resultado de una transacción a nivel de aplicación.
type TxResult struct {
	Events []Event
	Errors []*Error
}

// Agrega un potencial error al resultado.
func (s *TxResult) CaptureError(err error) {
	if err != nil {
		s.Errors = append(s.Errors, Err(err))
	}
}

// Reporta si el resultado contiene tal ErrorKey.
func (t TxResult) HasError(key ErrorKey) bool {
	for _, err := range t.Errors {
		if err.Contiene(key) {
			return true
		}
	}
	return false
}

// Reporta si el resultado contiene tal EventKey.
func (t TxResult) HasEvent(key EventKey) bool {
	for _, ev := range t.Events {
		if ev.EventKey == key {
			return true
		}
	}
	return false
}

// Reporta si el resultado contiene tal ErrorKey o EventKey.
func (t TxResult) HasKey(key string) bool {
	for _, err := range t.Errors {
		if err.Contiene(ErrorKey(key)) {
			return true
		}
	}
	for _, ev := range t.Events {
		if ev.EventKey == EventKey(key) {
			return true
		}
	}
	return false
}
