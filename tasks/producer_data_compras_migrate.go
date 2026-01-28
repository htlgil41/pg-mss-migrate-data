package tasks

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type DataComprasForMigrate struct {
	Numero_factura    string
	Numero_pedido     string
	Fecha_factura     string
	Descripcion       string
	Proveedor         string
	Tipo_cuenta       int32
	Codigo_producto   string
	Codigo_lote       string
	Fecha_vencimiento string
	Cantidad_factura  float64
	Peso_facturado    float64
	Unidad_empaque    float64
	Costo_compra      float64
}

func ProducerDataComprasForMigrate(
	connect_db *sql.DB,
	cha_send_data chan DataComprasForMigrate,
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

	var data_ *DataComprasForMigrate = &DataComprasForMigrate{}
	for data.Next() {

		data.Scan(
			&data_.Numero_factura,
			&data_.Numero_pedido,
			&data_.Fecha_factura,
			&data_.Descripcion,
			&data_.Proveedor,
			&data_.Tipo_cuenta,
			&data_.Codigo_producto,
			&data_.Codigo_lote,
			&data_.Fecha_vencimiento,
			&data_.Cantidad_factura,
			&data_.Peso_facturado,
			&data_.Unidad_empaque,
			&data_.Costo_compra,
		)
		cha_send_data <- *data_
	}
}

func RecibeDataComprasForMigrate(db_neon *pgx.Conn, r *[]DataComprasForMigrate) {

	rows := *r
	copy_count, err_copy_from := db_neon.CopyFrom(
		context.Background(),
		pgx.Identifier{"<NOMBRE DE LA TABLA>"},
		pgx.Identifier{"<CAMPOS DE LAS TABLAS>"},
		pgx.CopyFromSlice(len(rows), func(i int) ([]any, error) {

			return []any{
				rows[i].Numero_factura,
				rows[i].Numero_pedido,
				rows[i].Fecha_factura,
				rows[i].Descripcion,
				rows[i].Proveedor,
				rows[i].Tipo_cuenta,
				rows[i].Codigo_producto,
				rows[i].Codigo_lote,
				rows[i].Fecha_vencimiento,
				rows[i].Cantidad_factura,
				rows[i].Peso_facturado,
				rows[i].Unidad_empaque,
				rows[i].Costo_compra,
			}, nil
		}),
	)
	if err_copy_from != nil {

		fmt.Println("Ha fallado el copyfrom, ", err_copy_from.Error())
		return
	}

	fmt.Printf("Insert compras.....\n")
	fmt.Printf("Filas afectadas compras registradas #%d\n", copy_count)
}
