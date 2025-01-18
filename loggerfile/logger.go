package loggerfile

import (
	"bufio"
	"errors"
	"os"
	"path"
	"sync"
	"time"

	"github.com/pargomx/gecko/gko"
)

type logger struct {
	filepath  string
	file      *os.File
	writer    *bufio.Writer
	mu        sync.Mutex
	ch        chan string
	done      chan struct{}
	flushFreq time.Duration
}

// NewLogger instancia un nuevo logger que escribirá sus entradas
// en un archivo de texto. Escribirá al archivo con la frecuencia
// especificada para no saturar de operaciones IO cuando haya muchas
// solicitudes en poco tiempo.
func NewLogger(filePath string, flushFreq time.Duration) (*logger, error) {
	op := gko.Op("logfile.NewLogger").Ctx("file", filePath).Ctx("flushFreq", flushFreq)
	if filePath == "" {
		return nil, op.ErrDatoIndef().Str("log file path undefined")
	}
	if flushFreq < time.Second { // Mínimo 1s de intervalo.
		flushFreq = time.Second * 5 // Default 5s.
	}

	// Crear directorio si no existe.
	_, err := os.Stat(path.Dir(filePath))
	if errors.Is(err, os.ErrNotExist) {
		gko.LogInfof("Creado directorio %s", path.Dir(filePath))
		err := os.MkdirAll(path.Dir(filePath), 0755)
		if err != nil {
			return nil, op.Err(err)
		}
	} else if err != nil {
		return nil, op.Err(err)
	}

	// Abrir o crear archivo.
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		return nil, op.Err(err)
	}

	logger := &logger{
		filepath:  filePath,
		file:      file,
		writer:    bufio.NewWriter(file),
		ch:        make(chan string, 100), // Buffered channel for log entries
		done:      make(chan struct{}),    // For ticker
		flushFreq: flushFreq,
	}

	go logger.processEntries()
	go logger.periodicFlush()

	return logger, nil
}

// Log adds a new log entry to the channel synchronously.
func (l *logger) LogSync(entry string) {
	l.mu.Lock()
	_, err := l.writer.WriteString(entry + "\n")
	if err != nil {
		gko.Err(err).Op("logfile.LogSync")
	}
	l.mu.Unlock()
}

// Log adds a new log entry to the channel synchronously.
func (l *logger) LogBytes(entry []byte) {
	l.mu.Lock()
	_, err := l.writer.Write(entry)
	l.writer.WriteString("\n")
	if err != nil {
		gko.Err(err).Op("logfile.LogBytes")
	}
	l.mu.Unlock()
}

// Log adds a new log entry to the channel asynchronously.
func (l *logger) LogAsync(entry string) {
	l.ch <- entry
}

// processEntries handles log entries asynchronously
func (l *logger) processEntries() {
	for entry := range l.ch {
		l.mu.Lock()
		_, err := l.writer.WriteString(entry + "\n")
		if err != nil {
			gko.Err(err).Op("logfile.LogAsync")
		}
		l.mu.Unlock()
	}
}

// periodicFlush flushes the writer at regular intervals
func (l *logger) periodicFlush() {
	ticker := time.NewTicker(l.flushFreq)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if l.writer.Buffered() > 0 {
				l.mu.Lock()
				err := l.writer.Flush()
				if err != nil {
					gko.Err(err).Op("logfile.Flush")
				}
				l.mu.Unlock()
			}
		case <-l.done:
			return
		}
	}
}

// Close flushes and closes the file logger.
func (l *logger) Close() {
	close(l.ch)
	close(l.done)
	l.mu.Lock()
	defer l.mu.Unlock()
	err := l.writer.Flush()
	if err != nil {
		gko.Err(err).Op("logfile.Close").Ctx("logger", l.filepath)
	}
	err = l.file.Close()
	if err != nil {
		gko.Err(err).Op("logfile.Close").Ctx("logger", l.filepath)
	}
}
