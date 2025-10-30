package eventsqlite

import (
	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/sqlitedb"
)

func (*EventRepoSqlite) NuevoRepoWrite(db sqlitedb.Ejecutor) *RepoWrite {
	return &RepoWrite{
		db: db,
	}
}

type RepoWrite struct {
	db sqlitedb.Ejecutor
}

func (s *RepoWrite) Guardar(ev gko.RawEventRow) error {
	const op string = "InsertEvento"
	if ev.EventID == 0 {
		return gko.ErrDatoIndef.Str("pk_indefinida").Op(op).Msg("EventID sin especificar")
	}
	if ev.EventKey == "" {
		return gko.ErrDatoIndef.Str("required_sin_valor").Op(op).Msg("EventKey sin especificar")
	}
	if ev.Fecha.IsZero() {
		return gko.ErrDatoIndef.Str("required_sin_valor").Op(op).Msg("Fecha sin especificar")
	}
	_, err := s.db.Exec("INSERT INTO eventos "+
		"(event_id, responsable_id, event_key, fecha, data, metadata) "+
		"VALUES (?, ?, ?, ?, ?, ?) ",
		ev.EventID, ev.ResponsableID, ev.EventKey, ev.Fecha.Format(formatoTimestamp), ev.Data, ev.Metadata,
	)
	if err != nil {
		return gko.ErrAlEscribir.Err(err).Op(op)
	}
	return nil
}
