package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/EularGauss/bandlab-assignment/internal/app/database"
	"github.com/EularGauss/bandlab-assignment/internal/app/database/models"
	"net/http"
	"strconv"
	"time"
)

type PostInput struct {
	Caption   string `json:"caption"`
	LoadImage bool   `json:"image"`
}

type PostOutput struct{
	ID string `json:"id"`
	PreSignedURL string `json:"preSignedURL"`
}

type PostWithComments struct {
	Caption   string    `json:"caption"`
	ImageURL  string    `json:"imageUrl"`
	Comments  []Comment `json:"comments"`
	CreatedAt time.Time `json:"createdAt"`
}

type Comment struct {
	ID      string `json:"id" db:"id"`
	Content string `json:"content" db:"content"`
	Creator string `json:"creator" db:"creator"` // In production, we fetch this from JWT token
}

func storePostToDB(post *models.Post) error {
	db := database.GetDB()
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

func fetchPostsFromDB(limit int, cursor string) ([]PostWithComments, string, error) {
	db := database.GetDB()
	if db == nil {
		return nil, "", fmt.Errorf("database connection is nil")
	}
	defer db.Close()

	var query string
	if cursor != "" {
		query = `
            SELECT p.caption, p.image_url, 
                (SELECT json_group_array(json_object('content', c.content))
                 FROM (SELECT c.content 
                     FROM comments c 
                     WHERE c.post_id = p.id 
                     ORDER BY c.created_at DESC 
                     LIMIT 2)) AS last_comments
            FROM posts p
            WHERE p.created_at < (
                SELECT created_at FROM posts WHERE id = ?
            )
            ORDER BY p.created_at DESC 
            LIMIT ?`
	} else {
		query = `
            SELECT p.caption, p.image_url, 
                (SELECT json_group_array(json_object('content', c.content))
                 FROM (SELECT c.content 
                     FROM comments c 
                     WHERE c.post_id = p.id 
                     ORDER BY c.created_at DESC 
                     LIMIT 2)) AS last_comments
            FROM posts p
            ORDER BY p.created_at DESC 
            LIMIT ?`
	}

	rows, err := db.Query(query, cursor, limit)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var posts []PostWithComments
	var nextCursor string // For cursor-based pagination

	for rows.Next() {
		var post PostWithComments
		var lastComments string

		if err := rows.Scan(&post.Caption, &post.ImageURL, &lastComments); err != nil {
			return nil, "", err
		}

		if lastComments != "" {
			if err := json.Unmarshal([]byte(lastComments), &post.Comments); err != nil {
				return nil, "", err
			}
		}

		posts = append(posts, post)

		// Set nextCursor to the created_at of the last retrieved post for pagination
		nextCursor = post.CreatedAt.Format(time.RFC3339)
	}

	// Return retrieved posts and the next cursor value
	return posts, nextCursor, nil
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var postInput PostInput

	if err := json.NewDecoder(r.Body).Decode(&postInput); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if postInput.Caption == "" {
		respondWithError(w, http.StatusBadRequest, "Caption not provided")
		return
	}

	newPost := models.Post{
		ID:      generateID(),
		Caption: postInput.Caption,
	}

	presignedURL := ""
	if postInput.LoadImage == true {
		s3Service := GetS3Service()
		presignedURL, err := s3Service.GeneratePresignedURL(postInput.ImageKey)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to generate presigned URL")
			return
		}
	}

	if err := storePostToDB(&newPost); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to store post")
		return
	}
	output := PostOutput{ID: newPost.ID, PreSignedURL: presignedURL}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	limitStr := r.URL.Query().Get("limit")
	cursorStr := r.URL.Query().Get("cursor")

	limit := 10 // Default limit if not specified
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	posts, nextCursor, err := fetchPostsFromDB(limit, cursorStr)
	if err != nil {
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"posts":      posts,
		"nextCursor": nextCursor,
	})
}
