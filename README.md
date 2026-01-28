![Logo](https://intus.com.mx/wp-content/uploads/2022/08/Migracion-de-datos-ERP-CRM-Base-de-Datos.png)

# Historia 

En la empresa necesitaban migrar informacion/data a la nube y eran mas de 400K registros, la informacion necesitaba ser obtenida de una bases de datos SQLserver y migrarla a una bases de datos PostgreSQL.

el verdadero reto era hacerlo de la forma mas rapida posible y eficiente tanto para las bases de datos como para la ejecucion de la aplicacion y no tener problemas de rendimiento al leer y escribir datos en un mismo proceso.

# Problemas

- Leer los registros que se deseaban migrar de manera eficiente para no saturar la db en cuestion ya que era una db en produccion

- Distribuir la carga o la insercion de registros al estar migrando chucks de datos o completos en la db de destino

- Realizar la operacion bajo nivel de ejecucion es decir, que la aplicacion o script fuera lo suficientemente rapido para mover la data en cuestion de segundos.

# Intentos

- Mientras los datos se obtenian se constuia un string con el formato del insert por lo que al finalizar quedaba en memoria un string muy grande con 400K+ de registros en una sola insercion. (Colapsaba tanto PostgreSQL como la aplicacion)

- Crear chuck de estos insert por cada N registros, el error generado era igual que el anterior ya que se codifico para cada 1K registros crear un insert completo y ejecutar multiples. (Las operaciones se realizaban por ceparado y fallaban algunas por lo que la data quedaba inconsistente y con faltas de registros. Se podia manejar una transaccion pero entraba el problema de rendimineto al respecto de la memoria ya que la ejecucion se convertia un cuello de botella al sistema operativo)

- Crear un insert.txt, crear un archivo con los insert ceparados terminaba siendo una bomba al culminar el archivo ya que pesaba mucho y era imposible de solo copiar y pegar. (El tiempo de ejecucion de la aplicacion empezo a mejorar pero seguia sin tener un resultado final satisfactorio ya que no daba una fiabilidad que los datos se migraron correctamente a la db de PostgreSQL sino mas bien a un txt)

- Formato para usar **COPY** de PostgreSQL es decir, crear un archivo con la estructura necesaria para hacer la migracion con (Se genero el archivo mas se tenia que hacer la migracion manualmente **"No se realizo por tiempo pero era la mejor opcion"**):

``` sql

    COPY nombre_tabla (columna1, columna2....)
    FROM '/ruta/al/archivo.csv'
    DELIMITER ';'
    CSV HEADER;

```

# Descripcion de la app/script

Se escribio un programa en Go donde se levanta un worker en un hilo diferente del principal gracias a una goRoutine donde este buscaba los datos a migrar en la db SQLserver y los enviaba a travez de una canal con un buffer definido que se leian en el hilo principal de la ejecucion del programa para asi llamar a una funcion a la que se le pasaba la conexion a PostgreSQL y todos los datos almacenados en un slice (La colecion se hizo mediante la lectura del canal hasta cerrarce) en dicha funcion gracias a la libreria o paquete/modulo de goland **pgx** se ejcutaba **"Copy"** con todos los registros de manera totalmente eficiente segura **(Sin transacciones por lo menos en la primera version)** pero si algun error ocurria si la transaccion principal simplemente se detenia o fallaba por lo que no se completaba la migacion.


Nota: En la primer version se realizaron las pruebas y migracion total con exito en un cierto tiempo considerable (Mas de 1minuto) ya que la red era lenta.