package gko

import (
	"fmt"
	"reflect"
	"time"
)

// ================================================================ //
// ========== Evento ============================================== //

// EventKey identifica a la consecuencia de un comando sin importar sus
// parámetros específicos. Por ejemplo:
//
//	var (
//		AgeModified    gko.EventKey = "age_modified"
//		UserRegistered gko.EventKey = "user_registered"
//		ProductSold    gko.EventKey = "product_sold"
//	)
type EventKey string

type evento struct {
	EventID   string    `json:"id"`
	EventKey  EventKey  `json:"type"`
	Body      any       `json:"body"`
	Mensaje   string    `json:"-"`
	Timestamp time.Time `json:"ts"`
}

// Agrega una estructura con los argumentos del evento (body).
func (k EventKey) WithArgs(argsStruct any) evento {
	e := evento{
		EventKey: k,
		Body:     argsStruct,
	}
	return e
}

// Agrega un mensaje para el usuario describiendo el evento.
func (k EventKey) Msgf(format string, a ...any) evento {
	e := evento{
		EventKey: k,
		Mensaje:  fmt.Sprintf(format, a...),
	}
	return e
}

// Agrega un mensaje para el usuario describiendo el evento.
func (e evento) Msgf(format string, a ...any) evento {
	e.Mensaje = fmt.Sprintf(format, a...)
	return e
}

// Obtiene el mensaje para el usuari que describe al evento.
func (e evento) GetMensaje(format string, a ...any) string {
	return e.Mensaje
}

// ================================================================ //
// ========== Event Store ========================================= //

type EventStoreRepo interface {
	SaveEvent(ev evento) error
}

type EventStoreService struct {
	repo EventStoreRepo
}

// Registra un evento en el EventStore.
func (s *EventStoreService) Rise(ev evento) error {
	op := Op("EventStoreService.Rise")
	if s.repo == nil {
		return op.E(ErrNoDisponible).Str("servicio no instanciado correctamente")
	}
	val := reflect.ValueOf(ev)
	if val.Kind() != reflect.Struct {
		return op.Strf("evento debe ser struct, no %v", val.Kind().String())
	}
	if ev.EventKey == "" {
		return op.Strf("evento %v sin EventKey", val.Kind().String())
	}
	id, err := newID(16)
	if err != nil {
		LogError(err)
	}
	ev.EventID = id
	ev.Timestamp = time.Now()
	err = s.repo.SaveEvent(ev)
	if err != nil {
		return op.Err(err)
	}
	return nil
}
