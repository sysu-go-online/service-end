package entities

// FileStructure defines the structure of file
type FileStructure struct {
	ID         int             `json:"id"`
	Name       string          `json:"name"`
	Type       string          `json:"type"`
	Children   []FileStructure `json:"children"`
	Root       bool            `json:"root"`
	IsSelected bool            `json:"isSelected"`
}

type UserInfo struct {
	Name  string
	Icon  string
	Email string
}
