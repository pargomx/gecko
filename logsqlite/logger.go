package logsqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	_ "github.com/glebarez/go-sqlite"
	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
)

type Logger struct {
	dbPath    string
	db        *sql.DB
	entries   []gecko.LogEntry // Buffer for entries
	mu        sync.Mutex       // Buffer mutex
	flushFreq time.Duration
}

var createTableLogHTTP = `
CREATE TABLE loghttp (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  timestamp TEXT NOT NULL,
  host TEXT NOT NULL,
  method TEXT NOT NULL,
  ruta TEXT NOT NULL,
  uri TEXT NOT NULL,
  htmx INT NOT NULL,
  status INT NOT NULL,
  latency INT NOT NULL,
  bytes_in INT NOT NULL,
  bytes_out INT NOT NULL,
  error TEXT NOT NULL,
  remote_ip TEXT NOT NULL,
  sesion TEXT NOT NULL,
  user_agent TEXT NOT NULL,
  referer TEXT NOT NULL,
  hx_current_url TEXT NOT NULL,
  hx_target TEXT NOT NULL,
  hx_trigger TEXT NOT NULL,
  hx_boosted INT NOT NULL
);
`

var pragmaConfig = "?_pragma=foreign_keys(0)&_busy_timeout=1000"

// NewLogger instancia un nuevo logger que escribirá sus entradas
// en un archivo de texto. Escribirá al archivo con la frecuencia
// especificada para no saturar de operaciones IO cuando haya muchas
// solicitudes en poco tiempo.
func NewLogger(dbPath string, flushFreq time.Duration) (*Logger, error) {
	if dbPath == "" {
		return nil, errors.New("log db path no especificada")
	}

	// Crear directorio si no existe.
	_, err := os.Stat(path.Dir(dbPath))
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Creado directorio para", path.Dir(dbPath))
		err := os.MkdirAll(path.Dir(dbPath), 0755)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	// Verificar o crear archivo para base de datos.
	_, err = os.Stat(dbPath)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Creado archivo para base de datos", dbPath)
		err = os.WriteFile(dbPath, []byte{}, 0664)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	// Abrir base de datos.
	db, err := sql.Open("sqlite", dbPath+pragmaConfig)
	if err != nil {
		return nil, err
	}
	// Para evitar error database locked. https://github.com/mattn/go-sqlite3/issues/274
	db.SetMaxOpenConns(1)

	// Verificar o crear tabla para logs http.
	var tblExists bool
	err = db.QueryRow("SELECT count(name) FROM sqlite_master WHERE type='table' AND name='loghttp'").Scan(&tblExists)
	if !tblExists {
		if err != nil {
			gko.LogError(err) // TODO: ErrNoRows?
		}
		gko.LogEventof("Inicializando log http en sqlite")
		_, err = db.Exec(createTableLogHTTP)
		if err != nil {
			return nil, err
		}
	}

	// Instanciar logger.
	logger := &Logger{
		dbPath:    dbPath,
		db:        db,
		flushFreq: flushFreq,
		entries:   make([]gecko.LogEntry, 0, 8),
	}
	go logger.flushBufferToDB()
	go logger.periodicFlush()
	return logger, nil
}

// Close closes the logger and the underlying db
func (l *Logger) Close() error {
	l.flushBufferToDB()
	return l.db.Close()
}

// ================================================================ //

// Implementación de log http en sqlite.
func (l *Logger) InsertLogEntry(entry gecko.LogEntry) error {
	l.mu.Lock()
	l.entries = append(l.entries, entry)
	l.mu.Unlock()
	return nil
}

// Saves all entries in buffer at regular intervals.
func (l *Logger) periodicFlush() {
	if l.flushFreq < time.Second {
		l.flushFreq = time.Second // Mínimo 1s de intervalo.
	}
	ticker := time.NewTicker(l.flushFreq)
	defer ticker.Stop()
	for range ticker.C {
		l.flushBufferToDB()
	}
}

// Inserts to db all entries in buffer
func (l *Logger) flushBufferToDB() {
	if len(l.entries) == 0 {
		return
	}
	l.mu.Lock()
	tx, err := l.db.Begin() // Varios inserts en una sola transacción es más eficiente.
	if err != nil {
		fmt.Printf("Error begining log transaction: %v\n", err)
	}
	for _, entry := range l.entries {
		err := insertLogEntry(tx, entry)
		if err != nil {
			fmt.Printf("Error inserting log entry: %v\n", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		fmt.Printf("Error commiting log transaction: %v\n", err)
	}
	l.entries = l.entries[:0] // Don't dispose of the underlying array.
	l.mu.Unlock()
}

// Inserts one entry to the sqlite db.
func insertLogEntry(tx *sql.Tx, entr gecko.LogEntry) error {
	const op string = "InsertLogEntry"
	if entr.Timestamp.IsZero() {
		return gko.ErrDatoIndef().Op(op).Msg("Timestamp sin especificar").Str("pk_indefinida")
	}
	_, err := tx.Exec("INSERT INTO loghttp "+
		"(timestamp, host, method, ruta, uri, htmx, status, latency, bytes_in, bytes_out, error, remote_ip, sesion, user_agent, referer, hx_current_url, hx_target, hx_trigger, hx_boosted) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ",
		entr.Timestamp, entr.Host, entr.Method, entr.Ruta, entr.URI, entr.Htmx, entr.Status, entr.Latency, entr.BytesIn, entr.BytesOut, entr.Error, entr.RemoteIP, entr.Sesion, entr.UserAgent, entr.Referer, entr.HxCurrentURL, entr.HxTarget, entr.HxTrigger, entr.HxBoosted,
	)
	if err != nil {
		return gko.ErrAlEscribir().Err(err).Op(op)
	}
	return nil
}
