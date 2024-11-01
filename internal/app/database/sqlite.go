package database

import (
	"database/sql"
	"fmt"
	"github.com/EularGauss/bandlab-assignment/internal/app/database/models"
	"strings"
	"sync"
)

var db_conn *sql.DB = nil
var mutex = &sync.Mutex{}

func Connect(db_name string) (*sql.DB, error) {
	mutex.Lock()
	defer mutex.Unlock()
	db, err := sql.Open("sqlite3", db_name)
	if db_conn == nil {
		db_conn = db
	}
	return db, err
}

// Migrate creates tables in the database
func Migrate(db *sql.DB) error {
	for _, model := range models.RegisteredModels {

		tableName := model.TableName()
		fields := model.Fields()

		var fieldDefinitions []string
		for _, field := range fields {
			fieldDefinitions = append(fieldDefinitions, fmt.Sprintf("%s %s", field.Name, field.Type))
		}

		constraints := model.Constraints()
		var constraintsDefinition string
		if len(constraints) > 0 {
			constraintsDefinition = strings.Join(constraints, ", ")
		}

		// TODO check if the tables exists already, if so, skip creating the table
		// Reason being that if there are multiple services running, they might try to create the table at the same time
		var sql string
		if constraintsDefinition != "" {
			sql = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s, %s);", tableName, strings.Join(fieldDefinitions, ", "), constraintsDefinition)
		} else {
			sql = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", tableName, strings.Join(fieldDefinitions, ", "))
		}

		_, err := db.Exec(sql)
		if err != nil {
			return fmt.Errorf("failed to create table %s: %w", tableName, err)
		}
	}
	return nil
}

// TODO add connection pool library
func GetDB() *sql.DB {
	db, err := Connect("test.db")
	if err != nil {
		fmt.Errorf("Failed to connect to the database: %v", err)
		return nil
	}
	return db
}
