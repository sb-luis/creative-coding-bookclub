package model

import (
	"time"
)

// MemberSketchInfo holds information about a member and their sketches.
type MemberSketchInfo struct {
	Name     string       // Member's directory name (alias)
	Sketches []SketchInfo // List of sketches for this member
}

// SketchInfo holds information about a single sketch.
// It includes fields populated from file system context (Slug, URL, Alias)
// and fields populated from sketch metadata JSON (Title, Description, etc.).
type SketchInfo struct {
	Slug  string `json:"-"`
	URL   string `json:"-"` // URL to the sketch page
	Alias string `json:"-"` // Member's alias

	// Fields from metadata JSON
	Title       *string  `json:"title,omitempty"`
	Description *string  `json:"description,omitempty"`
	Keywords    *string  `json:"keywords,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// Sketch represents a sketch stored in the database
type Sketch struct {
	ID               int       `json:"id" db:"id"`
	MemberID         int       `json:"member_id" db:"member_id"`
	Slug             string    `json:"slug" db:"slug"`
	Title            string    `json:"title" db:"title"`
	Description      string    `json:"description" db:"description"`
	Keywords         string    `json:"keywords" db:"keywords"`
	Tags             []string  `json:"tags" db:"-"`          // Will be stored as JSON
	TagsJSON         string    `json:"-" db:"tags"`          // JSON string for database
	ExternalLibs     []string  `json:"external_libs" db:"-"` // Will be stored as JSON
	ExternalLibsJSON string    `json:"-" db:"external_libs"` // JSON string for database
	SourceCode       string    `json:"source_code" db:"source_code"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// CreateSketchRequest represents the data needed to create a new sketch
type CreateSketchRequest struct {
	Title        string     `json:"title" validate:"required,min=1,max=200"`
	Description  string     `json:"description" validate:"max=1000"`
	Keywords     string     `json:"keywords" validate:"max=500"`
	Tags         []string   `json:"tags" validate:"dive,min=1,max=50"`
	ExternalLibs []string   `json:"external_libs" validate:"dive,min=1,max=100"`
	SourceCode   string     `json:"source_code" validate:"required,min=1,max=1000000"` // 1MB max for UTF-8
	CreatedAt    *time.Time `json:"created_at,omitempty"`                              // Optional
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`                              // Optional
}

// UpdateSketchRequest represents the data that can be updated for an existing sketch
type UpdateSketchRequest struct {
	Title        *string  `json:"title,omitempty" validate:"omitempty,min=1,max=200"`
	Description  *string  `json:"description,omitempty" validate:"omitempty,max=1000"`
	Keywords     *string  `json:"keywords,omitempty" validate:"omitempty,max=500"`
	Tags         []string `json:"tags,omitempty" validate:"dive,min=1,max=50"`
	ExternalLibs []string `json:"external_libs,omitempty" validate:"dive,min=1,max=100"`
	SourceCode   *string  `json:"source_code,omitempty" validate:"omitempty,min=1,max=1000000"` // 1MB max for UTF-8
}
