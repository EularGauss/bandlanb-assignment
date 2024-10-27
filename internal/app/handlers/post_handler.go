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

	// Iterate through the rows
	for rows.Next() {
		var post PostWithComments
		var lastComments string // Temporary variable to hold JSON result of last comments

		if err := rows.Scan(&post.Caption, &post.ImageURL, &lastComments); err != nil {
			return nil, "", err
		}

		// Unmarshal the JSON string of last comments into the Comments slice
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

func GetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Read query parameters for pagination
	limitStr := r.URL.Query().Get("limit")   // Number of posts to retrieve
	cursorStr := r.URL.Query().Get("cursor") // Cursor for pagination

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
