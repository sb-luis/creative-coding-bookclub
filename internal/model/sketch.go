package model

// MemberSketchInfo holds information about a member and their sketches.
type MemberSketchInfo struct {
	Name     string       // Member's directory name (alias)
	Sketches []SketchInfo // List of sketches for this member
}

// SketchInfo holds information about a single sketch.
// It includes fields populated from file system context (Name, URL, Alias)
// and fields populated from sketch metadata JSON (Title, Description, etc.).
type SketchInfo struct {
	Name  string `json:"-"`
	URL   string `json:"-"` // URL to the sketch page
	Alias string `json:"-"` // Member's alias

	// Fields from metadata JSON
	Title       *string  `json:"title,omitempty"`
	Description *string  `json:"description,omitempty"`
	Keywords    *string  `json:"keywords,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}
