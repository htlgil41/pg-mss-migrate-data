![Logo](https://intus.com.mx/wp-content/uploads/2022/08/Migracion-de-datos-ERP-CRM-Base-de-Datos.png)

# Historia PG-MSS-MIGRATE-DATA 

En la empresa surgió la necesidad de migrar más de 400,000 registros desde SQL Server hacia PostgreSQL en la nube. 

El verdadero reto fue ejecutar el proceso de la forma más rápida y eficiente posible, optimizando tanto el consumo de recursos de las bases de datos como el rendimiento de la aplicación, logrando así una lectura y escritura de datos simultánea sin degradar el sistema.

# Descripcion de la app/script

Se escribio un programa en Go donde se levanta un worker en un hilo diferente del principal gracias a una goRoutine donde este buscaba los datos a migrar en la db SQLserver y los enviaba a travez de una canal con un buffer definido que se leian en el hilo principal de la ejecucion del programa para asi llamar a una funcion a la que se le pasaba la conexion a PostgreSQL y todos los datos almacenados en un slice **(La colecion se hizo mediante la lectura del canal hasta cerrarce)** en dicha funcion gracias a la libreria o paquete/modulo de goland **pgx** se ejcutaba **"Copy"** con todos los registros de manera totalmente eficiente segura **(Sin transacciones por lo menos en la primera version)** pero si algun error ocurria si la transaccion principal simplemente se detenia o fallaba por lo que no se completaba la migacion.


Nota: En la primer version se realizaron las pruebas y migracion total con exito en un cierto tiempo considerable (Mas de 1minuto) ya que la red era lenta.

# Resumen

Esta aplicación ofrece una solución rápida y robusta para la migración de miles o millones de registros entre bases de datos. El código está diseñado para ser altamente adaptable, permitiendo ajustes según la estructura de los datos y los motores involucrados.

Ya sea una migración relacional o una transición de SQL a NoSQL, el sistema proporciona una base sólida para reestructurar la lógica según las necesidades del proyecto. Una herramienta sencilla en su uso, pero potente y eficaz en su ejecución.

# Problemas

- Tratar de leer los registros que se deseaban migrar de manera eficiente para no saturar la db en cuestion ya que era una db en produccion

- Distribuir la carga o la insercion de registros al estar migrando chucks de datos o completos en la db de destino

- Realizar la operacion bajo nivel de ejecucion es decir, que la aplicacion o script fuera lo suficientemente rapido para mover la data en cuestion de segundos.

# Intentos

- Mientras los datos se obtenian se constuia un string con el formato del insert por lo que al finalizar quedaba en memoria un string muy grande con 400K+ de registros en una sola insercion. (Colapsaba tanto PostgreSQL como la aplicacion)

- Crear chuck de estos insert por cada N registros, el error generado era igual que el anterior ya que se codifico para cada 1K registros crear un insert completo y ejecutar multiples. (Las operaciones se realizaban por ceparado y fallaban algunas por lo que la data quedaba inconsistente y con faltas de registros. Se podia manejar una transaccion pero entraba el problema de rendimineto al respecto de la memoria ya que la ejecucion se convertia un cuello de botella al sistema operativo)

- Crear un insert.txt, un archivo con los insert separados terminaba siendo una bomba al culminar el archivo ya que pesaba mucho y era imposible de solo copiar y pegar. (El tiempo de ejecucion de la aplicacion empezo a mejorar pero seguia sin tener un resultado final satisfactorio ya que no daba una fiabilidad que los datos se migraron correctamente a la db de PostgreSQL sino mas bien a un txt)

- Formato para usar **COPY** de PostgreSQL es decir, crear un archivo con la estructura necesaria para hacer la migracion con (Se genero el archivo mas se tenia que hacer la migracion manualmente **"No se realizo por tiempo pero era la mejor opcion"**):

``` sql

    COPY nombre_tabla (columna1, columna2....)
    FROM '/ruta/al/archivo.csv'
    DELIMITER ';'
    CSV HEADER;

```

# Solución

En este punto, lo que cobraba mayor sentido era la metodología de inserción de datos: se debía realizar a través de un **COPY** en lugar de un **INSERT INTO TABLA VALUES ()()()...400k**.

Al iniciar la aplicación, el programa levantaba un Worker/Producer encargado de extraer los registros desde la base de datos SQL Server. Estos datos se transmitían a un canal con un buffer N específico, diseñado para gestionar la latencia y establecer un margen de flujo controlado. Desde el hilo principal, se recibía dicha información para consolidar la colección de datos y ejecutar la función COPY, aprovechando las capacidades avanzadas de la librería **Pgx**. 

Gracias a esta estructura, simple pero poderosa, se logró migrar la información de manera eficiente, incluso enfrentando una latencia de red considerable.