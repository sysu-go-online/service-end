package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sysu-go-online/service-end/controller/service"
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
		err = json.Unmarshal(body, project)
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
		dao.InsertProject(dbProject)
	} else {
		// Unauthorized user
		w.WriteHeader(401)
	}
}
