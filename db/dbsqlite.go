package db

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	"github.com/inaryzen/priotasks/common"
	_ "modernc.org/sqlite"
)

type DbSQLite struct {
	instance *sql.DB
}

func NewDbSQLite() *DbSQLite {
	return &DbSQLite{}
}

func (d *DbSQLite) Init() {
	dir, err := common.ResolveAppDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}
	file := filepath.Join(dir, "db.sqlite")

	db, err := sql.Open("sqlite", file)
	if err != nil {
		log.Fatal(err)
	}
	d.instance = db

	d.initMigration()
	d.initTasks()
	d.initSettings()
}

func (d *DbSQLite) columnExists(tableName, columnName string) bool {
	query := fmt.Sprintf("PRAGMA table_info(%s);", tableName)
	rows, err := d.instance.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name, ctype, notnull string
		var dfltValue interface{}
		var primaryKey int

		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dfltValue, &primaryKey); err != nil {
			panic(err)
		}

		if name == columnName {
			return true
		}
	}
	return false
}

func (d *DbSQLite) Close() {
	common.Debug("closing db...")
	d.instance.Close()
}
