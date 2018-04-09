package entities

type Project struct {
	ProjectName string
	ProjectID   int `xorm:"'project_id' pk autoincr"`
	UserID      int `xorm:"user_id notnull"`
	Language    int
}

func (p *Project) TableName() string {
	return "project"
}
