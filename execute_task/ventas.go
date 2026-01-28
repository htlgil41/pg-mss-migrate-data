package executetask

import (
	"context"
	"fmt"
	"pg_mss_migrate_data/dbs"
	"pg_mss_migrate_data/tasks"
	"time"
)

func MigrateVentasForTask() {

	fmt.Println("Migracion de ventas activa")

	defer fmt.Println("Migracion de ventas culminada")

	rows_ventas := []tasks.DataVentasForMigrate{}
	var cha_sent_data chan tasks.DataVentasForMigrate = make(chan tasks.DataVentasForMigrate, 1000)

	db_mss, err_db_mss_open := dbs.SetConnectionToDbSqlServer(
		"<URL URI PARA CONEXION SQLSERVER>",
		5,
		5,
		(5 * time.Minute),
		(1 * time.Minute),
	)
	if err_db_mss_open != nil {

		fmt.Println(err_db_mss_open.Error())
		return
	}
	defer db_mss.Close()

	go tasks.ProducerDataVentasForMigrate(
		db_mss,
		cha_sent_data,
		"<QUERY DE LA DATA A MIGRAR>",
	)

	db_neon, err_db_neon_open := dbs.SetConnectionToDbPostgres("<URL URI PARA CONEXION POSTGRESQL>")
	if err_db_neon_open != nil {

		fmt.Println(err_db_neon_open.Error())
		return
	}
	defer db_neon.Close(context.Background())

	fmt.Println("Agregando ventas ......")
	for v := range cha_sent_data {

		rows_ventas = append(rows_ventas, v)
	}

	tasks.RecibeDataVentasForMigrate(db_neon, &rows_ventas)
}
