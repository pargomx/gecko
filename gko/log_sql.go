package gko

import (
	"fmt"
	"strings"
)

func PrintableQuery(qry string, args ...any) string {
	numP := strings.Count(qry, "?") // Número de placeholders "?"
	numA := len(args)               // Número de argumentos
	argStr := cCyan + "%v" + reset
	if numP == numA {
		qry = strings.Replace(qry, "?", argStr, numA)
	} else {
		qry = fmt.Sprintf(qry+cRed+" PH:%v ARGS:%v "+reset+"{", numP, numA)
		argStr += reset + ", " + cCyan
		qry = qry + strings.Repeat(argStr, numA)
		qry = strings.TrimSuffix(qry, ", "+cCyan)
		qry += reset + "}"
	}
	return fmt.Sprintf(cBlue+"[QUERY] "+reset+qry, args...)
}
