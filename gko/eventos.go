package gko

import (
	"encoding/json"
	"math/rand/v2"
	"reflect"
	"time"
)

// ================================================================ //
// ========== Evento ============================================== //

// EventKey identifica a la consecuencia de un comando sin importar sus
// parámetros específicos. Es el tipo de evento. Por ejemplo:
//
//	var (
//		AgeModified    gko.EventKey = "age_modified"
//		UserRegistered gko.EventKey = "user_registered"
//		ProductSold    gko.EventKey = "product_sold"
//	)
type EventKey string

// EventData es una estructura que contiene los argumentos
// involucrados en el evento.
type EventData interface {
	// Mensaje para el usuario describiendo el evento.
	//
	// Pueden implementarse varios de mensajes según se necesite
	// Por ejemplo: "en", "es", "pasado", "futuro", etc.
	//
	// Siempre debe haber un mensaje default.
	ToMsg(tipo string) string
}

// Event representa la consecuencia de un comando y los parámetros
// con los que fue ejecutado.
type Event struct {
	EventID  uint
	EventKey EventKey
	Fecha    time.Time
	Data     EventData // Será serializado
	Metadata []byte
}

// RawEventRow corresponde a un elemento de la tabla 'eventos'.
type RawEventRow struct {
	EventID  uint      // `eventos.event_id`
	EventKey EventKey  // `eventos.event_key`
	Fecha    time.Time // `eventos.fecha`
	Data     []byte    // `eventos.data`
	Metadata []byte    // `eventos.metadata`
}

func (e Event) Mensaje() string {
	if e.Data == nil {
		return string(e.EventKey) + " (no data)"
	}
	return e.Data.ToMsg("")
}

// ================================================================ //
// ========== Event Store ========================================= //

type EventStore struct {
	Repo       EventStoreRepo // Persisitir eventos.
	Results    *TxResult      // Guardar en memoria durante la transacción.
	ConsoleLog bool           // Activar para mostrar mensajes en log.
}

type EventStoreRepo interface {
	Guardar(ev RawEventRow) error
}

// Registra un evento en los lugares configurados (Repo / TxResults / Log)
//
// Key: identificador del evento. Ej. "usuario_registrado".
// Data: estructura con los argumentos del evento que permite
// reproducirlo o mostrar un mensaje.
func (s *EventStore) Rise(key EventKey, data EventData) (*Event, error) {
	op := Op("EventStore.Rise").Msg("Hubo un problema en el servidor")

	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Struct {
		return nil, op.Strf("data debe ser struct, no %v", val.Kind().String())
	}
	if key == "" {
		return nil, op.Strf("evento %v sin EventKey", val.Kind().String())
	}

	// Crear nuevo evento
	ev := Event{
		// Mask porque sqlite no soporta uint64,
		// y para que en sqlite no aparezcan negativos.
		EventID:  rand.Uint() & 0x7FFFFFFFFFFFFFFF,
		EventKey: key,
		Fecha:    time.Now(), // TODO: always local?
		Data:     data,
	}

	// Serializar y guardar en repositorio.
	if s.Repo != nil {
		dataJson, err := json.Marshal(data)
		if err != nil {
			return nil, op.Err(err)
		}
		row := RawEventRow{
			EventID:  ev.EventID,
			EventKey: ev.EventKey,
			Fecha:    ev.Fecha,
			Data:     dataJson,
		}
		err = s.Repo.Guardar(row)
		if err != nil {
			return nil, op.Err(err)
		}
	}

	// Mantener en memoria como resultado de transacción.
	if s.Results != nil {
		s.Results.Events = append(s.Results.Events, ev)
	}

	// Log como fallback o si se configuró.
	if (s.Results == nil && s.Repo == nil) || s.ConsoleLog {
		LogInfof("%v %+v", key, data)
	}

	// Nota: podríamos retornar solo el EventID...
	return &ev, nil
}

// ================================================================ //
// ========== Read events ========================================= //

// Para deserializar un evento se necesita un mapa de EventKey hacia
// el tipo concreto de EventData que permita acceder a los argumentos
// y construir el mensaje particular del evento.
//
// Este mapa hace tal conexión mediante una función que retorna una
// instancia vacía del EventData type concreto.
var eventDataConstructors = make(map[EventKey]func() EventData)

// Para deserializar un evento se necesita que la aplicación registre
// qué EventKey corresponde a qué tipo EventData concreto.
// DeclareEvent registra un EventData type con un EventKey correspondiente.
//
// Un mismo EventData puede usarse para varios EventKey, pero no al revés.
//
// The `sample` parameter is used to infer the concrete type for unmarshaling.
//
// This function should be called during package initialization (e.g., in init() functions).
func DeclareEvent(key EventKey, sample EventData) {
	if key == "" {
		LogAlertf("Eventkey empty for %T. Ignoring.", sample)
		return
	}
	if _, exists := eventDataConstructors[key]; exists {
		LogAlertf("EventKey '%s' already declared. Overwriting.\n", key)
	}
	eventDataConstructors[key] = func() EventData {
		// Use reflection to create a new instance of the same type as 'sample'.
		// This is the most generic way to get an empty instance of an interface type.
		// Make sure 'sample' is a pointer if you want to unmarshal into a pointer.
		return createInstance(sample)
	}
}

// Helper to create a new instance of a type that implements EventData.
// It uses reflection to create a new value of the same type as the sample.
func createInstance(sample EventData) EventData {
	// If the sample is a pointer, create a new instance of the underlying type
	// and return a pointer to it. Typical for JSON unmarshaling into structs.
	val := makeNew(sample)
	return val.Interface().(EventData)
}

// makeNew is a helper for createInstance
func makeNew(sample any) reflect.Value {
	typ := reflect.TypeOf(sample)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return reflect.New(typ)
}

func ParseEvent(row RawEventRow) (*Event, error) {
	concreteEventConstructor, ok := eventDataConstructors[row.EventKey]
	if !ok {
		return nil, ErrAlLeer.Strf("evento desconocido: %v", row.EventKey)
	}
	eventData := concreteEventConstructor()
	err := json.Unmarshal(row.Data, &eventData)
	if err != nil {
		return nil, ErrAlLeer.Strf("can't unmarshall event into %v", row.EventKey)
	}
	// If constructor returns a pointer, we need to dereference it
	// if EventData methods are on value receiver.
	// Or ensure all EventData methods have pointer receivers.
	// For JSON, it's common to unmarshal into a pointer to a struct.
	return &Event{
		EventID:  row.EventID,
		EventKey: row.EventKey,
		Fecha:    row.Fecha,
		Data:     eventData,
	}, nil
}

func ParseEvents(rawEvents []RawEventRow) ([]Event, error) {
	op := Op("ParseEvents")
	events := make([]Event, 0, len(rawEvents))
	for _, rawEv := range rawEvents {
		ev, err := ParseEvent(rawEv)
		if err != nil {
			return nil, op.Err(err).Ctx("key", rawEv.EventKey).Ctx("id", rawEv.EventID)
		}
		events = append(events, *ev)
		// events[i] = *ev
	}
	return events, nil
}
