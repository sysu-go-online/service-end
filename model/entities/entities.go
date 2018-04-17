package entities

// FileStructure defines the structure of file
type FileStructure struct {
	ID         int             `json:"id"`
	Name       string          `json:"name"`
	EditStatus int          `json:"edit_status"`
	Type       string          `json:"type"`
	Children   []FileStructure `json:"children"`
}
