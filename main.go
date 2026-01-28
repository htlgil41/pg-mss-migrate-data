package main

import executetask "pg_mss_migrate_data/execute_task"

func main() {

	executetask.MigrateVentasForTask()
	executetask.MigrateComprasForTask()
	select {}
}
