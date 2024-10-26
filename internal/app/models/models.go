package models

import (
	"time"
)

type Model interface {
	TableName() string
	Fields() []FieldInfo
}

// FieldInfo struct for holding information about model fields
type FieldInfo struct {
	Name string
	Type string
}

type Post struct {
	ID        string    `json:"id" db:"id"`
	Caption   string    `json:"caption" db:"caption"`
	ImageURL  string    `json:"imageUrl" db:"image_url"`
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
		{Name: "created_at", Type: "TIMESTAMP"},
		{Name: "updated_at", Type: "TIMESTAMP"},
	}
}

type Comment struct {
	ID        string    `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"`
	Creator   string    `json:"creator" db:"creator"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

func (Comment) TableName() string {
	return "comments"
}

func (Comment) Fields() []FieldInfo {
	return []FieldInfo{
		{Name: "id", Type: "TEXT PRIMARY KEY"},
		{Name: "content", Type: "TEXT"},
		{Name: "creator", Type: "TEXT"},
		{Name: "created_at", Type: "TIMESTAMP"},
	}
}

// Register models in a slice
var RegisteredModels = []Model{
	Post{},
	Comment{},
}
