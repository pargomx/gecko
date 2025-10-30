package eventsqlite

import (
	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/sqlitedb"
)

const formatoTimestamp = "2006-01-02 15:04:05.99999-07:00"

const createTableEventos = `
CREATE TABLE eventos (
  event_id INT NOT NULL,
  responsable_id INT NOT NULL,
  event_key TEXT NOT NULL,
  fecha TEXT NOT NULL,
  data TEXT NOT NULL DEFAULT '',
  metadata TEXT NOT NULL DEFAULT '',
  PRIMARY KEY (event_id)
);
CREATE INDEX index_eventos_key ON eventos (event_key);
CREATE INDEX index_eventos_responsable ON eventos (responsable_id);
CREATE INDEX index_eventos_fecha ON eventos (fecha);
`

type EventRepoSqlite struct{}

func NuevoEventRepoSqlite(db sqlitedb.Ejecutor) (*EventRepoSqlite, error) {
	// Crear tabla eventos si no existe.
	const checkTableExists = "SELECT name FROM sqlite_master WHERE type='table' AND name='eventos';"
	var tableName string
	err := db.QueryRow(checkTableExists).Scan(&tableName)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return nil, gko.Err(err).Op("gko.EventosSqlite.NuevoRepo")
	}
	if tableName == "" {
		_, err := db.Exec(createTableEventos)
		if err != nil {
			return nil, gko.Err(err).Op("gko.EventosSqlite.NuevoRepo")
		}
		gko.LogInfo("EventosSqlite: nueva tabla preparada")
	}

	return &EventRepoSqlite{}, nil
}
