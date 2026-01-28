package tasks

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type DataVentasForMigrate struct {
	Numero_factura  string
	Cliente         string
	Nombre_cliente  string
	Codigo_producto string
	Codigo_barra    string
	Descripcion     string
	Numero_pedido   string
	Cantidad        int32
	Cesta           string
	Bulto           string
}

func ProducerDataVentasForMigrate(
	connect_db *sql.DB,
	cha_send_data chan DataVentasForMigrate,
	q string,
) {

	defer connect_db.Close()
	defer close(cha_send_data)

	data, err_feth_data := connect_db.Query(q)
	if err_feth_data != nil {

		fmt.Println("No se ha podido obtener los datos que migraran, ", err_feth_data.Error())
		close(cha_send_data)
		return
	}
	defer data.Close()

	var data_ *DataVentasForMigrate = &DataVentasForMigrate{}
	for data.Next() {

		data.Scan(
			&data_.Numero_factura,
			&data_.Cliente,
			&data_.Nombre_cliente,
			&data_.Codigo_producto,
			&data_.Codigo_barra,
			&data_.Descripcion,
			&data_.Numero_pedido,
			&data_.Cantidad,
			&data_.Cesta,
			&data_.Bulto,
		)
		cha_send_data <- *data_
	}
}

func RecibeDataVentasForMigrate(db_neon *pgx.Conn, r *[]DataVentasForMigrate) {

	rows := *r
	copy_count, err_copy_from := db_neon.CopyFrom(
		context.Background(),
		pgx.Identifier{"<NOMBRE DE LA TABLA>"},
		pgx.Identifier{"<CAMPOS DE LAS TABLAS>"},
		pgx.CopyFromSlice(len(rows), func(i int) ([]any, error) {

			return []any{
				rows[i].Numero_factura,
				rows[i].Cliente,
				rows[i].Nombre_cliente,
				rows[i].Codigo_producto,
				rows[i].Codigo_barra,
				rows[i].Descripcion,
				rows[i].Numero_pedido,
				rows[i].Cantidad,
				rows[i].Cesta,
				rows[i].Bulto,
			}, nil
		}),
	)
	if err_copy_from != nil {

		fmt.Println("Ha fallado el copyfrom, ", err_copy_from.Error())
		return
	}

	fmt.Printf("Insert ventas.....\n")
	fmt.Printf("Filas afectadas ventas registradas #%d\n", copy_count)
}
