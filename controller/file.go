package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/sysu-go-online/service-end/model"
)

// UpdateFileHandler is a handler for update file or rename
func UpdateFileHandler(w http.ResponseWriter, r *http.Request) error {
	// Read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	type temp struct {
		Operation string `json:"operation"`
		Content   string `json:"content"`
	}

	// get operation from body
	op := &temp{}
	err = json.Unmarshal(body, op)
	if err != nil {
		w.WriteHeader(400)
		return nil
	}

	// get username from jwt
	ok, username := GetUserNameFromToken(r.Header.Get("Authorization"))
	if !ok {
		w.WriteHeader(401)
		return nil
	}
	// Read project name and file path from uri
	vars := mux.Vars(r)
	projectName := vars["projectname"]
	filePath := vars["filepath"]

	// Get project information
	session := MysqlEngine.NewSession()
	u := model.User{Username: username}
	ok, err = u.GetWithUsername(session)
	if !ok {
		w.WriteHeader(400)
		return nil
	}
	if err != nil {
		return err
	}
	p := model.Project{Name: projectName, UserID: u.ID}
	has, err := p.GetWithUserIDAndName(session)
	if !has {
		w.WriteHeader(204)
		return nil
	}
	if err != nil {
		return err
	}

	// Check if the file path is valid
	ok = checkFilePath(filePath)

	if ok {
		switch op.Operation {
		case "update":
			err = model.UpdateFileContent(projectName, username, filePath, op.Content, false, false, p.Language)
		case "rename":
			err = model.RenameFile(projectName, username, filePath, op.Content, p.Language)
		default:
			w.WriteHeader(400)
			return nil
		}
		if err != nil {
			return err
		}
		w.WriteHeader(200)
	} else {
		w.WriteHeader(400)
	}
	return nil
}

// CreateFileHandler is a handler for create file
func CreateFileHandler(w http.ResponseWriter, r *http.Request) error {
	// Read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	type IsDir struct {
		Dir bool `json:"dir"`
	}
	// Judge if it is dir from body
	isDir := IsDir{}
	err = json.Unmarshal(body, &isDir)
	if err != nil {
		return err
	}
	dir := isDir.Dir

	// get username from jwt
	ok, username := GetUserNameFromToken(r.Header.Get("Authorization"))
	if !ok {
		w.WriteHeader(401)
		return nil
	}
	// Read project id and file path from uri
	vars := mux.Vars(r)
	projectName := vars["projectname"]
	filePath := vars["filepath"]

	// Check if the file path is valid
	ok = checkFilePath(filePath)

	// Get project information
	session := MysqlEngine.NewSession()
	u := model.User{Username: username}
	ok, err = u.GetWithUsername(session)
	if !ok {
		w.WriteHeader(400)
		return nil
	}
	if err != nil {
		return err
	}
	p := model.Project{Name: projectName, UserID: u.ID}
	has, err := p.GetWithUserIDAndName(session)
	if !has {
		w.WriteHeader(204)
		return nil
	}
	if err != nil {
		return err
	}

	if ok {
		// Save file
		err := model.UpdateFileContent(projectName, username, filePath, "", true, dir, p.Language)
		if err != nil {
			return err
		}
		w.WriteHeader(200)
	} else {
		w.WriteHeader(400)
	}
	return nil
}

// GetFileContentHandler is a handler for reading file content
func GetFileContentHandler(w http.ResponseWriter, r *http.Request) error {
	// get username from jwt
	ok, username := GetUserNameFromToken(r.Header.Get("Authorization"))
	if !ok {
		w.WriteHeader(401)
		return nil
	}

	// Read project id and file path from uri
	vars := mux.Vars(r)
	projectName := vars["projectname"]
	filePath := vars["filepath"]

	// Check if the file path is valid
	ok = checkFilePath(filePath)
	if ok {
		// Get project information
		session := MysqlEngine.NewSession()
		u := model.User{Username: username}
		ok, err := u.GetWithUsername(session)
		if !ok {
			w.WriteHeader(400)
			return nil
		}
		if err != nil {
			return err
		}
		p := model.Project{Name: projectName, UserID: u.ID}
		has, err := p.GetWithUserIDAndName(session)
		if !has {
			w.WriteHeader(204)
			return nil
		}
		if err != nil {
			return err
		}

		// Load file
		content, err := model.GetFileContent(projectName, username, filePath, p.Language)
		if err != nil {
			return err
		}
		w.WriteHeader(200)
		w.Write(content)
	} else {
		w.WriteHeader(400)
	}
	return nil
}

// DeleteFileHandler is a handler for delete file
func DeleteFileHandler(w http.ResponseWriter, r *http.Request) error {
	// get username from jwt
	ok, username := GetUserNameFromToken(r.Header.Get("Authorization"))
	if !ok {
		w.WriteHeader(401)
		return nil
	}
	// Read project id and file path from uri
	vars := mux.Vars(r)
	projectName := vars["projectname"]
	filePath := vars["filepath"]

	// Check if the file path is valid
	ok = checkFilePath(filePath)
	if ok {
		// Get project information
		session := MysqlEngine.NewSession()
		u := model.User{Username: username}
		ok, err := u.GetWithUsername(session)
		if !ok {
			w.WriteHeader(400)
			return nil
		}
		if err != nil {
			return err
		}
		p := model.Project{Name: projectName, UserID: u.ID}
		has, err := p.GetWithUserIDAndName(session)
		if !has {
			w.WriteHeader(204)
			return nil
		}
		if err != nil {
			return err
		}

		// Load file
		err = model.DeleteFile(projectName, username, filePath, p.Language)
		if err != nil {
			return err
		}
		w.WriteHeader(200)
	} else {
		w.WriteHeader(400)
	}
	return nil
}

// GetFileStructureHandler is handler for get project structure
func GetFileStructureHandler(w http.ResponseWriter, r *http.Request) error {
	// get username from jwt
	ok, username := GetUserNameFromToken(r.Header.Get("Authorization"))
	if !ok {
		w.WriteHeader(401)
		return nil
	}
	// Read project id
	vars := mux.Vars(r)
	projectName := vars["projectname"]

	// Get project information
	session := MysqlEngine.NewSession()
	u := model.User{Username: username}
	ok, err := u.GetWithUsername(session)
	if !ok {
		w.WriteHeader(400)
		return nil
	}
	if err != nil {
		return err
	}
	p := model.Project{Name: projectName, UserID: u.ID}
	has, err := p.GetWithUserIDAndName(session)
	if !has {
		w.WriteHeader(204)
		return nil
	}
	if err != nil {
		return err
	}

	// Get file structure
	structure, err := model.GetFileStructure(projectName, username, p.Language)
	if err != nil {
		return err
	}
	ret, err := json.Marshal(structure)
	if err != nil {
		return err
	}
	w.Write(ret)
	return nil
}
