package entities

type Project struct {
	ProjectName string `json:"projectName"`
	ProjectID   int    `xorm:"'project_id' pk autoincr" json:"projectID"`
	UserID      int    `xorm:"user_id notnull" json:"userID"`
	Language    int    `json:"language"`
}

type FileStructure struct {
	ID         int             `json:"id"`
	Name       string          `json:"name"`
	EditStatus int          `json:"edit_status"`
	Type       string          `json:"type"`
	Children   []FileStructure `json:"children"`
}

func (p *Project) TableName() string {
	return "project"
}
