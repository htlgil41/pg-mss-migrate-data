package dbs

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func SetConnectionToDbPostgres(uri string) (*pgx.Conn, error) {

	db_ope, err_db_open := pgx.Connect(context.Background(), uri)
	if err_db_open != nil {

		return nil, fmt.Errorf("No se ha podido abrir la conexion a la bases de datos de postgres: %q", err_db_open.Error())
	}
	return db_ope, nil
}
