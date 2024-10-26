package main

import (
	"fmt"
	"github.com/EularGauss/bandlab-assignment/internal/app/database"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main() {
	// Connect to SQLite database (it will be created if it doesn't exist)
	db := database.Connect("test.db")
	err := database.Migrate(db)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	defer db.Close()
	fmt.Print("Database migrated successfully")
}
