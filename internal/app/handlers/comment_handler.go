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

func deleteCommentFromDB(commentID, postID string) error {
	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM comments WHERE id = ? AND postId = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(commentID, postID)
	return err
}

func AddCommentToPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	postId := mux.Vars(r)["postId"]
	var commentInput Comment

	if err := json.NewDecoder(r.Body).Decode(&commentInput); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if commentInput.Content == "" {
		respondWithError(w, http.StatusBadRequest, "Comment content not provided")
		return
	}

	commentInput.ID = generateID()

	if err := storeCommentToDB(postId, commentInput); err != nil {
		http.Error(w, "Failed to store the comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(commentInput)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	commentID := mux.Vars(r)["commentId"]
	postID := mux.Vars(r)["postId"]

	err := deleteCommentFromDB(commentID, postID)
	if err != nil {
		http.Error(w, "Failed to delete the comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
