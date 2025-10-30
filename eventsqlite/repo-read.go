package eventsqlite

import (
	"database/sql"
	"strings"
	"time"

	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/sqlitedb"
)

func (*EventRepoSqlite) NuevoRepoRead(db sqlitedb.Ejecutor) *RepoRead {
	return &RepoRead{
		db: db,
	}
}

type RepoRead struct {
	db sqlitedb.Ejecutor
}

//  ================================================================  //
//  ========== CONSTANTES ==========================================  //

// Lista de columnas separadas por coma para usar en consulta SELECT
// en conjunto con scanRow o scanRows, ya que las columnas coinciden
// con los campos escaneados.
//
//	event_id,
//	event_key,
//	fecha,
//	data,
//	metadata
const columnasEvento string = "event_id, responsable_id, event_key, fecha, data, metadata"
const ColumnasEventoPrefix string = "ev.event_id, ev.responsable_id, ev.event_key, ev.fecha, ev.data, ev.metadata"

// Origen de los datos de gko.RawEventRow
//
//	FROM eventos
const fromEvento string = "FROM eventos "
const FromEventoPrefix string = "FROM eventos ev "

//  ================================================================  //
//  ========== SCAN ================================================  //

// ScanRowsEvento escanea cada row en la struct Evento
// y devuelve un slice con todos los items.
// Siempre se encarga de llamar rows.Close()
func ScanRowsEvento(rows *sql.Rows, op string) ([]gko.RawEventRow, error) {
	defer rows.Close()
	items := []gko.RawEventRow{}
	for rows.Next() {
		ev := gko.RawEventRow{}
		var fecha string
		err := rows.Scan(
			&ev.EventID, &ev.EventKey, &fecha, &ev.Data, &ev.Metadata,
		)
		if err != nil {
			return nil, gko.ErrInesperado.Err(err).Op(op)
		}
		ev.Fecha, err = time.Parse(formatoTimestamp, fecha)
		if err != nil {
			gko.ErrInesperado.Str("fecha no tiene formato correcto en db").Op("scanRowEvento").Err(err).Log()
		}
		items = append(items, ev)
	}
	return items, nil
}

//  ================================================================  //
//  ========== LIST ================================================  //

func (s *RepoRead) ListEventos() ([]gko.RawEventRow, error) {
	const op string = "ListEventos"
	rows, err := s.db.Query(
		"SELECT " + columnasEvento + " " + fromEvento,
	)
	if err != nil {
		return nil, gko.ErrInesperado.Err(err).Op(op)
	}
	return ScanRowsEvento(rows, op)
}

func (s *RepoRead) ListEventosByID(ids []uint) ([]gko.RawEventRow, error) {
	const op string = "ListEventosByID"
	if len(ids) == 0 {
		return []gko.RawEventRow{}, nil
	}
	placeholders := strings.Repeat("?,", len(ids))
	placeholders = strings.TrimSuffix(placeholders, ",")
	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	rows, err := s.db.Query("SELECT "+columnasEvento+" "+fromEvento+
		"WHERE event_id IN ("+placeholders+")", args...)
	if err != nil {
		return nil, gko.ErrInesperado.Err(err).Op(op)
	}
	return ScanRowsEvento(rows, op)
}

func (s *RepoRead) ListEventosByKey(eventKey gko.EventKey) ([]gko.RawEventRow, error) {
	const op string = "ListEventosByKey"
	rows, err := s.db.Query(
		"SELECT "+columnasEvento+" "+fromEvento+
			"WHERE event_key = ?", eventKey,
	)
	if err != nil {
		return nil, gko.ErrInesperado.Err(err).Op(op)
	}
	return ScanRowsEvento(rows, op)
}

func (s *RepoRead) ListEventosByKeys(keys ...gko.EventKey) ([]gko.RawEventRow, error) {
	const op string = "ListEventosByKeys"
	if len(keys) == 0 {
		return []gko.RawEventRow{}, nil
	}
	placeholders := strings.Repeat("?,", len(keys))
	placeholders = strings.TrimSuffix(placeholders, ",")
	args := make([]any, len(keys))
	for i, key := range keys {
		args[i] = key
	}
	rows, err := s.db.Query("SELECT "+columnasEvento+" "+fromEvento+
		"WHERE event_key IN ("+placeholders+")", args...)
	if err != nil {
		return nil, gko.ErrInesperado.Err(err).Op(op)
	}
	return ScanRowsEvento(rows, op)
}

func (s *RepoRead) ListLastEventos(n int) ([]gko.RawEventRow, error) {
	const op string = "ListLastEventos"
	if n == 0 {
		return nil, gko.ErrDatoInvalido.Str("no limit set")
	}
	if n > 1000 {
		return nil, gko.ErrDatoInvalido.Strf("limit to high: %v", n)
	}
	rows, err := s.db.Query(
		"SELECT "+columnasEvento+" "+fromEvento+
			"ORDER BY fecha DESC LIMIT ?", n,
	)
	if err != nil {
		return nil, gko.ErrInesperado.Err(err).Op(op)
	}
	return ScanRowsEvento(rows, op)
}

// func (s *RepoRead) ListEventosByEntidadID(entidadID int) ([]gko.RawEventRow, error) {
// 	const op string = "ListEventosByEntidadID"
// 	rows, err := s.db.Query(
// 		"SELECT "+ColumnasEventoPrefix+" "+FromEventoPrefix+
// 			"JOIN eventos_nodos ent ON ent.event_id = ev.event_id "+
// 			"WHERE ent.nodo_id = ?", entidadID,
// 	)
// 	if err != nil {
// 		return nil, gko.ErrInesperado.Err(err).Op(op)
// 	}
// 	return ScanRowsEvento(rows, op)
// }
