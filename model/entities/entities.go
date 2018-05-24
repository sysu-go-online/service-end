package entities

import "time"

// FileStructure defines the structure of file
type FileStructure struct {
	ID         int             `json:"id"`
	Name       string          `json:"name"`
	Type       string          `json:"type"`
	Children   []FileStructure `json:"children"`
	Root       bool            `json:"root"`
	IsSelected bool            `json:"isSelected"`
}

// UserInfo maps users table in the db
type UserInfo struct {
	Name       string `xorm:"'username' pk"`
	Icon       string
	Email      string
	CreateTime *time.Time
	Gender     int
	Age        int
	Token      string // Stores token from github
}

func (u UserInfo) TableName() string {
	return "users"
}
