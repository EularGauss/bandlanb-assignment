package services

import (
	"bandlab-assignment/models"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"sync"
	"time"
)

var (
	posts      = []models.Post{}
	postsMutex = &sync.Mutex{}
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse the request body
	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate the input (add image URL and other logic as necessary)
	if post.Caption == "" {
		http.Error(w, "Caption is required", http.StatusBadRequest)
		return
	}

	post.ID = generateID() // Function to generate unique ID
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	postsMutex.Lock()
	posts = append(posts, post)
	postsMutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	postsMutex.Lock()
	defer postsMutex.Unlock()

	json.NewEncoder(w).Encode(posts)
}

func generateID() string {
	id := uuid.New()
	postIDCounter++
	return fmt.Sprintf("post-%d", id)
}
