package loggerfile

import (
	"github.com/pargomx/gecko/gko"
)

func (s *fileLogger) PrintSQL(qry string, args ...any) {
	s.queque <- gko.PrintableQuery(qry, args...)
}
