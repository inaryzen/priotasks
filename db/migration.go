package db

import (
	"fmt"
	"time"

	"github.com/inaryzen/priotasks/common"
)

const (
	MIGRATION_TABLE_NAME = "migration"
)

func (d *DbSQLite) initMigration() {
	common.Debug("initMigration")
	var err error
	_, err = d.instance.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %v (
			id TEXT PRIMARY KEY,
			time TEXT
		)
	`, MIGRATION_TABLE_NAME))
	if err != nil {
		panic(err)
	}
}

func (d *DbSQLite) MigrationExists(id string) bool {
	rows, err := d.instance.Query("select * from "+MIGRATION_TABLE_NAME+" where id = ?", id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var result = rows.Next()
	common.Debug("migrationExist: id=%v, result=%v", id, result)
	return result
}

func (d *DbSQLite) RecordMigration(id string) {
	sql := "insert into " + MIGRATION_TABLE_NAME + " (id, time) values (?, ?)"
	args := []interface{}{id, time.Now()}
	common.Debug("recordMigration: sql: %v", sql)
	common.Debug("recordMigration: args: %v", args)

	_, err := d.instance.Exec(sql, args...)
	if err != nil {
		err := fmt.Sprintf("failed: %v: %v", id, err)
		panic(err)
	}
	common.Debug("RecordMigration: id=%v", id)
}
