package datasources

import (
	"database/sql"
	"strconv"

	graphql "github.com/graph-gophers/graphql-go"
)

type Resolver struct {
	DB *sql.DB
}

// This is nutty, but converting from an int32 to a string requires it to first
// be converted to an int. This library expects a graphql.ID, which a string can
// be converted to. This method is benchmark faster than `fmt.Sprint(i)`
func Int32ToGraphqlID(i int32) graphql.ID {
	return graphql.ID(strconv.Itoa(int(i)))
}
