package dbs

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

func SetConnectionToDbSqlServer(
	uri string,
	max_connection int,
	idle_connection int,
	life_connection time.Duration,
	idle_life_connection time.Duration,
) (*sql.DB, error) {

	db_ope, err_db_open := sql.Open("sqlserver", uri)
	if err_db_open != nil {

		return nil, fmt.Errorf("No se ha podido abrir la conexion a la bases de datos de postgres: %q", err_db_open.Error())
	}

	db_ope.SetMaxOpenConns(max_connection * int(time.Second))
	db_ope.SetMaxIdleConns(idle_connection * int(time.Second))
	db_ope.SetConnMaxIdleTime(life_connection)
	db_ope.SetConnMaxLifetime(idle_life_connection)

	return db_ope, nil
}
