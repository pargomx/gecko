package loggerfile

import (
	"github.com/pargomx/gecko/gko"
)

func (s *logger) LogQuerySQL(qry string, args ...any) {
	s.LogSync(gko.PrintableQuery(qry, args...))
}
