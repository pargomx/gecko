package gko

import (
	"reflect"
	"time"
)

// ================================================================ //
// ========== Instrucción ========================================= //

// Instruction represents the user's intent to do a command that the
// application needs to confirm, warn or inform the user before commiting
// the acction. The server must present those messages to the user
// and remember the instructions while waiting for the confirmation.
type Instruction struct {
	id   string
	time time.Time

	argsStruct any

	// Para dar información al usuario sobre lo que sucederá con su comando.
	consecuencias []evento

	// Para dar información al usuario sobre lo que no sucederá.
	advertencias []*Error
}

// const (
// 	tError int = iota
// 	tWarning
// 	tInfo
// 	tDebug
// )

// Para advertir al usuario sobre lo que sucederá (o no) con su comando.
// type advertencia struct {
// key   ErrorKey
// msg   string // mensaje para el usuario.
// level int    // warning, error
// }

// ================================================================ //
// ========== Servicio ============================================ //

// Servicio para ejecutar comandos luego de confirmar con el usuario.
//
// Las instrucciones sin confirmar se eliminan a los 20 minutos.
//
// Instanciar con NewSavedInstrService()
type SavedInstructionsService struct {
	pendientes map[string]*Instruction
	ticker     *time.Ticker
}

func NewSavedInstrService() *SavedInstructionsService {
	cs := &SavedInstructionsService{
		pendientes: make(map[string]*Instruction),
		ticker:     time.NewTicker(5 * time.Minute),
	}
	go func() {
		for range cs.ticker.C {
			if len(cs.pendientes) > 0 {
				now := time.Now()
				for id, cmd := range cs.pendientes {
					if now.Sub(cmd.time) > 20*time.Minute {
						LogDebugf("Command expired: %v", cmd.id)
						delete(cs.pendientes, id)
					}
				}
			}
		}
	}()
	return cs
}

var debugActive = false

func (s *SavedInstructionsService) ToggleDebug() {
	debugActive = !debugActive
}

// ================================================================ //
// ========== Life Cycle ========================================== //

// Guarda una instrucción para ejecutar un comando después de confirmar con el
// usuario.
//
// El argsStruct debe ser una estructura con los argumentos que el comando
// necesita para ejecutarse. Por ejemplo:
//
//	type ChangeAgeArgs struct {
//	    UserID string
//	    NewAge int
//	}
//
// Después de registrar la instrucción se deben agregar los mensajes necesarios
// según las validaciones que se hagan en la capa de aplicación. Por ejemplo:
//
//	cmd, _ := s.ComandosPendientes.Add(args)
//	if args.NewAge < 10 {
//	    cmd.AddWarning(AgeToSmall, "Too young: %v", args.NewAge)
//	}
func (s *SavedInstructionsService) Add(argsStruct any) (*Instruction, error) {
	op := Op("SavedInstructions.Add")
	if s.pendientes == nil || s.ticker == nil {
		return nil, op.E(ErrNoDisponible).Str("servicio no instanciado correctamente")
	}
	if len(s.pendientes) > 999 {
		return nil, op.E(ErrTooManyReq).Msg("Cola de comandos llena, esperar por favor")
	}
	if argsStruct == nil {
		return nil, op.Str("argsStruct no puede ser nil")
	}
	val := reflect.ValueOf(argsStruct)
	if val.Kind() != reflect.Struct {
		return nil, op.Strf("argsStruct debe ser struct, no %v", val.Kind().String())
	}
	id, err := newID(16)
	if err != nil {
		LogError(err)
	}
	if _, alreadyExists := s.pendientes[id]; alreadyExists {
		return nil, op.E(ErrInesperado).Str("uuid repetido para nueva instrucción")
	}
	cmd := &Instruction{
		id:         id,
		time:       time.Now(),
		argsStruct: argsStruct,
	}
	s.pendientes[id] = cmd
	return cmd, nil
}

// Devuelve una copia de la instrucción especificada para que se tome acción al
// respecto y la remueve la original del servicio, asumiendo que ya no será
// necesaria. Si desea modificar o conservar la instrucción debe reconstruir
// otra nueva desde cero.
//
// Utilizar ArgsStruct() para obtener a qué comando corresponde y sus parámetros.
func (s *SavedInstructionsService) Remove(id string) (*Instruction, error) {
	op := Op("SavedInstructions.Remove")
	if s.pendientes == nil || s.ticker == nil {
		return nil, op.E(ErrNoDisponible).Str("servicio no instanciado correctamente")
	}
	if id == "" {
		return nil, op.E(ErrDatoIndef).Str("id vacío")
	}
	cmd, ok := s.pendientes[id]
	if !ok {
		return nil, op.E(ErrNoEncontrado).Str("comando '%v' no encontrado")
	}
	copia := *cmd
	delete(s.pendientes, cmd.id)
	return &copia, nil
}

// ================================================================ //
// ========== API - Setters ======================================= //

// Used to construct the URL to confirm the command.
func (c *Instruction) CmdID() string {
	return c.id
}

// Para dar información al usuario sobre lo que sucederá con el comando.
// si se ejecuta con los argumentos proporcionados.
func (c *Instruction) AddEffect(ev evento) *Instruction {
	c.consecuencias = append(c.consecuencias, ev)
	return c
}

// Para dar información al desarrollador sobre lo que sucederá con el comando.
func (c *Instruction) AddDebug(ev evento) *Instruction {
	if debugActive {
		c.consecuencias = append(c.consecuencias, ev)
	}
	return c
}

// Para advertir al usuario sobre lo que sucederá (o no) con su comando.
func (c *Instruction) AddWarning(err *Error) *Instruction {
	c.advertencias = append(c.advertencias, err)
	return c
}

// Para advertir al usuario sobre lo que sucederá (o no) con su comando.
func (c *Instruction) AddError(err *Error) *Instruction {
	c.advertencias = append(c.advertencias, err)
	return c
}

// ================================================================ //
// ========== API - Getters ======================================= //

// Devuelve la estructura con argumentos originalmente proporcionada
// al crear la instrucción. Hacer type casting para averiguar a qué
// comando corresponde.
//
// Ejemplos:
//
//	args := cmd.ArgsStruct().(myapp.DoSomethingArgs)
//
//	switch args := cmd.ArgsStruct().(type) {
//		case myapp.DoSomethingArgs:
//	}
func (c *Instruction) ArgsStruct() any {
	return c.argsStruct
}

// ================================================================ //

// HasEffect reporta si la instrucción tiene tal advertencia o error.
func (c *Instruction) HasWarning(key ErrorKey) bool {
	for _, err := range c.advertencias {
		if err.Contiene(key) {
			return true
		}
	}
	return false
}

// HasEffect reporta si la instrucción tiene la consecuencia especificada.
func (c *Instruction) HasEffect(key EventKey) bool {
	for _, ev := range c.consecuencias {
		if ev.EventKey == key {
			return true
		}
	}
	return false
}

func (c *Instruction) HasMsg(key string) bool {
	for _, err := range c.advertencias {
		if err.Contiene(ErrorKey(key)) {
			return true
		}
	}
	for _, ev := range c.consecuencias {
		if ev.EventKey == EventKey(key) {
			return true
		}
	}
	return false
}

// Devuelve todos los mensajes tipo Info y Debug (si debugActive).
//
// Usar HasEffect(key) para saber si hay alguna en específico.
func (c *Instruction) Effects() []evento {
	return c.consecuencias
}

// Devuelve todos los mensajes tipo Warning y Error.
//
// Usar HasWarning(key) para saber si hay alguna en específico.
func (c *Instruction) Warnings() []*Error {
	return c.advertencias
}
