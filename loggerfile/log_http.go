package loggerfile

import (
	"encoding/json"

	"github.com/pargomx/gecko"
)

// Implementación simple de log http como JSON al stdout.
func (l *Logger) InsertLogEntry(entr gecko.LogEntry) error {
	logBytes, err := json.Marshal(entr)
	if err != nil {
		return err
	}
	l.mu.Lock()
	_, err = l.writer.Write(logBytes)
	if err != nil {
		return err
	}
	_, err = l.writer.WriteString("\n")
	if err != nil {
		return err
	}
	l.mu.Unlock()
	return nil
}
