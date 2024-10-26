package database

import (
	"database/sql"
	"github.com/EularGauss/bandlab-assignment/internal/app/models"
	"strings"
)

var db_conn *sql.DB = nil

func Connect(db_name string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", db_name)
	if db_conn
	db_conn = db
	return db, err
}

func Migrate(db *sql.DB) error {
	for _, model := range models.RegisteredModels {
		createTableSQL := `CREATE TABLE IF NOT EXISTS ` + model.TableName() + ` (`
		for _, field := range model.Fields() {
			createTableSQL += field.Name + ` ` + field.Type + `,`
		}
		createTableSQL = strings.TrimRight(createTableSQL, ",")
		createTableSQL += `);`
		if _, err := db.Exec(createTableSQL); err != nil {
			return err
		}
	}
	return nil
}

func GetDB() *sql.DB {
	return db
}
