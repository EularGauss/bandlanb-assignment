package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/EularGauss/bandlab-assignment/internal/app/database"
	"github.com/gorilla/mux"
	"net/http"
)

func storeCommentToDB(postId string, comment Comment) error {
	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	defer db.Close()

	// Prepare the SQL statement for inserting a comment
	stmt, err := db.Prepare("INSERT INTO comments (id, content, creator, postId) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement with the comment data
	_, err = stmt.Exec(comment.ID, comment.Content, comment.Creator, postId)
	return err
}

func AddCommentToPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract postId from the URL
	postId := mux.Vars(r)["postId"] // Assuming you're using Gorilla Mux for routing
	var commentInput Comment

	if err := json.NewDecoder(r.Body).Decode(&commentInput); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if commentInput.Content == "" {
		respondWithError(w, http.StatusBadRequest, "Comment content not provided")
		return
	}

	// Generate unique ID for the comment and set necessary fields
	commentInput.ID = generateID() // Implement this function to generate a unique ID

	if err := storeCommentToDB(postId, commentInput); err != nil {
		http.Error(w, "Failed to store the comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(commentInput)
}
