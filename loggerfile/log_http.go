package loggerfile

import (
	"encoding/json"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
)

// Implementación simple de log http como JSON al stdout.
func (l *logger) SaveLog(entr gecko.LogEntry) {
	logBytes, err := json.Marshal(entr)
	if err != nil {
		gko.Err(err).Op("loggerfile.SaveLog").Log()
		return
	}
	l.mu.Lock()
	_, err = l.writer.Write(logBytes)
	if err != nil {
		gko.Err(err).Op("loggerfile.SaveLog").Log()
		return
	}
	_, err = l.writer.WriteString("\n")
	if err != nil {
		gko.Err(err).Op("loggerfile.SaveLog").Log()
		return
	}
	l.mu.Unlock()
}
