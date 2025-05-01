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

func (d *DbSQLite) Init(dbFile string) {
	common.Debug("dbsql init...")
	if dbFile == "" {
		dir, err := common.ResolveAppDir()
		if err != nil {
			log.Fatalf("Failed to get home directory: %v", err)
		}
		dbFile = filepath.Join(dir, "db.sqlite")
	}
	common.Debug("Init: dbFile=%v", dbFile)

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	d.instance = db
	_, err = d.instance.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal(err)
	}

	d.initMigration()
	d.initTasks()
	d.initSettings()
	d.initTags()
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
