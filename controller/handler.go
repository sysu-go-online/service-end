package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

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

func UpdateFile(w http.ResponseWriter, r *http.Request) {
	// Read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
		w.WriteHeader(500)
		return
	}

	// Read project id and file path from uri
	vars := mux.Vars(r)
	projectID := vars["projectid"]
	filePath := vars["filepath"]

	// Check if the file path is valid
	ok := service.CheckFilePath(filePath)

	if ok {
		// Save file
		dao.UpdateFileContent(projectID, filePath, string(body))
		w.WriteHeader(200)
	} else {
		w.WriteHeader(400)
	}
}

func GetFileContent(w http.ResponseWriter, r *http.Request) {
	// Read project id and file path from uri
	vars := mux.Vars(r)
	projectID := vars["projectid"]
	filePath := vars["filepath"]

	// Check if the file path is valid
	ok := service.CheckFilePath(filePath)
	if ok {
		// Load file
		content := dao.GetFileContent(projectID, filePath)
		w.WriteHeader(200)
		w.Write(content)
	} else {
		w.WriteHeader(400)
	}
}
