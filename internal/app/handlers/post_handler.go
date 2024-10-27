package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/EularGauss/bandlab-assignment/internal/app/database"
	"github.com/EularGauss/bandlab-assignment/internal/app/database/models"
	"net/http"
)

type PostInput struct {
	Caption   string `json:"caption"`
	LoadImage bool   `json:"image"`
}

type Comment struct {
	ID      string `json:"id" db:"id"`
	Content string `json:"content" db:"content"`
	Creator string `json:"creator" db:"creator"` // In production, we fetch this from JWT token
}

func storePostToDB(post *models.Post) error {
	db := database.GetDB() // Assuming GetDB function is defined elsewhere
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	defer db.Close()

	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO posts (id, caption) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement with the post data
	_, err = stmt.Exec(post.ID, post.Caption)
	return err // Return the error if any
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var postInput PostInput

	// Decode the request body into the PostInput struct
	if err := json.NewDecoder(r.Body).Decode(&postInput); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if postInput.Caption == "" {
		respondWithError(w, http.StatusBadRequest, "Caption not provided")
		return
	}
	// Create a new Post instance
	newPost := models.Post{
		ID:      generateID(), // Implement this function to generate a unique ID
		Caption: postInput.Caption,
	}
	if err := storePostToDB(&newPost); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to store post")
		return
	}
	// Return the created post
	w.WriteHeader(http.StatusCreated)
}

func GetPost() {

}
