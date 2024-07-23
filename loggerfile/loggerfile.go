package loggerfile

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pargomx/gecko/gko"
)

// fileLogger es un servicio que escribe asincrónicamente las entradas a un archivo.
// Cada número determinado de llamadas se llama escribir en disco. Se utiliza
// para no saturar de llamadas a la base de datos o a la función de WriteFile.
type fileLogger struct {

	// printToConsole hace que los mensajes se muestren en la consola además de
	// guardarse en un archivo.
	printToConsole bool

	archivo  *os.File
	filename string

	// Buffer guarda las líneas de log en memoria mientras
	// se alcanza el bufferSize o se manda vaciar explícitamente.
	// Al vaciarlo se escribe su contenido en disco.
	buffer []string

	// Número de líneas de log que se mantendrán en memoria antes de mandarse
	// guardar en disco.
	//
	// Podría usar buffio, pero la implementación es básicamente igual, solo que
	// en vez de mantener un buffer de x bytes, aquí hay uno de x número de
	// líneas de Log. Es además más claro cuántas lineas de log se mantienen en
	// memoria y por ende cuántas se pueden perder si el servidor se apaga
	// abruptamente antes de escribir todos los contenidos del buffer.
	bufferSize int

	// Canal por el que se puede recibir hasta X número de mensajes de manera
	// concurrente mediante invocaciones a las funciones de log. Vale la pena
	// terminar rápido la función que reciba mensajes por el canal, ya que la
	// función que lo envía se bloquea.
	queque chan string
}

// NuevoLogger retorna un servicio que permite enviar entradas de log a un
// archivo de manera concurrente, y guardarlas cuando superen el tamaño del
// buffer dado. Debe deferirse el método Cerrar.
//
// Un bufferSize de 0 asigna el default de 10 entradas.
// Un bufferSize de 1 guarda cada entrada en disco, lo que efectivamente es no usar buffer.
func NuevoLogger(logFilePath string, buffSize int, enConsola bool) (*fileLogger, error) {
	if buffSize < 1 {
		buffSize = 10
	}
	if logFilePath == "" {
		logFilePath = "default.log"
	}
	_, err := os.Stat(logFilePath)
	if err != nil && !os.IsNotExist(err) {
		return nil, gko.Err(err).Op("loggerfile.NuevoLogger").Str("no se puede acceder al archivo de log")
	}
	s := fileLogger{
		filename:       logFilePath,
		bufferSize:     buffSize,
		buffer:         make([]string, 0, buffSize),
		queque:         make(chan string, 50),
		printToConsole: enConsola,
	}
	s.archivo, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		return nil, gko.Err(err).Op("loggerfile.NuevoLogger")
	}
	go s.handleLogQueque()
	return &s, nil
}

// Println recibe una línea de texto y la añade
// directamente al buffer sin pasar por el canal.
func (s *fileLogger) Println(logEntry string) {
	logEntry = time.Now().Format("2006-01-02 15:04:05") + " " + logEntry
	if s.printToConsole {
		fmt.Println(logEntry)
	}
	s.buffer = append(s.buffer, logEntry)
	if len(s.buffer) >= s.bufferSize {
		s.Flush()
	}
}

// handleLogQueque escucha el canal del servicio
// y guarda en el buffer los mensajes que llegan.
// En caso de que el buffer se llene, manda escribir
// en disco todos los contenidos del buffer.
//
// Añade la fecha.
func (s *fileLogger) handleLogQueque() {
	for logEntry := range s.queque {
		logEntry := time.Now().Format("2006-01-02 15:04:05") + " " + logEntry
		if s.printToConsole {
			fmt.Println(logEntry)
		}
		s.buffer = append(s.buffer, logEntry)
		if len(s.buffer) >= s.bufferSize {
			s.Flush()
		}
	}
	s.Flush() //Por si el canal termina.
}

// Flush escribe en disco el contenido del buffer.
func (s *fileLogger) Flush() {
	var porEscribir string // un solo I/O.
	for _, logLine := range s.buffer {
		porEscribir = porEscribir + strings.TrimSuffix(logLine, "\n") + "\n" // siempre \n al final.
	}
	_, err := s.archivo.Write([]byte(porEscribir))
	if err != nil {
		// A consola para evitar loops potenciales.
		fmt.Printf("Log(%v): %s\n", s.filename, err)
		fmt.Println(porEscribir)
	}
	// Vaciar el buffer. El colector de basura desaloja la memoria.
	s.buffer = nil
}

// Cerrar debe llamarse siempre antes de salir del programa
// para asegurarse de que se guardan las últimas entradas de
// log que aún están en el buffer.
// Añade el mensaje dado al último.
func (s *fileLogger) Cerrar() {
	// Vaciar buffer al disco para no perder información.
	s.Flush()
	// El error de close es la última oportunidad para que el sistema operativo
	// reporte una falla en la escritura de los datos.
	// Si el servidor se reinicia abruptamente, no pasa nada con los archivos abiertos
	// que no hayan sido cerrados, ya que la Kernel se encarga de los file descriptors:
	// "When a process terminates, the kernel releases the resources owned
	// by the process and notifies the child's parent of its demise."
	// No importa de qué manera termine la aplicación, siempre se cerrarán por el sistema:
	// "All of the file descriptors, directory streams, conversion descriptors,
	// and message catalog descriptors open in the calling process shall be closed."
	if err := s.archivo.Close(); err != nil {
		fmt.Printf("Log(%v): %s\n", s.filename, err)
		return
	}
}
