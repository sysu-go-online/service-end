package entities

type Project struct {
	ProjectName string `json:"projectName"`
	ProjectID   int    `xorm:"'project_id' pk autoincr" json:"projectID"`
	UserID      int    `xorm:"user_id notnull" json:"userID"`
	Language    int    `json:"language"`
}

func (p *Project) TableName() string {
	return "project"
}
