package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/sysu-go-online/service-end/model"
)

// var username = "golang"

// UpdateFileHandler is a handler for update file
func UpdateFileHandler(w http.ResponseWriter, r *http.Request) error {
	// Read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	// Read project id and file path from uri
	vars := mux.Vars(r)
	projectName := vars["projectname"]
	userName := vars["username"]
	filePath := vars["filepath"]

	// Check if the file path is valid
	ok := checkFilePath(filePath)

	if ok {
		// Save file
		err := model.UpdateFileContent(projectName, userName, filePath, string(body), false, false)
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
		dir bool
	}
	// Judge if it is dir from body
	isDir := IsDir{}
	err = json.Unmarshal(body, &isDir)
	if err != nil {
		return err
	}
	dir := isDir.dir

	// Read project id and file path from uri
	vars := mux.Vars(r)
	projectName := vars["projectname"]
	userName := vars["username"]
	filePath := vars["filepath"]

	// Check if the file path is valid
	ok := checkFilePath(filePath)

	if ok {
		// Save file
		err := model.UpdateFileContent(projectName, userName, filePath, "", true, dir)
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
	// Read project id and file path from uri
	vars := mux.Vars(r)
	projectName := vars["projectname"]
	userName := vars["username"]
	filePath := vars["filepath"]

	// Check if the file path is valid
	ok := checkFilePath(filePath)
	if ok {
		// Load file
		content, err := model.GetFileContent(projectName, userName, filePath)
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
	// Read project id and file path from uri
	vars := mux.Vars(r)
	projectName := vars["projectname"]
	userName := vars["username"]
	filePath := vars["filepath"]

	// Check if the file path is valid
	ok := checkFilePath(filePath)
	if ok {
		// Load file
		err := model.DeleteFile(projectName, userName, filePath)
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
	// Read project id
	vars := mux.Vars(r)
	projectName := vars["projectname"]
	userName := vars["username"]

	// Get file structure
	structure, err := model.GetFileStructure(projectName, userName)
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
