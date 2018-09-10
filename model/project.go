package model

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/go-xorm/xorm"
)

// Project corresponds to project table in db
type Project struct {
	ID          int        `xorm:"pk autoincr 'id'"`
	Language    int        `xorm:"notnull"`
	UserID      int        `xorm:"'user_id'"`
	Name        string     `xorm:"notnull"`
	CreateTime  *time.Time `xorm:"created"`
	Description string
	GitPath     string
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
func (p *Project) CreateProjectRoot(username string) error {
	userHome := path.Join("/home", username)
	switch p.Language {
	case 0:
		// golang
		projectPath := path.Join(userHome, "/go/src/github.com/", p.Name)
		importPath := path.Join(userHome, "/go/import")
		err := os.MkdirAll(projectPath, os.ModeDir)
		if err != nil {
			return err
		}
		err = os.MkdirAll(importPath, os.ModeDir)
		return err
	case 1:
		// cpp
		projectPath := path.Join(userHome, "/cpp/", p.Name)
		err := os.MkdirAll(projectPath, os.ModeDir)
		return err
	default:
		return errors.New("No such language type")
	}
}

// GetWithUserID returns projects with given user id
func (p *Project) GetWithUserID(session *xorm.Session) ([]Project, error) {
	var ps []Project
	err := session.Where("user_id = ?", p.UserID).Find(&ps)
	if err != nil {
		return nil, err
	}
	return ps, nil
}

// GetWithUserIDAndNmae returns project with given user id and project name
func (p *Project) GetWithUserIDAndName(session *xorm.Session) (bool, error) {
	return session.Where("user_id = ?", p.UserID).And("name = ?", p.Name).Get(p)
}

// GetWithID returns project with given project id
func (p *Project) GetWithID(session *xorm.Session) (bool, error) {
	return session.Where("id = ?", p.ID).Get(p)
}
