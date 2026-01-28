package executetask

import (
	"context"
	"fmt"
	"pg_mss_migrate_data/dbs"
	"pg_mss_migrate_data/tasks"
	"time"
)

func MigrateComprasForTask() {

	fmt.Println("Migracion de compras activa")

	defer fmt.Println("Migracion de compras culminada")

	rows_compras := []tasks.DataComprasForMigrate{}
	var cha_sent_data chan tasks.DataComprasForMigrate = make(chan tasks.DataComprasForMigrate, 1000)

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

	go tasks.ProducerDataComprasForMigrate(
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

	fmt.Println("Agregando compras ......")
	for v := range cha_sent_data {
		rows_compras = append(rows_compras, v)
	}

	tasks.RecibeDataComprasForMigrate(db_neon, &rows_compras)
}
