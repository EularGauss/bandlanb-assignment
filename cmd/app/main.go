package main

import (
	"database/sql"
	"fmt"
	"github.com/EularGauss/bandlab-assignment/internal/app/database"
	"github.com/EularGauss/bandlab-assignment/internal/app/handlers"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func setup_db() (*sql.DB, error) {
	db, err := database.Connect("test.db")
	if err != nil {
		log.Println("Failed to connect to the database:", err)
		return nil, err
	}
	err = database.Migrate(db)
	if err != nil {
		log.Println("Failed to migrate database:", err)
	}
	return db, err
}

func main() {
	// Connect to SQLite database (it will be created if it doesn't exist)
	_, err := setup_db()
	if err != nil {
		log.Fatal("Error connecting to database")
	}

	r := mux.NewRouter()

	// Define a simple hello world handler
	r.HandleFunc("/post", handlers.CreatePost).Methods("Post")

	// Start the server on port 8080
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
