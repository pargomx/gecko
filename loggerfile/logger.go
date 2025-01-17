package loggerfile

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/pargomx/gecko/gko"
)

// Logger struct to handle logging
//
//	func Example() {
//		logger, err := NewLogger("app.log", 5*time.Second)
//		if err != nil {
//			fmt.Printf("Error initializing logger: %v\n", err)
//			return
//		}
//		defer logger.Close()
//		logger.LogAsync("This is a log entry")
//		logger.LogSync("Another log entry 1")
//		logger.LogSync("Another log entry 2")
//		logger.LogSync("Another log entry 3")
//		time.Sleep(10 * time.Second) // Give some time for logs to process and flush
//	 }
type Logger struct {
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
func NewLogger(filePath string, flushFreq time.Duration) (*Logger, error) {
	op := gko.Op("loggerfile.NewLogger").Ctx("file", filePath).Ctx("flushFreq", flushFreq)

	// Crear directorio si no existe.
	_, err := os.Stat(path.Dir(filePath))
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Creado directorio para", path.Dir(filePath))
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

	logger := &Logger{
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
func (l *Logger) LogSync(entry string) {
	l.mu.Lock()
	_, err := l.writer.WriteString(entry + "\n")
	if err != nil {
		fmt.Printf("Error writing log entry: %v\n", err)
	}
	l.mu.Unlock()
}

// Log adds a new log entry to the channel synchronously.
func (l *Logger) LogBytes(entry []byte) {
	l.mu.Lock()
	_, err := l.writer.Write(entry)
	l.writer.WriteString("\n")
	if err != nil {
		fmt.Printf("Error writing log entry: %v\n", err)
	}
	l.mu.Unlock()
}

// Log adds a new log entry to the channel asynchronously.
func (l *Logger) LogAsync(entry string) {
	l.ch <- entry
}

// processEntries handles log entries asynchronously
func (l *Logger) processEntries() {
	for entry := range l.ch {
		l.mu.Lock()
		_, err := l.writer.WriteString(entry + "\n")
		if err != nil {
			fmt.Printf("Error writing log entry: %v\n", err)
		}
		l.mu.Unlock()
	}
}

// periodicFlush flushes the writer at regular intervals
func (l *Logger) periodicFlush() {
	ticker := time.NewTicker(l.flushFreq)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if l.writer.Buffered() > 0 {
				l.mu.Lock()
				err := l.writer.Flush()
				if err != nil {
					fmt.Printf("Error flushing log writer: %v\n", err)
				}
				l.mu.Unlock()
			}
		case <-l.done:
			return
		}
	}
}

// Close flushes and closes the file logger.
func (l *Logger) Close() {
	close(l.ch)
	close(l.done)
	l.mu.Lock()
	defer l.mu.Unlock()
	err := l.writer.Flush()
	if err != nil {
		gko.Err(err).Op("loggerfile.Close").Ctx("logger", l.filepath)
	}
	err = l.file.Close()
	if err != nil {
		gko.Err(err).Op("loggerfile.Close").Ctx("logger", l.filepath)
	}
}
