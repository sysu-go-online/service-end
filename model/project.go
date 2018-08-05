package model

import (
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
)

// Project corresponds to project table in db
type Project struct {
	ID         int        `xorm:"pk autoincr 'id'"`
	Language   int        `xorm:"notnull"`
	UserID     int        `xorm:"'user_id'"`
	Name       string     `xorm:"notnull"`
	CreateTime *time.Time `xorm:"created"`
}

// TableName defines table name
func (p Project) TableName() string {
	return "project"
}

// Insert insert a project to db
func (p *Project) Insert(session *xorm.Session) (int, error) {
	affected, err := session.Insert(p)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return int(affected), nil
}

// CreateProjectRoot create project root in the user home
// TODO:
func (p *Project) CreateProjectRoot() {}

// GetWithUserID returns project with given user id
func (p *Project) GetWithUserID(session *xorm.Session) {}

// GetWithUserIDAndNmae returns project with given user id and project name
func (p *Project) GetWithUserIDAndNmae(session *xorm.Session) {}

// GetWithID returns project with given project id
func (p *Project) GetWithID(session *xorm.Session) {}
