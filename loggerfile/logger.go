package loggerfile

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
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
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		return nil, err
	}

	logger := &Logger{
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

// Close closes the logger and the underlying file
func (l *Logger) Close() error {
	close(l.ch)
	close(l.done)
	l.mu.Lock()
	defer l.mu.Unlock()
	err := l.writer.Flush()
	if err != nil {
		return err
	}
	return l.file.Close()
}
