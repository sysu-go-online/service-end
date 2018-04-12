package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/unrolled/render"

	"github.com/sysu-go-online/service-end/controller/service"
	"github.com/sysu-go-online/service-end/model/entities"
	dao "github.com/sysu-go-online/service-end/model/service"
)

func CreateProjects(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// Get user token and judge if the token is valid
	token := r.Header.Get("token")
	ok, userid := service.ValidateToken(token)

	if ok {
		// Get post body with json format
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		r.Body.Close()
		var project Project
		err = json.Unmarshal(body, &project)
		if err != nil {
			panic(err)
		}
		// Judge if the userid match token
		if userid != project.UserID {
			w.WriteHeader(401)
			return
		}
		// Insert into db
		dbProject := NewProject(project)
		affected, err := dao.InsertProject(dbProject)
		if err != nil {
			panic(err)
		}
		if affected == 0 {
			// Can not add into db
			w.WriteHeader(400)
		}
	} else {
		// Unauthorized user
		w.WriteHeader(401)
	}
}

func GetProjectsID(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		// Get user token and judge if the token is valid

		userid := r.Form.Get("userID")
		projects, err := dao.FindProjectByUserID(userid)

		if err != nil {
			// can't use db
			w.WriteHeader(400)
		} else {
			formatter.JSON(w, http.StatusOK, struct {
				Content []entities.Project `json:"content"`
			}{
				Content: projects})
		}
	}
}

func TestGetProjectsID(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		// Get user token and judge if the token is valid
		fmt.Println(r.Form.Get("userID"))

		projects := []entities.Project{}
		projects = append(projects, entities.Project{
			ProjectName: "p1",
			ProjectID:   1,
			UserID:      2,
			Language:    0,
		})

		projects = append(projects, entities.Project{
			ProjectName: "p2",
			ProjectID:   2,
			UserID:      2,
			Language:    1,
		})

		projects = append(projects, entities.Project{
			ProjectName: "p3",
			ProjectID:   3,
			UserID:      2,
			Language:    2,
		})

		formatter.JSON(w, http.StatusOK, struct {
			Content []entities.Project `json:"content"`
		}{
			Content: projects,
		})
	}
}
