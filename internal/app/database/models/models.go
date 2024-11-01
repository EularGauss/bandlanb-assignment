// This is a model to show how we can make a connection to db
// This is basic implementation of ORM model. it is intentionally kept simple to show how we can use ORM model
package models
import (
	"time"
)

type Model interface {
	TableName() string
	Fields() []FieldInfo
	Constraints() []string
}


type FieldInfo struct {
	Name string
	Type string
}

type Post struct {
	ID        string    `json:"id" db:"id"`
	Caption   string    `json:"caption" db:"caption"`
	ImageURL  string    `json:"imageUrl" db:"image_url"`
	UserID    string    `json:"userId" db:"user_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

func (Post) TableName() string {
	return "posts"
}

func (Post) Fields() []FieldInfo {
	return []FieldInfo{
		{Name: "id", Type: "TEXT PRIMARY KEY"},
		{Name: "caption", Type: "TEXT"},
		{Name: "image_url", Type: "TEXT"},
		{Name: "created_at", Type: "TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL"},
		{Name: "updated_at", Type: "TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL"},
		{Name: "user_id", Type: "INTEGER"},
	}
}


//TODO Add foreign key constraint for user_id
func (Post) Constraints() []string {
	return []string{}
}

type Comment struct {
	ID        string    `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"`
	Creator   string    `json:"creator" db:"creator"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

func (Comment) TableName() string {
	return "comments"
}

func (Comment) Fields() []FieldInfo {
	return []FieldInfo{
		{Name: "id", Type: "TEXT PRIMARY KEY"},
		{Name: "post_id", Type: "TEXT"},
		{Name: "content", Type: "TEXT"},
		{Name: "creator", Type: "TEXT"},
		{Name: "created_at", Type: "TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL"},
		{Name: "updated_at", Type: "TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL"},
	}
}

func (Comment) Constraints() []string {
	return []string{
		"FOREIGN KEY(post_id) REFERENCES posts(id)",
	}
}

// Register models in a slice
var RegisteredModels = []Model{
	Post{},
	Comment{},
}
