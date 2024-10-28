package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/EularGauss/bandlab-assignment/internal/app/database"
	"github.com/EularGauss/bandlab-assignment/internal/app/database/models"
	"github.com/gorilla/mux"
	"net/http"
)

type CommentInput struct {
	PostId  string `json:"postId"`
	Content string `json:"content"`
	Creator string `json:"creator"`
}

func storeCommentToDB(comment *CommentInput, userId string) error {
	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	stmt, err := db.Prepare("INSERT INTO comments (id, content, creator, postId) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(generateID(), comment.Content, comment.Creator, comment.PostId)
	return err
}

func deleteCommentFromDB(commentID, postID string) error {
	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

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

	var commentInput CommentInput

	if err := json.NewDecoder(r.Body).Decode(&commentInput); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if commentInput.Content == "" {
		respondWithError(w, http.StatusBadRequest, "Comment content not provided")
		return
	}

	userId := "annonymous" // getch userid from jwt token in production

	if err := storeCommentToDB(&commentInput, userId); err != nil {
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
