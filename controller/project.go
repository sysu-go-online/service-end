package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sysu-go-online/service-end/model"
)

// ProjectController is controller for user
type ProjectController struct {
	model.Project
	model.User
	Name        string `json:"name"`
	Description string `json:"description"`
	Language    int    `json:"language"`
	GitPath     string `json:"git_path"`
}

// CreateProjectHandler create project
func CreateProjectHandler(w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	project := ProjectController{}
	if err := json.Unmarshal(body, &project); err != nil {
		return err
	}
	// TODO: check contend
	project.Project.Language = project.Language
	project.Project.Name = project.Name
	project.Project.Description = project.Description
	project.Project.GitPath = project.GitPath
	session := MysqlEngine.NewSession()
	// TODO: get user with jwt
	affected, err := project.Project.Insert(session)
	if err != nil {
		session.Rollback()
		return err
	}
	err = project.Project.CreateProjectRoot(project.User.Username)
	if err != nil {
		session.Rollback()
		return err
	}
	session.Commit()
	if affected == 0 {
		w.WriteHeader(400)
		return nil
	}
	return nil
}

func ListProjectsHandler(w http.ResponseWriter, r *http.Request) error {
	project := ProjectController{}
	session := MysqlEngine.NewSession()
	ps, err := project.Project.GetWithUserID(session)
	if err != nil {
		session.Rollback()
		return err
	}
	if len(ps) == 0 {
		w.WriteHeader(204)
		return nil
	}
	// TODO: construct return json message
	return nil
}
