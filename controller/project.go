package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

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

// ListProjectsResponse is response for list projects
type ListProjectsResponse struct {
	Name     string `json:"name"`
	Language int    `json:"language"`
}

// CreateProjectHandler create project
// TODO: Check if the same name exists
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
	// TODO: check content
	project.Project.Language = project.Language
	project.Project.Name = project.Name
	project.Project.Description = project.Description
	project.Project.GitPath = project.GitPath

	session := MysqlEngine.NewSession()
	project.User.Username = mux.Vars(r)["username"]
	has, err := project.User.GetWithUsername(session)
	if err != nil {
		session.Rollback()
		return err
	}
	if !has {
		w.WriteHeader(401)
		return nil
	}
	project.Project.UserID = project.User.ID
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

// ListProjectsHandler list projects
func ListProjectsHandler(w http.ResponseWriter, r *http.Request) error {
	project := ProjectController{}
	session := MysqlEngine.NewSession()
	project.User.Username = mux.Vars(r)["username"]
	has, err := project.User.GetWithUsername(session)
	if err != nil {
		session.Rollback()
		return err
	}
	if !has {
		w.WriteHeader(401)
		return nil
	}
	project.Project.UserID = project.User.ID
	ps, err := project.Project.GetWithUserID(session)
	if err != nil {
		session.Rollback()
		return err
	}
	if len(ps) == 0 {
		w.WriteHeader(204)
		return nil
	}

	ret := make([]ListProjectsResponse, 0)
	for _, v := range ps {
		tmp := ListProjectsResponse{v.Name, v.Language}
		ret = append(ret, tmp)
	}
	body, err := json.Marshal(ret)
	if err != nil {
		return err
	}
	w.Write(body)
	return nil
}
