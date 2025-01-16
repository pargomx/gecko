package gecko

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pargomx/gecko/gko"
)

// LogEntry corresponde a un elemento de la tabla 'loghttp'.
type LogEntry struct {
	Timestamp    time.Time     // `loghttp.timestamp`
	Host         string        // `loghttp.host`
	Method       string        // `loghttp.method`  Verbo HTTP de la solicitud
	Ruta         string        // `loghttp.ruta`  Patrón del router sobre la que se enrutó
	URI          string        // `loghttp.uri`  Path con query params
	Htmx         bool          // `loghttp.htmx`  Es una solicitud HTMX
	Status       int           // `loghttp.status`  Código HTTP de la respuesta
	Latency      time.Duration // `loghttp.latency`  Tiempo de procesamiento y envío de la respuesta
	BytesIn      uint64        // `loghttp.bytes_in`  Request "Content-Length"
	BytesOut     uint64        // `loghttp.bytes_out`  Bytes escritos como respuesta
	Error        error         // `loghttp.error`  En caso de haber error
	RemoteIP     string        // `loghttp.remote_ip`
	Sesion       string        // `loghttp.sesion`
	UserAgent    string        // `loghttp.user_agent`
	Referer      string        // `loghttp.referer`
	HxCurrentURL string        // `loghttp.hx_current_url`
	HxTarget     string        // `loghttp.hx_target`
	HxTrigger    string        // `loghttp.hx_trigger`
	HxBoosted    bool          // `loghttp.hx_boosted`
}

type HTTPLogger interface {
	InsertLogEntry(entr LogEntry) error
}

func (g *Gecko) logHTTP(c *Context, err error) {
	bytesIn, _ := TxtUint64(c.request.Header.Get("Content-Length"))
	origin := c.request.Header.Get("Origin")
	logEnt := LogEntry{
		Timestamp:    c.time,
		RemoteIP:     c.RealIP(),
		UserAgent:    c.request.UserAgent(),
		Host:         c.request.Host,
		Status:       c.response.Status,
		Method:       c.request.Method,
		Ruta:         c.path,
		URI:          c.request.RequestURI,
		Error:        err,
		Latency:      time.Since(c.time),
		BytesIn:      bytesIn,
		BytesOut:     c.response.Size,
		Htmx:         c.EsHTMX(),
		Referer:      strings.TrimPrefix(c.request.Referer(), origin),
		HxCurrentURL: strings.TrimPrefix(c.request.Header.Get("Hx-Current-Url"), origin),
		HxTarget:     c.request.Header.Get("Hx-Target"),
		HxTrigger:    c.request.Header.Get("Hx-Trigger"),
		HxBoosted:    c.request.Header.Get("Hx-Boosted") == "true",
	}
	if len(c.SesionID) > 6 {
		logEnt.Sesion = c.SesionID[:6] // Conocer usuario sin exponer sesión.
	}
	logErr := g.HTTPLogger.InsertLogEntry(logEnt)
	if logErr != nil {
		gko.Err(logErr).Op("LogHTTP").Log()
	}
}

// ================================================================ //
// ================================================================ //

type HTTPLoggerJSON struct{}

// Implementación simple de log http como JSON al stdout.
func (l *HTTPLoggerJSON) InsertLogEntry(entr LogEntry) error {
	log, err := json.Marshal(entr)
	if err != nil {
		return err
	}
	fmt.Println(string(log))
	return nil
}
