package database

import (
	"database/sql"
	"github.com/EularGauss/bandlab-assignment/internal/app/models"
	"log"
	"strings"
)

func Connect(db_name string) *sql.DB {
	db, err := sql.Open("sqlite3", db_name)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	return db
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
