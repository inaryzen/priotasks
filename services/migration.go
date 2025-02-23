package services

import (
	"github.com/inaryzen/priotasks/db"
)

func Init() {
	migrationTaskValue()
}

func migrationTaskValue() {
	d := db.DB()
	id := "update_task_value"
	if !d.MigrationExists(id) {
		tasks, err := d.Tasks()
		if err != nil {
			panic(err)
		}
		for _, t := range tasks {
			SaveTask(t)
		}
		d.RecordMigration(id)
	}
}
